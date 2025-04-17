package router

import (
	"github.com/guruorgoru/newsguru/pkg/handler"
	"github.com/guruorgoru/newsguru/pkg/models"
	"github.com/justinas/alice"
	"net/http"
)

func NewsRouter(app *models.NewsModel) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/", handler.RootHandler())
	router.HandleFunc("GET /news", handler.GetNewsHandler(app))
	router.HandleFunc("GET /news/{id}", handler.GetNewsByIdHandler(app))
	router.HandleFunc("POST /news", handler.PostNewsHandler(app))
	router.HandleFunc("DELETE /news/{id}", handler.DeleteNewsHandler(app))

	standard := alice.New(recoverPanic, logRequest, secureHeaders, enableCORS)
	return standard.Then(router)

}
