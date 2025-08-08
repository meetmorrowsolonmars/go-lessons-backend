package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain"
)

//go:generate minimock -i UserStore -g

type UserStore interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

type UserController struct {
	userStore UserStore
}

func NewUserController(userStore UserStore) *UserController {
	return &UserController{
		userStore: userStore,
	}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type CreateUserResponse struct {
	ID int64 `json:"id"`
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	req, err := DecodeRequest[CreateUserRequest](r)
	if err != nil {
		slog.Error("Decode create user request", slog.String("error", err.Error()))

		err = WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body: %s", err)
		if err != nil {
			slog.Error("Write error response", slog.String("error", err.Error()))
		}

		return
	}

	user, err := c.userStore.Create(r.Context(), domain.User{
		Email:    req.Email,
		FullName: req.FullName,
	})
	if err != nil {
		slog.Error("Create user", slog.String("error", err.Error()))

		err = WriteErrorResponse(w, http.StatusInternalServerError, "Create user error")
		if err != nil {
			slog.Error("Write error response", slog.String("error", err.Error()))
		}

		return
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(CreateUserResponse{
		ID: user.ID,
	})
	if err != nil {
		slog.Error("Write create user response", slog.String("error", err.Error()))
	}
}

type User struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	FullName   string    `json:"full_name"`
	CreateTime time.Time `json:"create_time"`
}

type GetByEmailResponse struct {
	User
}

func (c *UserController) GetByEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	const pathKey = "email"

	email := r.PathValue(pathKey)
	if email == "" {
		slog.Error("Get user by email", slog.String("error", "email is empty"))

		err := WriteErrorResponse(w, http.StatusBadRequest, "Email is required")
		if err != nil {
			slog.Error("Write error response", slog.String("error", err.Error()))
		}

		return
	}

	user, err := c.userStore.GetByEmail(r.Context(), email)
	if err != nil {
		slog.Error("Get user by email", slog.String("error", err.Error()))

		err = WriteErrorResponse(w, http.StatusInternalServerError, "Get user by email error")
		if err != nil {
			slog.Error("Write error response", slog.String("error", err.Error()))
		}

		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(GetByEmailResponse{
		User: User{
			ID:         user.ID,
			Email:      user.Email,
			FullName:   user.FullName,
			CreateTime: user.CreateTime,
		},
	})
	if err != nil {
		slog.Error("Write get user by email response", slog.String("error", err.Error()))
	}
}

func (c *UserController) Register(mux *http.ServeMux) {
	// TODO: Show different ways to design API (REST, RPC).
	mux.HandleFunc("POST /v1/users", c.CreateUser)
	mux.HandleFunc("GET /v1/users/{email}", c.GetByEmail)
}
