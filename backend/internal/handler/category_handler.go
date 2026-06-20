package handler

import (
	"net/http"

	"kafe-omborxona/internal/domain"
	"kafe-omborxona/internal/repository"
)

type CategoryHandler struct {
	repo *repository.CategoryRepo
}

func NewCategoryHandler(repo *repository.CategoryRepo) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	cats, err := h.repo.GetAll()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, cats)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCategoryRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}
	if req.Name == "" {
		Error(w, http.StatusBadRequest, "kategoriya nomi kiritilishi shart")
		return
	}
	cat := &domain.Category{Name: req.Name}
	if err := h.repo.Create(cat); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, cat)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri ID")
		return
	}
	var req domain.CreateCategoryRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}
	cat := &domain.Category{ID: id, Name: req.Name}
	if err := h.repo.Update(cat); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, cat)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri ID")
		return
	}
	if err := h.repo.Delete(id); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"message": "o'chirildi"})
}
