package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"github.com/McFlip/go-meme-vault/components"
	"github.com/McFlip/go-meme-vault/internal/models"
)

type TestHooks struct {
	DB *gorm.DB
}
type errRes struct {
	Err string `json:"error"`
}

func main() {
	dbPath := flag.String("db-path", ":memory:", "path to sqlite file")
	p := flag.Int("p", 8080, "port to listen on")
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
	testHooksRtr.Post("/memes", testHooks.HandleCreateMeme)
	testHooksRtr.Post("/memes/{memeId}/addtag/{tagId}", testHooks.HandleAddTag)
	r.Mount("/api/testhooks", testHooksRtr)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		allTags, err := tagsModel.GetAll()
		if err != nil {
			respondWithErr(w, 500, "Error getting all tags")
		}
		idxComponent := components.Index(allTags)
		if isHtmx(r) {
			idxComponent.Render(r.Context(), w)
		} else {
			components.Layout(idxComponent).Render(r.Context(), w)
		}
	})

	r.Get("/memes", func(w http.ResponseWriter, r *http.Request) {
		qStr := r.URL.Query()
		tagIds := qStr["tag"]
		// log.Println(tagIds)
		tagSlice := make([]models.Tag, 0, len(tagIds))
		tagMapBool := make(map[uint]bool)
		var memesSlice []models.Meme
		memesMap := make(map[uint]models.Meme)
		memesCountMap := make(map[uint]int)

		for _, id := range tagIds {
			id, err := strconv.Atoi(id)
			if err != nil {
				respondWithErr(w, 400, "tag query param must be int")
				return
			}
			tagMapBool[uint(id)] = false
			t, err := tagsModel.GetByID(uint(id))
			if err != nil {
				respondWithErr(w, 500, "error getting tag by id")
				return
			}
			tagSlice = append(tagSlice, t)
			for _, m := range t.Memes {
				// log.Println(m)
				memesMap[m.ID] = *m
				if _, ok := memesCountMap[m.ID]; ok {
					memesCountMap[m.ID] = memesCountMap[m.ID] + 1
				} else {
					memesCountMap[m.ID] = 1
				}
			}
		}
		// log.Println(memesCountMap)

		for k, v := range memesCountMap {
			if v == len(tagSlice) {
				memesSlice = append(memesSlice, memesMap[k])
			}
		}

		/*
				for _, m := range memesMap {
		      hasAllTags := false
		      for t := range tagMapBool {
		        tagMapBool[t] = false
		      }
		      for _, t := range m.Tags {
		        log.Println("FUBAR")
		        log.Println(t)
		        log.Println(tagMapBool)
		        if v, ok := tagMapBool[t.ID]; ok {
		          v = true
		          log.Printf("v: %v, map: %v", v, tagMapBool[t.ID])
		          tagMapBool[t.ID] = true
		        }
		      }
		      for _, t := range tagMapBool {
		        hasAllTags = hasAllTags && t
		      }
		      if hasAllTags {memesSlice = append(memesSlice, m)}
				}
		*/

		memesComponent := components.Memes(tagSlice, memesSlice)
		if isHtmx(r) {
			err := memesComponent.Render(r.Context(), w)
			if err != nil {
				respondWithErr(w, 500, err.Error())
			}
		} else {
			err := components.Layout(memesComponent).Render(r.Context(), w)
			if err != nil {
				respondWithErr(w, 500, err.Error())
			}
		}
	})

	r.Get("/memes/untagged", func(w http.ResponseWriter, r *http.Request) {
		var tagSlice []models.Tag
		memesSlice, err := memesModel.GetUntagged()
		if err != nil {
			respondWithErr(w, 500, err.Error())
		}
		memesComponent := components.Memes(tagSlice, memesSlice)
		if isHtmx(r) {
			err := memesComponent.Render(r.Context(), w)
			if err != nil {
				respondWithErr(w, 500, err.Error())
			}
		} else {
			err := components.Layout(memesComponent).Render(r.Context(), w)
			if err != nil {
				respondWithErr(w, 500, err.Error())
			}
		}
	})

	r.Get("/memes/{memeId}", func(w http.ResponseWriter, r *http.Request) {
		memeId := chi.URLParam(r, "memeId")
		id, err := strconv.Atoi(memeId)
		if err != nil {
			respondWithErr(w, 400, "memeId must be an int")
			return
		}
		meme, err := memesModel.GetByID(uint(id))
		if err != nil {
			respondWithErr(w, 500, "error getting meme by id")
			return
		}
		err = components.MemeModal(meme).Render(r.Context(), w)
		if err != nil {
			respondWithErr(w, 500, err.Error())
		}
	})

	r.Get("/memes/new", func(w http.ResponseWriter, r *http.Request) {
		newMemesComponent := components.New_Memes()
		if isHtmx(r) {
			newMemesComponent.Render(r.Context(), w)
		} else {
			components.Layout(newMemesComponent).Render(r.Context(), w)
		}
	})

	r.Post("/memes/scan", func(w http.ResponseWriter, _ *http.Request) {
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

	r.Post("/tags", func(w http.ResponseWriter, r *http.Request) {
		memeId, err := strconv.Atoi(r.PostFormValue("memeId"))
		if err != nil {
			respondWithErr(w, 400, "bad form data")
			return
		}
		meme, err := memesModel.GetByID(uint(memeId))
		if err != nil {
			respondWithErr(w, 404, "missing meme ref")
			return
		}

		tag := models.Tag{
			Name:  r.PostFormValue("search"),
			Memes: []*models.Meme{&meme},
		}
		res := tagsModel.Create(&tag)
		if res.Error != nil {
			respondWithErr(w, 500, res.Error.Error())
			return
		}
		newTag := components.TagParams{
			MemeId: strconv.Itoa(int(meme.ID)),
			TagId:  strconv.Itoa(int(tag.ID)),
			Name:   tag.Name,
		}
		err = components.Tag(newTag).Render(r.Context(), w)
		if err != nil {
			respondWithErr(w, 500, err.Error())
		}
	})

	r.Post("/tags/search", func(w http.ResponseWriter, r *http.Request) {
		qStr := r.PostFormValue("search")
		if qStr == "" {
			w.WriteHeader(200)
			return
		}
		memeIdStr := r.PostFormValue("memeId")
		memeId := 0
		if memeIdStr != "" {
			memeId, err = strconv.Atoi(memeIdStr)
			if err != nil {
				respondWithErr(w, 400, "bad form data")
				return
			}
		}

		tags, err := tagsModel.Search(qStr)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				noTag := "<p>No tag found</p>"
				w.Write([]byte(noTag))
				return
			}
			respondWithErr(w, 500, err.Error())
			return
		}

		if memeId == 0 {
			components.TagsList(tags).Render(r.Context(), w)
		} else {
			// filter out tags already attached to meme
			filteredTags, err := memesModel.FilterTags(uint(memeId), tags)
			if err != nil {
				respondWithErr(w, 500, err.Error())
				return
			}
			// render AddTagsList
			components.AddTagsList(memeIdStr, filteredTags).Render(r.Context(), w)
		}
	})

	r.Patch("/memes/{memeId}/tags/{tagId}", func(w http.ResponseWriter, r *http.Request) {
		memeIdStr := chi.URLParam(r, "memeId")
		memeId, err := strconv.Atoi(memeIdStr)
		if err != nil {
			respondWithErr(w, 400, "memeId must be an int")
			return
		}
		tagIdStr := chi.URLParam(r, "tagId")
		tagId, err := strconv.Atoi(tagIdStr)
		if err != nil {
			respondWithErr(w, 400, "tagId must be an int")
			return
		}
		tag, err := tagsModel.GetByID(uint(tagId))
		if err != nil {
			respondWithErr(w, 500, err.Error())
			return
		}

		_, err = memesModel.AddTag(uint(memeId), tag)
		if err != nil {
			respondWithErr(w, 500, err.Error())
			return
		}

		newTag := components.TagParams{
			MemeId: memeIdStr,
			TagId:  tagIdStr,
			Name:   tag.Name,
		}
		components.Tag(newTag).Render(r.Context(), w)
	})

	// remove a tag from a meme
	r.Delete("/memes/{memeId}/tags/{tagId}", func(w http.ResponseWriter, r *http.Request) {
		memeId, err := strconv.Atoi(chi.URLParam(r, "memeId"))
		if err != nil {
			respondWithErr(w, 400, "memeId must be an int")
			return
		}
		tagId, err := strconv.Atoi(chi.URLParam(r, "tagId"))
		if err != nil {
			respondWithErr(w, 400, "tagId must be an int")
			return
		}

		tag, err := tagsModel.GetByID(uint(tagId))
		if err != nil {
			respondWithErr(w, 500, err.Error())
		}

		err = memesModel.RemoveTag(uint(memeId), tag)
		if err != nil {
			respondWithErr(w, 500, err.Error())
		}
	})

	staticHandler := http.StripPrefix("/public", http.FileServer(http.Dir("./public/")))
	r.Handle("/public", staticHandler)
	r.Handle("/public/*", staticHandler)

	log.Printf("🚀 Server launched on http://localhost%s", port)
	log.Printf("💾 Sqlite database saved at: %s", *dbPath)
	log.Fatal(http.ListenAndServe(port, r))
}

