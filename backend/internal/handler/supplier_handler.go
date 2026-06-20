package handler

import (
	"net/http"

	"kafe-omborxona/internal/domain"
	"kafe-omborxona/internal/repository"
)

type SupplierHandler struct {
	repo *repository.SupplierRepo
}

func NewSupplierHandler(repo *repository.SupplierRepo) *SupplierHandler {
	return &SupplierHandler{repo: repo}
}

func (h *SupplierHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	sups, err := h.repo.GetAll()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, sups)
}

func (h *SupplierHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateSupplierRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}
	if req.Name == "" {
		Error(w, http.StatusBadRequest, "ta'minotchi nomi kiritilishi shart")
		return
	}
	sup := &domain.Supplier{Name: req.Name, Phone: req.Phone, Address: req.Address}
	if err := h.repo.Create(sup); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, sup)
}

func (h *SupplierHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri ID")
		return
	}
	var req domain.CreateSupplierRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}
	sup := &domain.Supplier{ID: id, Name: req.Name, Phone: req.Phone, Address: req.Address}
	if err := h.repo.Update(sup); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, sup)
}

func (h *SupplierHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
