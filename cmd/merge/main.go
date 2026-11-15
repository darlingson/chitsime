package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/darlingson/chitsime/internal/raw"
)

type entry struct {
	ID     int    `json:"id"`
	Source string `json:"source"`
	raw.Job
}

type document struct {
	GeneratedAt string  `json:"generated_at"`
	Count       int     `json:"count"`
	Jobs        []entry `json:"jobs"`
}

func main() {
	inputs := flag.String("in", "data/careerad-jobs-raw.json,data/jobsearch-jobs-raw.json", "comma-separated raw capture files")
	outPath := flag.String("out", "data/jobs-merged.json", "merged output file")
	flag.Parse()

	files := splitList(*inputs)
	if len(files) == 0 {
		log.Fatal("no input files")
	}

	jobs := make([]entry, 0, 700)
	for _, file := range files {
		doc, err := load(file)
		if err != nil {
			log.Fatalf("%s: %v", file, err)
		}

		name := sourceName(doc, file)
		for _, job := range doc.Jobs {
			jobs = append(jobs, entry{Source: name, Job: normalize(job)})
		}
		log.Printf("%s: %d jobs", name, len(doc.Jobs))
	}

	sort.SliceStable(jobs, func(i, j int) bool {
		return postedAfter(jobs[i].PostedAt, jobs[j].PostedAt)
	})
	for i := range jobs {
		jobs[i].ID = i + 1
	}

	doc := document{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Count:       len(jobs),
		Jobs:        jobs,
	}
	if err := save(*outPath, doc); err != nil {
		log.Fatalf("save: %v", err)
	}

	log.Printf("merged %d jobs -> %s", doc.Count, *outPath)
	reportDuplicates(jobs)
}

func load(path string) (raw.Output, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return raw.Output{}, err
	}

	var doc raw.Output
	if err := json.Unmarshal(body, &doc); err != nil {
		return raw.Output{}, fmt.Errorf("decode: %w", err)
	}
	if doc.Jobs == nil {
		doc.Jobs = []raw.Job{}
	}
	return doc, nil
}

func splitList(value string) []string {
	parts := strings.Split(value, ",")
	files := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			files = append(files, part)
		}
	}
	return files
}

func sourceName(doc raw.Output, path string) string {
	if u, err := url.Parse(doc.Source); err == nil {
		host := strings.TrimPrefix(u.Hostname(), "www.")
		switch {
		case strings.Contains(host, "careerad"):
			return "careerad"
		case strings.Contains(host, "jobsearch"):
			return "jobsearch"
		}
	}

	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return strings.TrimSuffix(base, "-jobs-raw")
}

func normalize(job raw.Job) raw.Job {
	if job.Categories == nil {
		job.Categories = []string{}
	}
	if job.Links == nil {
		job.Links = []raw.Link{}
	}
	return job
}

func postedAfter(a, b string) bool {
	ta, errA := time.Parse(time.RFC3339, a)
	tb, errB := time.Parse(time.RFC3339, b)

	switch {
	case errA == nil && errB == nil:
		return ta.After(tb)
	case errA == nil:
		return true
	case errB == nil:
		return false
	default:
		return false
	}
}

func reportDuplicates(jobs []entry) {
	urls := make(map[string]int, len(jobs))
	titles := make(map[string]int, len(jobs))
	for _, job := range jobs {
		urls[job.SourceURL]++
		key := strings.ToLower(strings.TrimSpace(job.Title) + " | " + strings.TrimSpace(job.Employer))
		titles[key]++
	}

	log.Printf("duplicate source URLs: %d", countRepeated(urls))
	log.Printf("duplicate title/employer pairs: %d", countRepeated(titles))
}

func countRepeated(values map[string]int) int {
	total := 0
	for _, count := range values {
		if count > 1 {
			total += count - 1
		}
	}
	return total
}

func save(path string, doc document) error {
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
