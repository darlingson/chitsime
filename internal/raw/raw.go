package raw

type Link struct {
	Text string `json:"text"`
	Href string `json:"href"`
}

type Job struct {
	SourceURL        string   `json:"source_url"`
	Title            string   `json:"title"`
	Employer         string   `json:"employer"`
	PostedAt         string   `json:"posted_at"`
	Deadline         string   `json:"deadline"`
	EmploymentType   string   `json:"employment_type"`
	Location         string   `json:"location"`
	Categories       []string `json:"categories"`
	ApplicationText  string   `json:"application_text"`
	DescriptionText  string   `json:"description_text"`
	DescriptionHTML  string   `json:"description_html"`
	RequirementsText string   `json:"requirements_text"`
	RequirementsHTML string   `json:"requirements_html"`
	Links            []Link   `json:"links"`
	Gated            bool     `json:"gated"`
	FetchedAt        string   `json:"fetched_at"`
}

type Failure struct {
	SourceURL string `json:"source_url"`
	Error     string `json:"error"`
	At        string `json:"at"`
}

type Output struct {
	ScrapedAt string    `json:"scraped_at"`
	Source    string    `json:"source"`
	Sitemap   string    `json:"sitemap"`
	Jobs      []Job     `json:"jobs"`
	Failures  []Failure `json:"failures"`
}
