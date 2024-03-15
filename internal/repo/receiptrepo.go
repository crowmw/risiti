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
	Description string    `json:"filename"`
	Filename    string    `json:"date"`
	Date        time.Time `json:"string"`
}

type IReceiptRepo interface {
	Add(r Receipt) error
	GetAll() ([]Receipt, error)
	GetByText(text string) ([]Receipt, error)
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
	stmt := `INSERT INTO receipt(name, description, filename, date) VALUES($1, $2, $3, $4)`

	_, err := rr.db.Exec(stmt, receipt.Name, receipt.Description, receipt.Filename, receipt.Date.Format(YYYYMMDD))
	if err != nil {
		slog.Error("Cannot add new receipt to store", err)
		return err
	}
	return nil
}

func (rr *ReceiptRepo) GetAll() ([]Receipt, error) {
	query := `SELECT id, name, description, filename, date FROM receipt`

	result, err := rr.db.Query(query)
	if err != nil {
		slog.Error("Cannot get all receipts from store", err)
		return nil, err
	}

	defer result.Close()

	receipts := []Receipt{}
	for result.Next() {
		result.Scan(&rr.receipt.ID, &rr.receipt.Name, &rr.receipt.Description, &rr.receipt.Filename, &rr.receipt.Date)
		receipts = append(receipts, rr.receipt)
	}

	return receipts, nil
}

func (rr *ReceiptRepo) GetByText(text string) ([]Receipt, error) {
	query := "SELECT id, name, description, filename, date FROM receipt WHERE name like '%" + text + "?%' OR description like '%" + text + "%' OR filename like '%" + text + "%' OR date like '%" + text + "%'"

	result, err := rr.db.Query(query)
	if err != nil {
		slog.Error("Cannot get receipts matching "+text+" from store", err)
		return nil, err
	}

	defer result.Close()

	receipts := []Receipt{}
	for result.Next() {
		result.Scan(&rr.receipt.ID, &rr.receipt.Name, &rr.receipt.Description, &rr.receipt.Filename, &rr.receipt.Date)
		receipts = append(receipts, rr.receipt)
	}

	return receipts, nil
}
