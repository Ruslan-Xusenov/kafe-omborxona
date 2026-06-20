package handler

import (
	"net/http"

	"kafe-omborxona/internal/domain"
	"kafe-omborxona/internal/repository"
)

type ProductHandler struct {
	repo *repository.ProductRepo
}

func NewProductHandler(repo *repository.ProductRepo) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.GetWithStock()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, products)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateProductRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}
	if req.Name == "" {
		Error(w, http.StatusBadRequest, "mahsulot nomi kiritilishi shart")
		return
	}
	if req.Unit == "" {
		req.Unit = "dona"
	}
	p := &domain.Product{
		Name: req.Name, Unit: req.Unit, CategoryID: req.CategoryID,
		CostPrice: req.CostPrice, SalePrice: req.SalePrice,
		MinStock: req.MinStock, Barcode: req.Barcode,
	}
	if err := h.repo.Create(p); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, p)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri ID")
		return
	}
	var req domain.CreateProductRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}
	p := &domain.Product{
		ID: id, Name: req.Name, Unit: req.Unit, CategoryID: req.CategoryID,
		CostPrice: req.CostPrice, SalePrice: req.SalePrice,
		MinStock: req.MinStock, Barcode: req.Barcode,
	}
	if err := h.repo.Update(p); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, p)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
