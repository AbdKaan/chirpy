package main

import (
	"encoding/json"
	"net/http"

	"github.com/AbdKaan/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body    string `json:"body"`
		User_ID string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	user_id, err := uuid.Parse(params.User_ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID", err)
		return
	}

	post, err := cfg.db.CreatePost(r.Context(), database.CreatePostParams{
		Body:   params.Body,
		UserID: user_id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create posts", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Post{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Body:      cencorProfane(post.Body),
		User_ID:   post.UserID.String(),
	})
}

func (cfg *apiConfig) handlerGetPost(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirp ID", err)
		return
	}

	post, err := cfg.db.GetPost(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get post", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Post{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Body:      cencorProfane(post.Body),
		User_ID:   post.UserID.String(),
	})
}

func (cfg *apiConfig) handlerGetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := cfg.db.GetPosts(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts", err)
		return
	}

	var postsArr []Post
	for _, post := range posts {
		postsArr = append(postsArr, Post{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			Body:      cencorProfane(post.Body),
			User_ID:   post.UserID.String(),
		})
	}

	respondWithJSON(w, http.StatusOK, postsArr)
}