func (hooks *TestHooks) HandleNuke(w http.ResponseWriter, r *http.Request) {
	hooks.DB.Unscoped().Exec("DELETE FROM tags")
	hooks.DB.Unscoped().Exec("DELETE FROM memes")
	hooks.DB.Unscoped().Exec("DELETE FROM meme_tags")

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

func (hooks *TestHooks) HandleCreateMeme(w http.ResponseWriter, r *http.Request) {
	var meme models.Meme
	memesModel := models.MemesModel{DB: hooks.DB}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&meme)
	if err != nil {
		respondWithErr(w, 500, "unable to decode meme")
		return
	}
	memesModel.Create(&meme)
	w.WriteHeader(200)
}

func (hooks *TestHooks) HandleAddTag(w http.ResponseWriter, r *http.Request) {
	tagsModel := models.TagsModel{DB: hooks.DB}
	memesModel := models.MemesModel{DB: hooks.DB}
	memeId, err := strconv.Atoi(chi.URLParam(r, "memeId"))
	if err != nil {
		respondWithErr(w, 400, "memeId must be an int")
		return
	}
	tagId, err := strconv.Atoi(chi.URLParam(r, "tagId"))
	if err != nil {
		respondWithErr(w, 400, "tagId must be an int")
		return
	}

	tag, err := tagsModel.GetByID(uint(tagId))
	if err != nil {
		respondWithErr(w, 500, err.Error())
	}
	meme, err := memesModel.AddTag(uint(memeId), tag)
	respondWithJSON(w, 200, meme)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		respondWithErr(w, 500, fmt.Sprintf("Error marshaling json in respondWithJSON: %s", err))
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
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

func isHtmx(r *http.Request) bool {
	isHtmx := false
	htmxHeader := r.Header["Hx-Request"]
	if len(htmxHeader) > 0 {
		if htmxHeader[0] == "true" {
			isHtmx = true
		}
	}
	return isHtmx
}
