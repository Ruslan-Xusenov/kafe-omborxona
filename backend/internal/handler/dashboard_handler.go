package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"kafe-omborxona/internal/domain"
	"kafe-omborxona/internal/repository"
	"kafe-omborxona/internal/service"
)

type DashboardHandler struct {
	txnRepo  *repository.TransactionRepo
	debtRepo *repository.DebtRepo
	tgSvc    *service.TelegramService
}

func NewDashboardHandler(txnRepo *repository.TransactionRepo, debtRepo *repository.DebtRepo, tgSvc *service.TelegramService) *DashboardHandler {
	return &DashboardHandler{txnRepo: txnRepo, debtRepo: debtRepo, tgSvc: tgSvc}
}

func (h *DashboardHandler) Summary(w http.ResponseWriter, r *http.Request) {
	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	summary, err := h.txnRepo.GetSummary(dateFrom, dateTo)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, summary)
}

func (h *DashboardHandler) Inventory(w http.ResponseWriter, r *http.Request) {
	items, err := h.txnRepo.GetInventory()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, items)
}

func (h *DashboardHandler) Profit(w http.ResponseWriter, r *http.Request) {
	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	report, err := h.txnRepo.GetProfitReport(dateFrom, dateTo)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, report)
}

func (h *DashboardHandler) TopProducts(w http.ResponseWriter, r *http.Request) {
	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}

	tops, err := h.txnRepo.GetTopProducts(limit, dateFrom, dateTo)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, tops)
}

func (h *DashboardHandler) Alerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.txnRepo.GetAlerts()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, alerts)
}

func (h *DashboardHandler) TriggerReport(w http.ResponseWriter, r *http.Request) {
	if h.tgSvc == nil {
		JSON(w, http.StatusOK, map[string]string{"message": "telegram bot o'chirilgan"})
		return
	}

	dateStr := time.Now().Format("2006-01-02")
	summary, err := h.txnRepo.GetSummary(dateStr, dateStr)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	alerts, _ := h.txnRepo.GetAlerts()
	debts, _ := h.debtRepo.GetAll()

	var totalUnpaid float64
	for _, d := range debts {
		if d.Status != domain.DebtPaid {
			totalUnpaid += (d.TotalDebt - d.PaidAmount)
		}
	}

	msg := fmt.Sprintf("📊 <b>Kunlik Hisobot (%s)</b>\n\n", dateStr)
	msg += fmt.Sprintf("📈 Jami Sotuv: <b>%g so'm</b>\n", summary.TotalSales)
	msg += fmt.Sprintf("📉 Jami Kirim: <b>%g so'm</b>\n", summary.TotalPurchases)
	msg += fmt.Sprintf("💸 Qaytarilgan: <b>%g so'm</b>\n", summary.TotalReturns)
	msg += fmt.Sprintf("🗑 Yo'qotishlar: <b>%g so'm</b>\n\n", summary.TotalWriteOffs)
	
	msg += fmt.Sprintf("🤝 Jami to'lanmagan qarzlar: <b>%g so'm</b>\n\n", totalUnpaid)

	if len(alerts) > 0 {
		msg += "⚠️ <b>Ogohlantirishlar:</b>\n"
		for _, a := range alerts {
			if a.Type == "low_stock" {
				msg += fmt.Sprintf("➖ %s qoldig'i kam qoldi (%s)\n", a.ProductName, a.Value)
			} else {
				msg += fmt.Sprintf("⏳ %s muddati tugamoqda (%s)\n", a.ProductName, a.Value)
			}
		}
	} else {
		msg += "✅ Omborxonada muammo yo'q."
	}

	h.tgSvc.SendAdminMessage(msg)
	JSON(w, http.StatusOK, map[string]string{"message": "hisobot yuborildi"})
}
