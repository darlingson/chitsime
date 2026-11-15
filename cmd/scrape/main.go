package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"

	"github.com/darlingson/chitsime/internal/raw"
)

type urlSet struct {
	URLs []struct {
		Loc string `xml:"loc"`
	} `xml:"url"`
}

type source struct {
	baseURL string
	sitemap string
	parse   func(string, []byte) (raw.Job, error)
}

var sources = map[string]source{
	"careerad": {
		baseURL: "https://www.careeradmw.com",
		sitemap: "https://www.careeradmw.com/sitemap.xml",
		parse:   parseCareerAd,
	},
	"jobsearch": {
		baseURL: "https://jobsearchmalawi.com",
		sitemap: "https://jobsearchmalawi.com/job_listing-sitemap.xml",
		parse:   parseJobSearch,
	},
}

var blockTags = map[string]bool{
	"address": true, "article": true, "aside": true, "blockquote": true,
	"div": true, "dl": true, "dt": true, "dd": true, "fieldset": true,
	"figcaption": true, "figure": true, "footer": true, "form": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"header": true, "hr": true, "li": true, "main": true, "nav": true,
	"ol": true, "p": true, "pre": true, "section": true, "table": true,
	"tbody": true, "td": true, "tfoot": true, "th": true, "thead": true,
	"tr": true, "ul": true,
}

func main() {
	sourceName := flag.String("source", "careerad", "site to scrape: careerad or jobsearch")
	outPath := flag.String("out", "", "output JSON file (default data/<source>-jobs-raw.json)")
	minutes := flag.Int("minutes", 10, "maximum run time in minutes")
	delay := flag.Duration("delay", 300*time.Millisecond, "delay between page requests")
	limit := flag.Int("limit", 0, "maximum pages to fetch this run (0 = all)")
	sitemapURL := flag.String("sitemap", "", "override the job sitemap URL")
	fresh := flag.Bool("fresh", false, "start a new output file instead of resuming")
	flag.Parse()

	src, ok := sources[*sourceName]
	if !ok {
		log.Fatalf("unknown source %q, want careerad or jobsearch", *sourceName)
	}
	if *sitemapURL != "" {
		src.sitemap = *sitemapURL
	}
	if *outPath == "" {
		*outPath = filepath.Join("data", *sourceName+"-jobs-raw.json")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(*minutes)*time.Minute)
	defer cancel()

	if *fresh {
		if err := os.Remove(*outPath); err != nil && !os.IsNotExist(err) {
			log.Fatalf("remove output: %v", err)
		}
	}

	doc := loadOutput(*outPath, src.sitemap)
	doc.Source = src.baseURL
	doc.Sitemap = src.sitemap

	client := &http.Client{Timeout: 30 * time.Second}

	urls, err := fetchSitemap(ctx, client, src.sitemap)
	if err != nil {
		log.Fatalf("sitemap: %v", err)
	}
	if *limit > 0 && len(urls) > *limit {
		urls = urls[:*limit]
	}

	seen := make(map[string]bool, len(doc.Jobs))
	for _, job := range doc.Jobs {
		seen[job.SourceURL] = true
	}

	log.Printf("source %s: %d job URLs, %d already saved", *sourceName, len(urls), len(seen))

	started := time.Now()
	fetched, skipped, failed := 0, 0, 0

	for i, pageURL := range urls {
		if ctx.Err() != nil {
			break
		}
		if seen[pageURL] {
			skipped++
			continue
		}

		body, err := fetch(ctx, client, pageURL)
		if err != nil {
			failed++
			doc.Failures = append(doc.Failures, raw.Failure{
				SourceURL: pageURL,
				Error:     err.Error(),
				At:        time.Now().UTC().Format(time.RFC3339),
			})
			log.Printf("[%d/%d] FAIL %s: %v", i+1, len(urls), pageURL, err)
			_ = save(*outPath, &doc)
			continue
		}

		job, err := src.parse(pageURL, body)
		if err != nil {
			failed++
			doc.Failures = append(doc.Failures, raw.Failure{
				SourceURL: pageURL,
				Error:     err.Error(),
				At:        time.Now().UTC().Format(time.RFC3339),
			})
			log.Printf("[%d/%d] PARSE FAIL %s: %v", i+1, len(urls), pageURL, err)
			_ = save(*outPath, &doc)
			continue
		}

		doc.Jobs = append(doc.Jobs, job)
		seen[pageURL] = true
		fetched++
		log.Printf("[%d/%d] OK %s | %s", i+1, len(urls), job.Title, job.Employer)

		if fetched%5 == 0 {
			if err := save(*outPath, &doc); err != nil {
				log.Printf("save: %v", err)
			}
		}

		if i < len(urls)-1 {
			select {
			case <-ctx.Done():
			case <-time.After(*delay):
			}
		}
	}

	doc.Failures = pruneFailures(doc)

	sort.SliceStable(doc.Jobs, func(i, j int) bool {
		if doc.Jobs[i].PostedAt != doc.Jobs[j].PostedAt {
			return doc.Jobs[i].PostedAt > doc.Jobs[j].PostedAt
		}
		return doc.Jobs[i].Title < doc.Jobs[j].Title
	})

	if err := save(*outPath, &doc); err != nil {
		log.Fatalf("save: %v", err)
	}

	log.Printf("done in %s: fetched=%d skipped=%d failed=%d saved=%d -> %s",
		time.Since(started).Round(time.Second), fetched, skipped, failed, len(doc.Jobs), *outPath)
}

