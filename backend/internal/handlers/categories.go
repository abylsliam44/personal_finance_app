package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"finance_project/internal/models"
	"finance_project/internal/services"
)

type CategoryHandler struct {
	Service *services.CategoryService
}

// NewCategoryHandler создает новый обработчик для категорий.
func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{Service: service}
}

// GetAllCategoriesHandler возвращает список всех категорий.
// @Summary Список категорий
// @Description Возвращает список всех категорий
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {string} string "Failed to retrieve categories"
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Service.GetAllCategories()
	if err != nil {
		http.Error(w, "Failed to retrieve categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// CreateCategoryHandler создает новую категорию.
// @Summary Создание категории
// @Description Создает новую категорию
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category body"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create category"
// @Router /categories/create [post]
func (h *CategoryHandler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateCategory(category); err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetCategoryByIDHandler возвращает категорию по ID.
// @Summary Получение категории по ID
// @Description Возвращает данные категории по её ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id query int true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {string} string "Invalid category ID"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Failed to retrieve category"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.Service.GetCategoryByID(id)
	if err != nil {
		http.Error(w, "Failed to retrieve category", http.StatusInternalServerError)
		return
	}
	if category == nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// UpdateCategoryHandler обновляет категорию.
// @Summary Обновление категории
// @Description Обновляет данные категории
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category body"
// @Success 200 {string} string "Updated"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to update category"
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateCategory(category); err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteCategoryHandler удаляет категорию.
// @Summary Удаление категории
// @Description Удаляет категорию по её ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id query int true "Category ID"
// @Success 200 {string} string "Deleted"
// @Failure 400 {string} string "Invalid category ID"
// @Failure 500 {string} string "Failed to delete category"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteCategory(id); err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
