package main

import (
	"net/http"
	"strconv"

	"github.com/McFlip/go-meme-vault/components"
	"github.com/McFlip/go-meme-vault/internal/models"
	"github.com/go-chi/chi/v5"
)

type get_Meme_By_Id_Params struct {
	idx []string
	id  uint
}

var get_Meme_By_Id_Parser = func(r *http.Request) (get_Meme_By_Id_Params, error) {
	qStr := r.URL.Query()
	// TODO: validate idx
	idx := qStr["idx"]
	memeId := chi.URLParam(r, "memeId")
	id, err := strconv.Atoi(memeId)
	if err != nil {
		return get_Meme_By_Id_Params{}, err
	}
	return get_Meme_By_Id_Params{
		idx: idx,
		id:  uint(id),
	}, nil
}

// http.HandlerFunc
var get_Meme_By_Id = func(memesModel models.MemesModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := get_Meme_By_Id_Parser(r)
		if err != nil {
			respondWithErr(w, 400, "memeId must be an int")
			return
		}
		meme, err := memesModel.GetByID(params.id)
		if err != nil {
			respondWithErr(w, 500, "error getting meme by id")
			return
		}
		components.MemeModal(meme, params.idx[0]).Render(r.Context(), w)
	}
}