func fetchSitemap(ctx context.Context, client *http.Client, sitemapURL string) ([]string, error) {
	body, err := fetch(ctx, client, sitemapURL)
	if err != nil {
		return nil, err
	}

	var set urlSet
	if err := xml.Unmarshal(body, &set); err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(set.URLs))
	seen := make(map[string]bool, len(set.URLs))
	for _, entry := range set.URLs {
		loc := strings.TrimSpace(entry.Loc)
		if !strings.Contains(loc, "/job/") && !strings.Contains(loc, "/jobs/") {
			continue
		}
		if seen[loc] {
			continue
		}
		seen[loc] = true
		urls = append(urls, loc)
	}
	if len(urls) == 0 {
		return nil, errors.New("no job URLs found")
	}
	return urls, nil
}

func fetch(ctx context.Context, client *http.Client, pageURL string) ([]byte, error) {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Referer", "https://jobsearchmalawi.com/")

		resp, err := client.Do(req)
		if err == nil {
			body, readErr := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
			resp.Body.Close()
			if readErr != nil {
				lastErr = readErr
			} else if resp.StatusCode == http.StatusOK {
				return body, nil
			} else {
				lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
				if resp.StatusCode != http.StatusForbidden &&
					resp.StatusCode != http.StatusTooManyRequests &&
					resp.StatusCode < 500 {
					return nil, lastErr
				}
			}
		} else {
			lastErr = err
		}

		if attempt < 3 {
			wait := time.Duration(attempt) * 5 * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
		}
	}
	return nil, lastErr
}

type jobPosting struct {
	Title              string `json:"title"`
	DatePosted         string `json:"datePosted"`
	ValidThrough       string `json:"validThrough"`
	EmploymentType     any    `json:"employmentType"`
	HiringOrganization struct {
		Name string `json:"name"`
	} `json:"hiringOrganization"`
	JobLocation struct {
		Address struct {
			Locality string `json:"addressLocality"`
			Region   string `json:"addressRegion"`
			Country  string `json:"addressCountry"`
		} `json:"address"`
	} `json:"jobLocation"`
}

