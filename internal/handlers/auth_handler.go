package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"todo-list/internal/models"
	"todo-list/internal/repos"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo       *repos.UserRepository
	sessionManager *scs.SessionManager
}

func NewAuthHandler(userRepo *repos.UserRepository, sessionManager *scs.SessionManager) *AuthHandler {
	return &AuthHandler{userRepo, sessionManager}
}

// Register creates a new user
func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)
		user.Username = string(creds.Username)

		// Save user to the database
		if err := h.userRepo.CreateUser(&user); err != nil {
			log.Println("Error:", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "User registered successfully!")
	}
}

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintln(w, "Home page of the server!")
		w.Write([]byte("Home page of the server"))
	}
}

// Login authenticates a user and starts a session
func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := h.userRepo.GetUser(creds.Username, &user); err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Verify password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		fmt.Println("actual user id in db is value from db after login:::: ", user.ID)

		// Start session
		h.sessionManager.Put(r.Context(), "username", user.Username)
		h.sessionManager.Put(r.Context(), "userID", user.ID)

		// sessionUsername := sessionManager.GetString(r.Context(), "username")
		// sessionUserID := sessionManager.Get(r.Context(), "userID")
		// fmt.Println("sessionUsername value from db after login:::: ", sessionUsername)
		// fmt.Println("sessionUserID value from db after login:::: ", sessionUserID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Login successful!")
	}
}

// Logout ends the user's session sessionManager *scs.SessionManager
func (h *AuthHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.sessionManager.Destroy(r.Context())
		// session := sessionManager.GetString(r.Context(), "username")
		// fmt.Println("session value from db after Logout and sessionManager.Destroy: ", session)
		if err != nil {
			http.Error(w, "Failed to log out", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Logged out successfully!")
	}
}
