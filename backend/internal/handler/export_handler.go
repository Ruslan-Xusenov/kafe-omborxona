package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
	"kafe-omborxona/internal/domain"
	"kafe-omborxona/internal/repository"
)

type ExportHandler struct {
	txnRepo *repository.TransactionRepo
}

func NewExportHandler(txnRepo *repository.TransactionRepo) *ExportHandler {
	return &ExportHandler{txnRepo: txnRepo}
}

func (h *ExportHandler) ExportTransactions(w http.ResponseWriter, r *http.Request) {
	filter := domain.TransactionFilter{
		Type:     domain.TransactionType(r.URL.Query().Get("type")),
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
	}
	if pid := r.URL.Query().Get("product_id"); pid != "" {
		if id, err := strconv.Atoi(pid); err == nil {
			filter.ProductID = id
		}
	}

	txns, err := h.txnRepo.GetAll(filter)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	f := excelize.NewFile()
	defer f.Close()
	sheet := "Tranzaksiyalar"
	f.SetSheetName("Sheet1", sheet)

	// Set Headers
	headers := []string{"ID", "Sana", "Turi", "Mahsulot", "Ta'minotchi", "Miqdor", "Narx", "Jami Summa", "Izoh", "Xodim"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Style header
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4F46E5"}, Pattern: 1},
	})
	f.SetRowStyle(sheet, 1, 1, headerStyle)

	// Write data
	for i, t := range txns {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), t.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), t.CreatedAt.Format("2006-01-02 15:04"))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), string(t.Type))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), t.ProductName)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), t.SupplierName)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), t.Quantity)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), t.UnitPrice)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), t.TotalAmount)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), t.Note)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), t.UserName)
	}

	filename := fmt.Sprintf("Hisobot_%s.xlsx", time.Now().Format("2006-01-02"))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	if err := f.Write(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
