package handler

// import (
// 	"net/http"

// 	"github.com/gleb-korostelev/short-url.git/internal/handler/business"
// )

// func HandleRequest(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost && r.URL.Path == "/" {
// 		business.PostShorter(w, r)
// 	} else if r.Method == http.MethodGet {
// 		business.GetOriginal(w, r)
// 	} else {
// 		http.Error(w, "Not Found", http.StatusBadRequest)
// 	}
// }
