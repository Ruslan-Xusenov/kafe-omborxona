package handler

import (
	"net/http"

	"kafe-omborxona/internal/domain"
	"kafe-omborxona/internal/repository"
)

type RecipeHandler struct {
	repo *repository.RecipeRepo
}

func NewRecipeHandler(repo *repository.RecipeRepo) *RecipeHandler {
	return &RecipeHandler{repo: repo}
}

func (h *RecipeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.repo.GetAll()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, recipes)
}

func (h *RecipeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateRecipeRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}
	if req.ProductID == 0 || req.Name == "" {
		Error(w, http.StatusBadRequest, "mahsulot va retsept nomi kiritilishi shart")
		return
	}

	rec := &domain.Recipe{
		ProductID:   req.ProductID,
		Name:        req.Name,
		Ingredients: req.Ingredients,
	}

	if err := h.repo.Create(rec); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, rec)
}

func (h *RecipeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri ID")
		return
	}

	var req domain.CreateRecipeRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}

	rec := &domain.Recipe{
		ID:          id,
		ProductID:   req.ProductID,
		Name:        req.Name,
		Ingredients: req.Ingredients,
	}

	if err := h.repo.Update(rec); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, rec)
}

func (h *RecipeHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
