package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/joaquinamado/gobank/docs"
	env "github.com/joaquinamado/gobank/internal/app/env"
	"github.com/joaquinamado/gobank/internal/app/repositories"
	types "github.com/joaquinamado/gobank/internal/app/types"
	"github.com/swaggo/http-swagger"
)

var jwtSecret = env.GetString("JW_TOKEN", "")

type APIServer struct {
	listenAddr string
	repo       repositories.Repositories
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func NewApiServer(listenAddr string, repo repositories.Repositories) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		repo:       repo,
	}
}

func (s *APIServer) Mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {

		// === Docs ===
		r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/v1/docs/index.html", http.StatusFound)
		})
		r.Get("/docs/*", httpSwagger.Handler(
			httpSwagger.URL("docs/doc.json"),
		))

		// === Health ===
		r.Get("/health", makeHttpHandleFunc(s.handleHealth))

		// === Auth ===
		r.Post("/login", makeHttpHandleFunc(s.handleLogin))

		// === Account ===
		r.Route("/account", func(r chi.Router) {
			r.Get("/", makeHttpHandleFunc(s.handleGetAccount))
			r.Post("/", makeHttpHandleFunc(s.handleCreateAccount))
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", withJWTAuth(makeHttpHandleFunc(s.handleGetAccountById), s.repo))
				r.Delete("/", withJWTAuth(makeHttpHandleFunc(s.handleDeleteAccount), s.repo))
			})
		})

		// === Transfer ===
		r.Route("/transfer", func(r chi.Router) {
			r.Post("/", withJWTAuth(makeHttpHandleFunc(s.handleTransfer), s.repo))
		})
	})

	return r
}

func (s *APIServer) Run(mux http.Handler) {
	host := env.GetString("API_HOST", "localhost")
	docs.SwaggerInfo.Version = env.GetString("API_VERSION", "1.0")
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", host, s.listenAddr)
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.listenAddr),
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Println("JSON API server running on port: ", s.listenAddr)

	srv.ListenAndServe()
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getId(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id %s", idStr)
	}
	return id, nil
}

func withJWTAuth(handlerFunc http.HandlerFunc, s repositories.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, error := validateJWT(tokenString)

		if error != nil || !token.Valid {
			permmisionDenied(w)
			return
		}

		idStr, err := getId(r)
		if err != nil {
			permmisionDenied(w)
			return
		}

		account, err := s.Account.GetAccountByID(idStr)
		if err != nil {
			permmisionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if account.Number != int64(claims["accountNumber"].(float64)) {
			permmisionDenied(w)
			return
		}
		handlerFunc(w, r)
	}
}

func permmisionDenied(w http.ResponseWriter) {
	WriteJson(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
}

func createJWT(account *types.Account) (string, error) {

	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
