package httpserver

import (
	apiservice "memesearch/internal/api"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func GetRouter(api *apiservice.API) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	// router.Use(auth.GetAuthMiddleware(plog, cfg.SecretAuth))

	router.Get("/about", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("t.me/sprff_code")) })

	router.Post("/memes", PostMeme(api))
	router.Get("/memes/{id}", GetMemeByID(api))
	router.Put("/memes/{id}", PutMeme(api))
	router.Delete("/memes/{id}", DeleteMeme(api))

	router.Put("/media/{id}", PutMedia(api))
	router.Get("/media/{id}", GetMedia(api))

	// router.Post("/board")

	// router.Get("/search/")

	return router
}
