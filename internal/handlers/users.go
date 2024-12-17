package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"finance_project/internal/models"
	"finance_project/internal/services"
)

// UserHandler представляет обработчики пользователей.
type UserHandler struct {
	Service *services.UserService
}

// NewUserHandler создает новый обработчик пользователей.
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// CreateUserHandler добавляет нового пользователя.
// @Summary Регистрация пользователя
// @Description Создает нового пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User body"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create user"
// @Router /users/create [post]
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateUser(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetAllUsersHandler возвращает список всех пользователей.
// @Summary Список пользователей
// @Description Возвращает список всех зарегистрированных пользователей
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByIDHandler возвращает пользователя по ID.
// @Summary Получение пользователя по ID
// @Description Возвращает данные конкретного пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {string} string "User not found"
// @Router /users/get [get]
func (h *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUserHandler обновляет информацию о пользователе.
// @Summary Обновление пользователя
// @Description Обновляет данные пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User body"
// @Success 200 {string} string "Updated"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to update user"
// @Router /users/update [put]
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateUser(user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteUserHandler удаляет пользователя.
// @Summary Удаление пользователя
// @Description Удаляет пользователя по указанному ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Success 200 {string} string "Deleted"
// @Failure 400 {string} string "Invalid user ID"
// @Failure 500 {string} string "Failed to delete user"
// @Router /users/delete [delete]
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteUser(id); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
