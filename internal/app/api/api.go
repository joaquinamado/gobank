package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	. "github.com/joaquinamado/gobank/internal/app/storage"
	. "github.com/joaquinamado/gobank/internal/app/types"
	. "github.com/joaquinamado/gobank/internal/app/utils"
)

var jwtSecret = GetEnvInstance().EnvVariables.JwtSecret

type APIServer struct {
	listenAddr string
	store      Storage
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func NewApiServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// === AUTH ===
	router.HandleFunc("/login", makeHttpHandleFunc(s.handleLogin))

	// === ACCOUNT ===
	router.HandleFunc("/account", makeHttpHandleFunc(s.handleAccount))
	router.HandleFunc(
		"/account/{id}",
		withJWTAuth(makeHttpHandleFunc(s.handleGetAccountById), s.store))

	// === TRANSFER ===
	router.HandleFunc("/transfer", makeHttpHandleFunc(s.handleTransfer))

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
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

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
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
		account, err := s.GetAccountByID(idStr)
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
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(jwtSecret), nil
	})
}

func createJWT(account *Account) (string, error) {

	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
