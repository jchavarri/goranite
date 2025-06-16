package config

import (
	"encoding/json"
	"html/template"
	"os"
	"time"
)

type Site struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	URL         string `json:"url"`
}

type Build struct {
	OutputDir    string `json:"output_dir"`
	PostsPerPage int    `json:"posts_per_page"`
}

type Social struct {
	Twitter string `json:"twitter"`
	GitHub  string `json:"github"`
}

type Config struct {
	Site   Site   `json:"site"`
	Build  Build  `json:"build"`
	Social Social `json:"social"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set defaults
	if config.Build.OutputDir == "" {
		config.Build.OutputDir = "public"
	}
	if config.Build.PostsPerPage == 0 {
		config.Build.PostsPerPage = 10
	}

	return &config, nil
}

type TemplateData struct {
	SiteTitle       string
	SiteDescription template.HTML
	SiteURL         string
	Author          string
	TwitterHandle   string
	Year            int
	Title           string
	Description     string
	Posts           []Post
	Post            *Post
	Page            *Page
}

type Post struct {
	Title       string
	Date        time.Time
	Tags        []string
	Summary     string
	Image       string        `json:"image"` // Featured image for social sharing
	Content     template.HTML `json:"content"`
	URL         string
	ReadingTime int
	Slug        string
}

type Page struct {
	Title       string
	Description string // Meta description for SEO
	Content     template.HTML
	URL         string
	Slug        string
	Template    string // Which template to use (e.g., "about.html", "contact.html")
}

func NewTemplateData(config *Config) *TemplateData {
	return &TemplateData{
		SiteTitle:       config.Site.Title,
		SiteDescription: template.HTML(config.Site.Description),
		SiteURL:         config.Site.URL,
		Author:          config.Site.Author,
		TwitterHandle:   config.Social.Twitter,
		Year:            time.Now().Year(),
	}
}