func parseCareerAd(pageURL string, body []byte) (raw.Job, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return raw.Job{}, fmt.Errorf("parse HTML: %w", err)
	}

	page := doc.Selection
	jp, _ := findJobPosting(page)
	title := strings.TrimSpace(jp.Title)
	if title == "" {
		title = cleanText(page.Find("h1").First())
	}
	if title == "" {
		return raw.Job{}, errors.New("title not found")
	}

	description := contentAfter(page, "Job Description")
	if description.Length() == 0 {
		description = page.Find("div.prose").Eq(0)
	}
	requirements := contentAfter(page, "Requirements")
	if requirements.Length() == 0 {
		requirements = page.Find("div.prose").Eq(1)
	}

	description.Find("script, style, ins, noscript, iframe, .code-block").Remove()
	requirements.Find("script, style, ins, noscript, iframe, .code-block").Remove()

	descriptionHTML, _ := description.Html()
	requirementsHTML, _ := requirements.Html()

	application := applicationText(description)
	if application == "" {
		application = applicationText(requirements)
	}

	links := make([]raw.Link, 0, 4)
	linkSeen := make(map[string]bool)
	addLinks(description, &links, linkSeen)
	addLinks(requirements, &links, linkSeen)

	employer := strings.TrimSpace(jp.HiringOrganization.Name)
	if employer == "" {
		employer = cleanText(page.Find("p.truncate.text-base.font-semibold").First())
	}

	location := pageLocation(page)
	if location == "" {
		location = joinLocation(jp)
	}

	return raw.Job{
		SourceURL:        pageURL,
		Title:            title,
		Employer:         employer,
		PostedAt:         strings.TrimSpace(jp.DatePosted),
		Deadline:         strings.TrimSpace(jp.ValidThrough),
		EmploymentType:   humanizeEmployment(jp.EmploymentType),
		Location:         location,
		Categories:       categories(page),
		ApplicationText:  application,
		DescriptionText:  renderSelection(description),
		DescriptionHTML:  strings.TrimSpace(descriptionHTML),
		RequirementsText: renderSelection(requirements),
		RequirementsHTML: strings.TrimSpace(requirementsHTML),
		Links:            links,
		Gated:            strings.Contains(cleanText(page.Find("body")), "Full details are for subscribers"),
		FetchedAt:        time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func findJobPosting(doc *goquery.Selection) (jobPosting, bool) {
	var jp jobPosting
	found := false

	doc.Find(`script[type="application/ld+json"]`).EachWithBreak(func(_ int, s *goquery.Selection) bool {
		var raw map[string]any
		if err := json.Unmarshal([]byte(s.Text()), &raw); err != nil {
			return true
		}
		if !isJobPosting(raw["@type"]) {
			return true
		}
		encoded, err := json.Marshal(raw)
		if err != nil {
			return true
		}
		if err := json.Unmarshal(encoded, &jp); err != nil {
			return true
		}
		found = true
		return false
	})

	return jp, found
}

func isJobPosting(v any) bool {
	switch t := v.(type) {
	case string:
		return t == "JobPosting"
	case []any:
		for _, item := range t {
			if s, ok := item.(string); ok && s == "JobPosting" {
				return true
			}
		}
	}
	return false
}

func contentAfter(doc *goquery.Selection, heading string) *goquery.Selection {
	h := doc.Find("h2").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return strings.EqualFold(strings.TrimSpace(s.Text()), heading)
	}).First()
	if h.Length() == 0 {
		return h
	}
	if sel := h.NextAll().Filter("div.prose").First(); sel.Length() > 0 {
		return sel
	}
	return h.NextAll().Filter("p.text-gray-700").First()
}

func applicationText(sel *goquery.Selection) string {
	var out string
	sel.Find("h2").EachWithBreak(func(_ int, h *goquery.Selection) bool {
		if strings.EqualFold(strings.TrimSpace(h.Text()), "How to apply") {
			out = renderSelection(h.NextUntil("h2, h3"))
			return false
		}
		return true
	})
	return out
}

func categories(doc *goquery.Selection) []string {
	out := []string{}
	seen := make(map[string]bool)
	doc.Find(`a[href*="/jobs?category="]`).Each(func(_ int, a *goquery.Selection) {
		text := cleanText(a)
		if text == "" || seen[text] {
			return
		}
		seen[text] = true
		out = append(out, text)
	})
	return out
}

func pageLocation(doc *goquery.Selection) string {
	meta := doc.Find("div.mt-4.flex.flex-wrap.gap-x-4.gap-y-2").First()
	if meta.Length() == 0 {
		return ""
	}
	icon := meta.Find("i.fa-map-marker-alt, i.fa-globe-africa, i.fa-location-dot").First()
	if icon.Length() == 0 {
		return ""
	}
	return cleanText(icon.Parent())
}

