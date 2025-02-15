package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Model struct {
	Name     string
	Path     string
	Category string
	Section  string
	BlobURL  string // Add this field for client-side uploads
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
				Path:     "/models/" + strings.Join(parts, "/"),
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
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Create the static/models directory if it doesn't exist
	os.MkdirAll("static/models", 0755)

	// Add a new route for direct file viewing
	r.Post("/view/upload", func(w http.ResponseWriter, r *http.Request) {
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

		// Create a temporary directory if it doesn't exist
		tempDir := filepath.Join("static", "temp")
		os.MkdirAll(tempDir, 0755)

		// Create a temporary file
		tempFile := filepath.Join(tempDir, handler.Filename)
		dst, err := os.Create(tempFile)
		if err != nil {
			http.Error(w, "Error creating temporary file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		// Create a model for the temporary file
		model := Model{
			Name: strings.TrimSuffix(handler.Filename, ".glb"),
			Path: "/temp/" + handler.Filename,
		}

		tmpl.ExecuteTemplate(w, "viewer.html", model)

		// Schedule cleanup of the temporary file
		go func() {
			time.Sleep(1 * time.Hour) // Keep the file for 1 hour
			os.Remove(tempFile)
		}()
	})

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

	r.Get("/view/*", func(w http.ResponseWriter, r *http.Request) {
		modelPath := chi.URLParam(r, "*")

		// Check if it's a blob URL
		if strings.HasPrefix(modelPath, "blob:") {
			model := Model{
				Name:    "Uploaded Model",
				BlobURL: modelPath,
			}
			tmpl.ExecuteTemplate(w, "viewer.html", model)
			return
		}

		// Remove any leading /models/ from the path
		modelPath = strings.TrimPrefix(modelPath, "models/")
		modelPath = strings.TrimPrefix(modelPath, "/models/")

		// Extract category and filename from the path
		parts := strings.Split(modelPath, "/")
		if len(parts) < 2 {
			http.Error(w, "Invalid model path", http.StatusBadRequest)
			return
		}

		category := parts[0]
		filename := parts[len(parts)-1]
		section := ""
		if len(parts) > 2 {
			section = parts[1]
		}

		// Verify the file exists
		fullPath := filepath.Join("static/models", modelPath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			http.Error(w, "Model not found", http.StatusNotFound)
			return
		}

		// Create the model with the correct path
		model := Model{
			Name:     strings.TrimSuffix(filename, ".glb"),
			Path:     "/models/" + modelPath,
			Category: category,
			Section:  section,
		}

		tmpl.ExecuteTemplate(w, "viewer.html", model)
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

		// Redirect to the viewer page without adding extra /models/
		modelPath := category + "/" + handler.Filename
		http.Redirect(w, r, "/view/"+modelPath, http.StatusSeeOther)
	})

	log.Println("Server starting on :3000")
	http.ListenAndServe(":3000", r)
}
