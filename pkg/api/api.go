package api

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"gonews/pkg/database"
	"gonews/pkg/models"

	"github.com/gorilla/mux"
)

type API struct {
	router *mux.Router
	db     *database.DB
}

func NewAPI(db *database.DB) *API {
	api := &API{router: mux.NewRouter(), db: db}
	api.setupRoutes()
	return api
}

func (a *API) setupRoutes() {
	a.router.Use(RequestIDMiddleware)
	a.router.Use(LoggingMiddleware)

	a.router.HandleFunc("/news/{n}", a.getNews).Methods("GET", "OPTIONS")
	a.router.HandleFunc("/news/detail/{id}", a.getNewsDetail).Methods("GET")
	a.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (a *API) getNews(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.Atoi(mux.Vars(r)["n"])
	if err != nil || n <= 0 {
		http.Error(w, `{"error":"invalid n"}`, http.StatusBadRequest)
		return
	}
	if n > 100 {
		n = 100
	}

	search := r.URL.Query().Get("s")
	page := 1
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	const perPage = 15
	offset := (page - 1) * perPage

	total, _ := a.db.GetPostsCount(search)
	posts, err := a.db.GetPostsPaginated(search, perPage, offset)
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPages == 0 {
		totalPages = 1
	}

	resp := models.NewsResponse{
		Posts: posts,
		Pagination: models.Pagination{
			TotalPages:   totalPages,
			CurrentPage:  page,
			ItemsPerPage: perPage,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(resp)
}

func (a *API) getNewsDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	// GetPostByID возвращает ОДНУ структуру models.Post, а не слайс
	post, err := a.db.GetPostByID(id)
	if err != nil {
		// Если запись не найдена, sql.ErrNoRows будет возвращен драйвером
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(post) // Кодируем саму структуру, а не posts[0]
}

func (a *API) GetRouter() *mux.Router { return a.router }
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
