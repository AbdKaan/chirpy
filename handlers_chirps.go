package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/AbdKaan/chirpy/internal/auth"
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

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get token from header", nil)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate user ID", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	post, err := cfg.db.CreatePost(r.Context(), database.CreatePostParams{
		Body:   params.Body,
		UserID: userId,
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

func (cfg *apiConfig) handlerDeletePost(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	post, err := cfg.db.GetPost(r.Context(), chirpID)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if post.UserID != userId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = cfg.db.DeletePost(r.Context(), chirpID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
