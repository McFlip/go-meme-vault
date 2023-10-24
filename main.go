package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/McFlip/go-meme-vault/internal/models"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type TestHooks struct {
	DB *gorm.DB
}
type errRes struct {
	Err string `json:"error"`
}

func main() {
	var dbPath = flag.String("db-path", ":memory:", "path to sqlite file")
	var p = flag.Int("p", 8080, "port to listen on")
	flag.Parse()
	port := fmt.Sprintf(":%d", *p)
	const memePath = "public/img"

	database, err := gorm.Open(sqlite.Open(*dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	database.AutoMigrate(&models.Tag{})
	database.AutoMigrate(&models.Meme{})
	tagsModel := models.TagsModel{DB: database}
	memesModel := models.MemesModel{DB: database}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	testHooks := TestHooks{DB: database}
	testHooksRtr := chi.NewRouter()
	testHooksRtr.Delete("/nuke", testHooks.HandleNuke)
	testHooksRtr.Post("/tags", testHooks.HandleCreateTag)
	r.Mount("/api/testhooks", testHooksRtr)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmplFiles := []string{"templates/index.html", "templates/partials/menu.html"}
		tmpl := template.Must(template.ParseFiles(tmplFiles...))
		allTags, err := tagsModel.GetAll()
		if err != nil {
			respondWithErr(w, 500, "Error getting all tags")
		}
		tmpl.Execute(w, allTags)
	})

	r.Get("/memes", func(w http.ResponseWriter, r *http.Request) {
		qStr := r.URL.Query()
		tagIds := qStr["tag"]
		// log.Println(tagIds)
		tagSlice := make([]models.Tag, 0, len(tagIds))
		for _, id := range tagIds {
			id, err := strconv.Atoi(id)
			if err != nil {
				respondWithErr(w, 400, "tag query param must be int")
				return
			}
			t, err := tagsModel.GetByID(uint(id))
			if err != nil {
				respondWithErr(w, 500, "error getting tag by id")
				return
			}
			tagSlice = append(tagSlice, t)
		}

		tmplFiles := []string{"templates/memes.html", "templates/partials/selected_tags.html"}
		tmpl := template.Must(template.ParseFiles(tmplFiles...))
		tmpl.Execute(w, tagSlice)
	})

	r.Get("/memes/new", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/new_memes.html"))
		tmpl.Execute(w, nil)
	})

	r.Post("/memes/scan", func(w http.ResponseWriter, r *http.Request) {
		memeFullPath := fmt.Sprintf("%s/full", memePath)
		freshMemes, err := memesModel.Scan(memeFullPath)
		if err != nil {
			respondWithErr(w, 500, err.Error())
			return
		}

		tmplFiles := []string{"templates/memes_list.html", "templates/partials/meme_tn.html"}
		tmpl := template.Must(template.ParseFiles(tmplFiles...))
		tmpl.Execute(w, freshMemes)
	})

	staticHandler := http.StripPrefix("/public", http.FileServer(http.Dir("./public/")))
	r.Handle("/public", staticHandler)
	r.Handle("/public/*", staticHandler)

	log.Printf("🚀 Server launched on localhost%s", port)
	log.Printf("💾 Sqlite database saved at: %s", *dbPath)
	log.Fatal(http.ListenAndServe(port, r))
}

func (hooks *TestHooks) HandleNuke(w http.ResponseWriter, r *http.Request) {
	// GORM safeguards against unconditional deletes
	// Do a get all then for ea. Id delete
	// This is a 'soft' delete
	tm := models.TagsModel{DB: hooks.DB}
	allTags, err := tm.GetAll()
	if err != nil {
		log.Println(err)
		respondWithErr(w, 500, "error getting all tags")
		return
	}

	for _, t := range allTags {
		tm.DB.Delete(&t)
	}

	w.WriteHeader(200)
}

func (hooks *TestHooks) HandleCreateTag(w http.ResponseWriter, r *http.Request) {
	var tag models.Tag
	tagsModel := models.TagsModel{DB: hooks.DB}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tag)
	if err != nil {
		respondWithErr(w, 500, "unable to decode tag")
		return
	}
	tagsModel.Create(&tag)
	w.WriteHeader(200)
}

func respondWithErr(w http.ResponseWriter, code int, msg string) {
	resBody := errRes{
		Err: msg,
	}
	dat, err := json.Marshal(resBody)
	if err != nil {
		log.Printf("Error marshaling json body in respondWithErr: %s", err)
	}
	w.WriteHeader(code)
	w.Write(dat)
}
