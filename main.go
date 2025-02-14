package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"test-go-htmx/templates"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Serve static files
	e.Static("/static", "static")

	// Routes - Remove column-related routes
	e.GET("/", handleIndex)
	e.GET("/category/:category", handleFilter)
	e.GET("/view/:category/:name", handleView)

	// Start server
	log.Printf("Starting server on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

func loadModels() ([]templates.Model, []string) {
	var models []templates.Model
	categories := make(map[string]bool)

	// Log the current working directory
	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting working directory: %v", err)
	} else {
		log.Printf("Current working directory: %s", pwd)
	}

	// Check if the models directory exists
	if _, err := os.Stat("static/models"); os.IsNotExist(err) {
		log.Printf("static/models directory does not exist!")
	}

	err = filepath.Walk("static/models", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v", path, err)
			return err
		}

		// Log all files/directories being processed
		log.Printf("Processing path: %s (isDir: %v)", path, info.IsDir())

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(info.Name()), ".glb") {
			return nil
		}

		// Get the path relative to models directory
		relPath, _ := filepath.Rel("static/models", filepath.Dir(path))
		pathParts := strings.Split(relPath, string(filepath.Separator))

		// The first directory is the category
		category := pathParts[0]

		// The section is either the second directory or empty
		section := ""
		if len(pathParts) > 1 {
			section = pathParts[1]
		}

		relPath = "/" + filepath.ToSlash(path)

		model := templates.Model{
			ID:       strings.TrimSuffix(info.Name(), ".glb"),
			Name:     strings.TrimSuffix(info.Name(), ".glb"),
			Path:     relPath,
			Category: category,
			Section:  section,
			URL:      relPath,
		}

		log.Printf("Found model: %+v", model)

		models = append(models, model)
		categories[category] = true
		return nil
	})
	if err != nil {
		log.Printf("Error loading models: %v", err)
	}

	var categoryList []string
	for category := range categories {
		categoryList = append(categoryList, category)
	}
	sort.Strings(categoryList)

	log.Printf("Loaded categories: %v", categoryList)
	log.Printf("Total models loaded: %d", len(models))

	return models, categoryList
}

func handleFilter(c echo.Context) error {
	models, categories := loadModels()
	category := c.Param("category")

	log.Printf("Filter request - Category: %s", category)
	log.Printf("Available categories: %v", categories)
	log.Printf("Total models before filtering: %d", len(models))

	// Check if this is an HTMX request
	if c.Request().Header.Get("HX-Request") == "true" {
		var filteredModels []templates.Model
		if category == "all" {
			filteredModels = models
		} else {
			for _, model := range models {
				if strings.EqualFold(model.Category, category) {
					filteredModels = append(filteredModels, model)
				}
			}
		}
		// Use fixed column count of 5
		return templates.ModelGrid(filteredModels, 5).Render(c.Request().Context(), c.Response().Writer)
	}

	// For regular browser requests, return the full page
	if category == "all" || category == "" {
		log.Printf("Returning all %d models", len(models))
		return templates.Index(models, categories, "all").Render(c.Request().Context(), c.Response().Writer)
	}

	var filteredModels []templates.Model
	for _, model := range models {
		log.Printf("Comparing model category '%s' with requested category '%s'", model.Category, category)
		if strings.EqualFold(model.Category, category) {
			filteredModels = append(filteredModels, model)
		}
	}

	log.Printf("Found %d models in category %s", len(filteredModels), category)
	if len(filteredModels) == 0 {
		log.Printf("WARNING: No models found for category: %s", category)
	}

	return templates.Index(filteredModels, categories, category).Render(c.Request().Context(), c.Response().Writer)
}

func handleIndex(c echo.Context) error {
	models, categories := loadModels()
	log.Printf("Index request - Available categories: %v", categories)
	return templates.Index(models, categories, "all").Render(c.Request().Context(), c.Response().Writer)
}

func handleView(c echo.Context) error {
	models, _ := loadModels()
	category := c.Param("category")
	name := c.Param("name")

	log.Printf("View request - Category: %s, Name: %s", category, name)

	var targetModel templates.Model
	for _, model := range models {
		if strings.EqualFold(model.Category, category) && strings.EqualFold(model.Name, name) {
			targetModel = model
			break
		}
	}

	if targetModel.Name == "" {
		log.Printf("Model not found - Category: %s, Name: %s", category, name)
		return c.String(http.StatusNotFound, "Model not found")
	}

	viewData := templates.ModelViewData{
		Model: targetModel,
		Title: targetModel.Name,
	}
	return templates.ModelViewer(viewData).Render(c.Request().Context(), c.Response().Writer)
}
