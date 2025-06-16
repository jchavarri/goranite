package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jchavarri/goranite/internal/content"
	"github.com/jchavarri/goranite/internal/generator"
)

func main() {
	var (
		buildCmd = flag.Bool("build", false, "Build the static site")
		serveCmd = flag.Bool("serve", false, "Start development server")
		newCmd   = flag.String("new", "", "Create new post with given title")
		siteDir  = flag.String("site", ".", "Path to site directory")
	)
	flag.Parse()

	switch {
	case *buildCmd:
		fmt.Println("🔨 Building static site...")
		if err := buildSite(*siteDir); err != nil {
			log.Fatalf("Build failed: %v", err)
		}
		fmt.Println("✅ Site built successfully!")

	case *serveCmd:
		fmt.Println("🚀 Starting development server...")
		if err := serveSite(*siteDir); err != nil {
			log.Fatalf("Server failed: %v", err)
		}

	case *newCmd != "":
		fmt.Printf("📝 Creating new post: %s\n", *newCmd)
		if err := createNewPost(*newCmd); err != nil {
			log.Fatalf("Failed to create post: %v", err)
		}
		fmt.Println("✅ Post created successfully!")

	default:
		fmt.Println("🪨 Goranite - Static Site Generator")
		fmt.Println("Solid as granite, powered by Go")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  -build       Build the static site")
		fmt.Println("  -serve       Start development server with auto-rebuild")
		fmt.Println("  -new 'Title' Create new post")
		fmt.Println("  -site path   Path to site directory (default: .)")
		flag.PrintDefaults()
	}
}

func buildSite(siteDir string) error {
	configPath := filepath.Join(siteDir, "config.json")
	contentDir := filepath.Join(siteDir, "content")
	staticDir := filepath.Join(siteDir, "static")
	outputDir := filepath.Join(siteDir, "public")

	// Find templates directory relative to executable
	templatesDir, err := getTemplatesDir()
	if err != nil {
		return fmt.Errorf("failed to find templates directory: %w", err)
	}

	// Clean output directory first
	fmt.Println("🧹 Cleaning output directory...")
	if err := os.RemoveAll(outputDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clean output directory: %w", err)
	}

	gen, err := generator.New(configPath, templatesDir)
	if err != nil {
		return err
	}

	// Build the site first
	if err := gen.Build(contentDir, staticDir, outputDir); err != nil {
		return err
	}

	// Generate chroma CSS after the build
	if err := content.GenerateChromaCSS(outputDir); err != nil {
		return fmt.Errorf("failed to generate chroma CSS: %w", err)
	}

	return nil
}

func serveSite(siteDir string) error {
	// Build initially
	fmt.Println("🔨 Initial build...")
	if err := buildSite(siteDir); err != nil {
		return fmt.Errorf("failed to build site: %w", err)
	}

	// Start file watcher in background
	go watchAndRebuild(siteDir)

	// Set up file server
	publicDir := filepath.Join(siteDir, "public")
	fmt.Printf("🌐 Serving site at http://localhost:8080\n")
	fmt.Printf("📁 Serving files from: %s\n", publicDir)
	fmt.Println("👀 Watching for changes... Edit files and refresh browser!")
	fmt.Println("Press Ctrl+C to stop")

	// Add cache-busting headers for development
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Disable caching for development
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		http.FileServer(http.Dir(publicDir)).ServeHTTP(w, r)
	})

	return http.ListenAndServe(":8080", handler)
}

func watchAndRebuild(siteDir string) {
	contentDir := filepath.Join(siteDir, "content")
	configFile := filepath.Join(siteDir, "config.json")

	// Get templates directory dynamically
	templatesDir, err := getTemplatesDir()
	if err != nil {
		fmt.Printf("❌ Warning: Could not find templates directory for watching: %v\n", err)
		templatesDir = "templates" // fallback
	}

	// Get initial modification times
	lastMod := getLastModTime(contentDir, configFile, templatesDir)

	for {
		time.Sleep(1 * time.Second) // Check every second

		currentMod := getLastModTime(contentDir, configFile, templatesDir)

		if currentMod.After(lastMod) {
			fmt.Println("📝 Changes detected, rebuilding...")
			if err := buildSite(siteDir); err != nil {
				fmt.Printf("❌ Build error: %v\n", err)
			} else {
				fmt.Println("✅ Rebuilt successfully! Refresh your browser.")
			}
			lastMod = currentMod
		}
	}
}

func getLastModTime(paths ...string) time.Time {
	var latest time.Time

	for _, path := range paths {
		filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}

			// Only check relevant files
			if filepath.Ext(filePath) == ".md" ||
				filepath.Ext(filePath) == ".json" ||
				filepath.Ext(filePath) == ".html" {
				if info.ModTime().After(latest) {
					latest = info.ModTime()
				}
			}
			return nil
		})
	}

	return latest
}

func createNewPost(title string) error {
	// TODO: Implement post creation
	fmt.Printf("Creating post '%s'... (not implemented yet)\n", title)
	return nil
}

func getTemplatesDir() (string, error) {
	// Get the path of the current executable
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	execDir := filepath.Dir(execPath)

	// Try templates relative to executable first
	templatesDir := filepath.Join(execDir, "templates")
	if _, err := os.Stat(templatesDir); err == nil {
		return templatesDir, nil
	}

	// Fallback to current working directory (for development)
	templatesDir = "templates"
	if _, err := os.Stat(templatesDir); err == nil {
		return templatesDir, nil
	}

	return "", fmt.Errorf("templates directory not found")
}
