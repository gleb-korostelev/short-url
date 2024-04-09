package service

import "net/http"

type APIServiceI interface {
	GetOriginal(w http.ResponseWriter, r *http.Request)
	Ping(w http.ResponseWriter, r *http.Request)
	PostShorter(w http.ResponseWriter, r *http.Request)
	PostShorterJSON(w http.ResponseWriter, r *http.Request)
	ShortenBatchHandler(w http.ResponseWriter, r *http.Request)
	GetUserURLs(w http.ResponseWriter, r *http.Request)
	DeleteURLsHandler(w http.ResponseWriter, r *http.Request)
}
