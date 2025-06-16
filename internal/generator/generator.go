package generator

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jchavarri/goranite/internal/config"
	"github.com/jchavarri/goranite/internal/content"
)

type Generator struct {
	config    *config.Config
	templates *template.Template
}

func New(configPath, templatesDir string) (*Generator, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	tmpl, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &Generator{
		config:    cfg,
		templates: tmpl,
	}, nil
}

func NewWithEmbedded(configPath string, embeddedFS embed.FS) (*Generator, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	tmpl := template.New("")

	// Read all .html files from the embedded filesystem
	err = fs.WalkDir(embeddedFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".html" {
			content, err := embeddedFS.ReadFile(path)
			if err != nil {
				return err
			}

			// Use just the filename as the template name
			templateName := filepath.Base(path)
			_, err = tmpl.New(templateName).Parse(string(content))
			if err != nil {
				return fmt.Errorf("failed to parse template %s: %w", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse embedded templates: %w", err)
	}

	return &Generator{
		config:    cfg,
		templates: tmpl,
	}, nil
}

func (g *Generator) Build(contentDir, staticDir, outputDir string) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Copy static files
	if err := g.copyStaticFiles(staticDir, outputDir); err != nil {
		return fmt.Errorf("failed to copy static files: %w", err)
	}

	// Load posts
	posts, err := content.LoadPosts(filepath.Join(contentDir, "posts"))
	if err != nil {
		return fmt.Errorf("failed to load posts: %w", err)
	}

	// Load pages
	pages, err := content.LoadPages(filepath.Join(contentDir, "pages"))
	if err != nil {
		// Pages directory is optional, so don't fail if it doesn't exist
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to load pages: %w", err)
		}
		pages = []config.Page{} // Empty slice if no pages directory
	}

	// Generate index page
	if err := g.generateIndex(posts, outputDir); err != nil {
		return fmt.Errorf("failed to generate index: %w", err)
	}

	// Generate individual post pages
	for _, post := range posts {
		if err := g.generatePost(&post, outputDir); err != nil {
			return fmt.Errorf("failed to generate post %s: %w", post.Title, err)
		}
	}

	// Generate individual pages
	for _, page := range pages {
		if err := g.generatePage(&page, outputDir); err != nil {
			return fmt.Errorf("failed to generate page %s: %w", page.Title, err)
		}
	}

	return nil
}

func (g *Generator) generateIndex(posts []config.Post, outputDir string) error {
	data := config.NewTemplateData(g.config)
	data.Posts = posts

	file, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return g.templates.ExecuteTemplate(file, "index.html", data)
}

func (g *Generator) generatePost(post *config.Post, outputDir string) error {
	data := config.NewTemplateData(g.config)
	data.Post = post
	data.Title = post.Title
	data.Description = post.Summary

	// Create post directory
	postDir := filepath.Join(outputDir, post.Slug)
	if err := os.MkdirAll(postDir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(postDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return g.templates.ExecuteTemplate(file, "post.html", data)
}

func (g *Generator) generatePage(page *config.Page, outputDir string) error {
	data := config.NewTemplateData(g.config)
	data.Page = page
	data.Title = page.Title
	data.Description = "" // Pages don't have summaries like posts do

	// Create page directory
	pageDir := filepath.Join(outputDir, page.Slug)
	if err := os.MkdirAll(pageDir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(pageDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return g.templates.ExecuteTemplate(file, page.Template, data)
}

func (g *Generator) copyStaticFiles(staticDir, outputDir string) error {
	return filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(staticDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(outputDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, src, info.Mode())
	})
}