func joinLocation(jp jobPosting) string {
	addr := jp.JobLocation.Address
	parts := []string{}
	for _, part := range []string{addr.Locality, addr.Region, addr.Country} {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		duplicate := false
		for _, existing := range parts {
			if strings.EqualFold(existing, part) {
				duplicate = true
				break
			}
		}
		if !duplicate {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, ", ")
}

func humanizeEmployment(v any) string {
	var raw string
	switch t := v.(type) {
	case string:
		raw = t
	case []any:
		for _, item := range t {
			if s, ok := item.(string); ok {
				raw = s
				break
			}
		}
	}
	raw = strings.ReplaceAll(raw, "_", " ")
	words := strings.Fields(raw)
	for i, word := range words {
		lower := strings.ToLower(word)
		words[i] = strings.ToUpper(lower[:1]) + lower[1:]
	}
	return strings.Join(words, " ")
}

func parseJobSearch(pageURL string, body []byte) (raw.Job, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return raw.Job{}, fmt.Errorf("parse HTML: %w", err)
	}

	title := cleanText(doc.Find("h1.entry-title").First())
	if title == "" {
		return raw.Job{}, errors.New("title not found")
	}

	postedAt := doc.Find("time.entry-date").First().AttrOr("datetime", "")
	if postedAt == "" {
		postedAt = doc.Find(".date-posted time").First().AttrOr("datetime", "")
	}

	description := doc.Find(".job_description").First()
	description.Find("script, style, ins, noscript, iframe, .code-block").Remove()

	descriptionHTML, _ := description.Html()
	application := doc.Find(".job_application .application_details").First()

	links := make([]raw.Link, 0, 4)
	linkSeen := make(map[string]bool)
	addLinks(description, &links, linkSeen)
	addLinks(application, &links, linkSeen)

	employer := cleanText(doc.Find(".company p.name").First())
	if employer == "" {
		employer = cleanText(doc.Find(".company .company_header").First())
	}

	return raw.Job{
		SourceURL:       pageURL,
		Title:           title,
		Employer:        employer,
		PostedAt:        postedAt,
		EmploymentType:  cleanText(doc.Find("li.job-type").First()),
		Location:        cleanText(doc.Find("li.location").First()),
		Categories:      []string{},
		ApplicationText: cleanText(application),
		DescriptionText: renderSelection(description),
		DescriptionHTML: strings.TrimSpace(descriptionHTML),
		Links:           links,
		FetchedAt:       time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func addLinks(sel *goquery.Selection, links *[]raw.Link, seen map[string]bool) {
	sel.Find("a[href]").Each(func(_ int, a *goquery.Selection) {
		href, ok := a.Attr("href")
		if !ok {
			return
		}
		href = strings.TrimSpace(href)
		if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "javascript:") {
			return
		}
		text := cleanText(a)
		key := href + "|" + text
		if seen[key] {
			return
		}
		seen[key] = true
		*links = append(*links, raw.Link{Text: text, Href: href})
	})
}

func cleanText(sel *goquery.Selection) string {
	return strings.Join(strings.Fields(sel.Text()), " ")
}

func renderSelection(sel *goquery.Selection) string {
	var b strings.Builder
	for _, node := range sel.Nodes {
		renderNode(node, &b)
	}
	return normalizeText(b.String())
}

func renderNode(n *html.Node, b *strings.Builder) {
	if n.Type == html.TextNode {
		b.WriteString(n.Data)
		return
	}
	if n.Type != html.ElementNode {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			renderNode(child, b)
		}
		return
	}

	switch n.Data {
	case "script", "style", "noscript", "iframe", "svg", "ins":
		return
	case "br":
		b.WriteString("\n")
		return
	}

	block := blockTags[n.Data]
	if block && b.Len() > 0 {
		b.WriteString("\n")
	}
	if n.Data == "li" {
		b.WriteString("- ")
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		renderNode(child, b)
	}

	if block {
		b.WriteString("\n")
	}
}

func normalizeText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	blank := false

	for _, line := range lines {
		line = strings.Join(strings.Fields(line), " ")
		if line == "" {
			if len(out) > 0 && !blank {
				out = append(out, "")
			}
			blank = true
			continue
		}
		blank = false
		out = append(out, line)
	}

	return strings.TrimSpace(strings.Join(out, "\n"))
}

func pruneFailures(doc raw.Output) []raw.Failure {
	seen := make(map[string]bool, len(doc.Jobs))
	for _, job := range doc.Jobs {
		seen[job.SourceURL] = true
	}
	out := doc.Failures[:0]
	for _, failure := range doc.Failures {
		if !seen[failure.SourceURL] {
			out = append(out, failure)
		}
	}
	return out
}

func loadOutput(path, sitemap string) raw.Output {
	doc := raw.Output{Jobs: []raw.Job{}, Failures: []raw.Failure{}}
	body, err := os.ReadFile(path)
	if err != nil {
		return doc
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		log.Printf("existing output is not valid JSON, starting fresh: %v", err)
		return raw.Output{Jobs: []raw.Job{}, Failures: []raw.Failure{}}
	}
	if doc.Sitemap != "" && doc.Sitemap != sitemap {
		log.Printf("existing output belongs to %s, starting fresh", doc.Sitemap)
		return raw.Output{Jobs: []raw.Job{}, Failures: []raw.Failure{}}
	}
	if doc.Jobs == nil {
		doc.Jobs = []raw.Job{}
	}
	if doc.Failures == nil {
		doc.Failures = []raw.Failure{}
	}
	return doc
}

func save(path string, doc *raw.Output) error {
	doc.ScrapedAt = time.Now().UTC().Format(time.RFC3339)

	body, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, body, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
