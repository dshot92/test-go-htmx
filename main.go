package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Model struct {
	Name     string
	Path     string
	Category string
	Section  string
}

type PageData struct {
	Categories     []string
	Models         []Model
	ActiveCategory string
}

func getCategories() []string {
	var categories []string
	entries, err := os.ReadDir("static/models")
	if err != nil {
		return categories
	}

	for _, entry := range entries {
		if entry.IsDir() {
			categories = append(categories, entry.Name())
		}
	}
	return categories
}

func getModels(category string) []Model {
	var result []Model
	basePath := filepath.Join("static/models", category)

	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".glb") {
			relPath := strings.TrimPrefix(path, "static")
			parts := strings.Split(strings.TrimPrefix(relPath, "/models/"), "/")

			section := ""
			if len(parts) > 2 {
				section = parts[1]
			}

			model := Model{
				Name:     strings.TrimSuffix(info.Name(), ".glb"),
				Path:     relPath,
				Category: category,
				Section:  section,
			}
			result = append(result, model)
		}
		return nil
	})

	return result
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	}))

	// Create templates
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	// Serve static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		categories := getCategories()
		var activeCategory string
		if len(categories) > 0 {
			activeCategory = categories[0]
		}
		data := PageData{
			Categories:     categories,
			Models:         getModels(activeCategory),
			ActiveCategory: activeCategory,
		}
		tmpl.ExecuteTemplate(w, "index.html", data)
	})

	r.Get("/models", func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("category")
		if category == "" {
			categories := getCategories()
			if len(categories) > 0 {
				category = categories[0]
			}
		}
		models := getModels(category)
		tmpl.ExecuteTemplate(w, "grid.html", models)
	})

	r.Post("/upload/{category}", func(w http.ResponseWriter, r *http.Request) {
		category := chi.URLParam(r, "category")

		// Parse the multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("model")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Ensure the category directory exists
		categoryPath := filepath.Join("static/models", category)
		os.MkdirAll(categoryPath, 0755)

		// Create the destination file
		dst, err := os.Create(filepath.Join(categoryPath, handler.Filename))
		if err != nil {
			http.Error(w, "Error creating destination file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		// Return the updated grid
		models := getModels(category)
		tmpl.ExecuteTemplate(w, "grid.html", models)
	})

	log.Println("Server starting on :3000")
	http.ListenAndServe(":3000", r)
}
