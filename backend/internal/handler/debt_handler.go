package handler

import (
	"net/http"

	"kafe-omborxona/internal/domain"
	"kafe-omborxona/internal/repository"
)

type DebtHandler struct {
	repo *repository.DebtRepo
}

func NewDebtHandler(repo *repository.DebtRepo) *DebtHandler {
	return &DebtHandler{repo: repo}
}

func (h *DebtHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	debts, err := h.repo.GetAll()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, debts)
}

func (h *DebtHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateDebtRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}
	if req.SupplierID == 0 || req.TotalDebt <= 0 {
		Error(w, http.StatusBadRequest, "ta'minotchi va qarz miqdori kiritilishi shart")
		return
	}

	d := &domain.Debt{
		SupplierID:    req.SupplierID,
		TransactionID: req.TransactionID,
		TotalDebt:     req.TotalDebt,
		DueDate:       req.DueDate,
	}

	if err := h.repo.Create(d); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, d)
}

func (h *DebtHandler) Pay(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri ID")
		return
	}

	var req domain.PayDebtRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}
	if req.Amount <= 0 {
		Error(w, http.StatusBadRequest, "to'lov miqdori noldan katta bo'lishi kerak")
		return
	}

	if err := h.repo.Pay(id, req.Amount); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"message": "to'lov qabul qilindi"})
}

func (h *DebtHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
