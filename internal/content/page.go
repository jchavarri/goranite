package content

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"

	"github.com/javierchavarri/goranite/internal/config"
)

type PageMatter struct {
	Title    string `json:"title"`
	Template string `json:"template"` // Which template to use (defaults to "page.html")
}

func LoadPage(path string) (*config.Page, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse frontmatter
	frontmatterData, body, err := parseFrontmatter(content)
	if err != nil {
		return nil, err
	}

	// Convert map to our struct
	var matter PageMatter
	if frontmatterData != nil {
		jsonBytes, err := json.Marshal(frontmatterData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal frontmatter: %w", err)
		}

		if err := json.Unmarshal(jsonBytes, &matter); err != nil {
			return nil, fmt.Errorf("failed to unmarshal frontmatter: %w", err)
		}
	}

	// Default template if not specified
	if matter.Template == "" {
		matter.Template = "page.html"
	}

	// Convert markdown to HTML with syntax highlighting
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithXHTML(),
		),
	)

	// Add our custom syntax highlighting renderer
	md.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(newSyntaxHighlightRenderer(), 200),
		),
	)

	var buf strings.Builder
	if err := md.Convert(body, &buf); err != nil {
		return nil, err
	}

	htmlContent := buf.String()

	// Generate slug from filename
	filename := filepath.Base(path)
	slug := strings.TrimSuffix(filename, ".md")

	return &config.Page{
		Title:    matter.Title,
		Content:  template.HTML(htmlContent),
		URL:      "/" + slug + "/",
		Slug:     slug,
		Template: matter.Template,
	}, nil
}

func LoadPages(contentDir string) ([]config.Page, error) {
	var pages []config.Page

	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		page, err := LoadPage(path)
		if err != nil {
			return fmt.Errorf("failed to load page %s: %w", path, err)
		}

		pages = append(pages, *page)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return pages, nil
}
