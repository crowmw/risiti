package receiptrepo

import (
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	YYYYMMDD = "2006-01-02"
)

type Receipt struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Filename    string    `json:"filename"`
	Date        time.Time `json:"date"`
	Description string    `json:"string"`
}

type IReceiptRepo interface {
	Add(r Receipt) error
	GetAll() ([]Receipt, error)
}

type ReceiptRepo struct {
	receipt Receipt
	db      *sql.DB
}

func NewReceiptRepo(r Receipt, db *sql.DB) *ReceiptRepo {
	return &ReceiptRepo{
		receipt: r,
		db:      db,
	}
}

func (rr *ReceiptRepo) Add(receipt Receipt) error {
	stmt := `INSERT INTO receipt(name, description, date, filename) VALUES($1, $2, $3, $4)`

	_, err := rr.db.Exec(stmt, receipt.Name, receipt.Description, receipt.Date.Format(YYYYMMDD), receipt.Filename)
	if err != nil {
		slog.Error("Cannot add new receipt to store", err)
		return err
	}
	return nil
}

func (rr *ReceiptRepo) GetAll() ([]Receipt, error) {
	query := `SELECT r.id, r.name, r.description, r.date, r.filename FROM receipt r`

	result, err := rr.db.Query(query)
	if err != nil {
		slog.Error("Cannot get all receipts from store", err)
		return nil, err
	}

	defer result.Close()

	receipts := []Receipt{}
	for result.Next() {
		result.Scan(&rr.receipt.ID, &rr.receipt.Name, &rr.receipt.Description, &rr.receipt.Filename, &rr.receipt.Date)
		slog.Info(rr.receipt.Date.String())
		receipts = append(receipts, rr.receipt)
	}

	return receipts, nil
}
