package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain"
)

func TestUserController_CreateUser(t *testing.T) {
	t.Run("Should create a new user", func(t *testing.T) {
		mc := minimock.NewController(t)

		// Prepare.
		userStore := NewUserStoreMock(mc)

		user := domain.User{
			ID:         200_001,
			Email:      "user_1@example.com",
			FullName:   "Ivan Ivanov",
			CreateTime: time.Date(2025, time.November, 10, 23, 0, 0, 0, time.UTC),
		}

		userStore.CreateMock.
			When(minimock.AnyContext, domain.User{
				Email:    "user_1@example.com",
				FullName: "Ivan Ivanov",
			}).
			Then(user, nil)

		// Act.
		controller := UserController{
			userStore: userStore,
		}

		server := httptest.NewServer(http.HandlerFunc(controller.CreateUser))
		defer server.Close()

		request, _ := json.Marshal(CreateUserRequest{
			Email:    "user_1@example.com",
			FullName: "Ivan Ivanov",
		})

		response, err := server.Client().Post(server.URL, "application/json", bytes.NewReader(request))
		defer func() {
			_ = response.Body.Close()
		}()

		// Check.
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.StatusCode)

		payload := CreateUserResponse{}

		err = json.NewDecoder(response.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, payload.ID)
	})

	t.Run("Should fail to create a new user", func(t *testing.T) {
		mc := minimock.NewController(t)

		// Prepare.
		userStore := NewUserStoreMock(mc)

		// Act.
		controller := UserController{
			userStore: userStore,
		}

		server := httptest.NewServer(http.HandlerFunc(controller.CreateUser))
		defer server.Close()

		request := []byte(`<user><email>user_1@example.com</email></user>`)

		response, err := server.Client().Post(server.URL, "application/json", bytes.NewReader(request))
		defer func() {
			_ = response.Body.Close()
		}()

		// Check.
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)

		payload := DefaultError{}

		err = json.NewDecoder(response.Body).Decode(&payload)
		assert.NoError(t, err)
		assert.NotEmpty(t, payload.Message)
	})
}
