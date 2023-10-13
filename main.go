package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/McFlip/go-meme-vault/internal/tags"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

const port = ":8080"

func main() {
	var dbPath = flag.String("db-path", ":memory:", "path to sqlite file")
	flag.Parse()

	database, err := gorm.Open(sqlite.Open(*dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	database.AutoMigrate(&tags.Tag{})
	tagsModel := tags.TagsModel{DB: database}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		allTags := tagsModel.GetAll()
		var tagNames []string
		for _, t := range allTags {
			tagNames = append(tagNames, t.Name)
		}
		tmpl.Execute(w, tagNames)
	})

	r.Get("/memes", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/memes.html"))
		tmpl.Execute(w, nil)
	})

	log.Printf("🚀 Server launched on localhost%s", port)
	log.Printf("💾 Sqlite database saved at: %s", *dbPath)
	log.Fatal(http.ListenAndServe(port, r))
}
