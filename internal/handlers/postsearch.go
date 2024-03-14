package handlers

import (
	"net/http"

	"github.com/crowmw/risiti/internal/components"
	receiptrepo "github.com/crowmw/risiti/internal/repo"
)

type PostSearchHandler struct {
	repo receiptrepo.IReceiptRepo
}

func NewPostSearchHandler(r receiptrepo.IReceiptRepo) *PostSearchHandler {
	return &PostSearchHandler{
		repo: r,
	}
}

func (h *PostSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("search")

	receipts, err := h.repo.GetByText(text)
	if err != nil {
		OnError(w, err, "Internal Server Error", http.StatusInternalServerError)
	}
	RenderView(w, r, components.ReceiptsList(receipts), "/")
}
