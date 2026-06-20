package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"kafe-omborxona/internal/domain"
	"kafe-omborxona/internal/middleware"
	"kafe-omborxona/internal/repository"
	"kafe-omborxona/internal/service"
)

type TransactionHandler struct {
	repo  *repository.TransactionRepo
	tgSvc *service.TelegramService
}

func NewTransactionHandler(repo *repository.TransactionRepo, tgSvc *service.TelegramService) *TransactionHandler {
	return &TransactionHandler{repo: repo, tgSvc: tgSvc}
}

func (h *TransactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	filter := domain.TransactionFilter{
		Type:      domain.TransactionType(r.URL.Query().Get("type")),
		DateFrom:  r.URL.Query().Get("date_from"),
		DateTo:    r.URL.Query().Get("date_to"),
	}
	if pid := r.URL.Query().Get("product_id"); pid != "" {
		if id, err := strconv.Atoi(pid); err == nil {
			filter.ProductID = id
		}
	}

	txns, err := h.repo.GetAll(filter)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, txns)
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateTransactionRequest
	if err := Decode(r, &req); err != nil {
		Error(w, http.StatusBadRequest, "noto'g'ri format")
		return
	}
	if req.ProductID == 0 || req.Quantity <= 0 {
		Error(w, http.StatusBadRequest, "mahsulot va miqdor kiritilishi shart")
		return
	}
	validTypes := map[domain.TransactionType]bool{
		domain.TransactionPurchase: true, domain.TransactionReturn: true,
		domain.TransactionSale: true, domain.TransactionWriteOff: true,
	}
	if !validTypes[req.Type] {
		Error(w, http.StatusBadRequest, "noto'g'ri tranzaksiya turi")
		return
	}

	t := &domain.Transaction{
		ProductID:  req.ProductID,
		SupplierID: req.SupplierID,
		UserID:     middleware.GetUserID(r),
		Type:       req.Type,
		Quantity:   req.Quantity,
		UnitPrice:  req.UnitPrice,
		Note:       req.Note,
		ExpiryDate: req.ExpiryDate,
	}
	if err := h.repo.Create(t); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Telegramga yuborish
	if h.tgSvc != nil {
		typeLabels := map[domain.TransactionType]string{
			domain.TransactionPurchase: "Kirim (Xarid)",
			domain.TransactionReturn:   "Qaytarish (Vozvrat)",
			domain.TransactionSale:     "Sotuv (Chiqim)",
			domain.TransactionWriteOff: "Spisaniye (Yo'qotish)",
		}
		
		// To get the product name
		if fullT, err := h.repo.GetByID(t.ID); err == nil {
			msg := fmt.Sprintf("🔔 <b>Yangi Tranzaksiya</b>\n\n"+
				"Turi: <b>%s</b>\n"+
				"Mahsulot: <b>%s</b>\n"+
				"Miqdor: <b>%g</b>\n"+
				"Narx: <b>%g</b>\n"+
				"Jami: <b>%g</b>",
				typeLabels[t.Type], fullT.ProductName, t.Quantity, t.UnitPrice, t.TotalAmount)
			
			if t.Note != "" {
				msg += fmt.Sprintf("\nIzoh: %s", t.Note)
			}
			h.tgSvc.SendAdminMessage(msg)
		}
	}

	JSON(w, http.StatusCreated, t)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
