package handlers

import (
	"net/http"

	"github.com/crowmw/risiti/internal/components"
	receiptrepo "github.com/crowmw/risiti/internal/repo"
)

type GetReceiptsHandler struct {
	repo receiptrepo.IReceiptRepo
}

func NewGetReceiptsHandler(r receiptrepo.IReceiptRepo) *GetReceiptsHandler {
	return &GetReceiptsHandler{
		repo: r,
	}
}

func (h *GetReceiptsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	receipts, err := h.repo.GetAll()
	if err != nil {
		OnError(w, err, "Internal Server Error", http.StatusInternalServerError)
	}
	RenderView(w, r, components.ReceiptsList(receipts), "/")
}
