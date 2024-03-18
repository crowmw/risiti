package service

import (
	"database/sql"
	"log/slog"

	"github.com/crowmw/risiti/model"
)

const (
	YYYYMMDD = "2006-01-02"
)

type IReceiptService interface {
	Create(receipt model.Receipt) (model.Receipt, error)
	ReadAll() ([]model.Receipt, error)
	ReadByText(text string) ([]model.Receipt, error)
}

type ReceiptService struct {
	db *sql.DB
}

func NewReceiptService(db *sql.DB) *ReceiptService {
	return &ReceiptService{
		db: db,
	}
}

func (s *ReceiptService) Create(receipt model.Receipt) (model.Receipt, error) {
	//TODO SANITIZE AND VALIDATE
	stmt := `INSERT INTO receipt(name, description, filename, date) VALUES($1, $2, $3, $4)`

	_, err := s.db.Exec(stmt, receipt.Name, receipt.Description, receipt.Filename, receipt.Date.Format(YYYYMMDD))
	if err != nil {
		slog.Error("Cannot add new receipt to store", err)
		return model.Receipt{}, err
	}

	// TODO PARSE RESULT
	return model.Receipt{}, nil
}

func (s *ReceiptService) ReadAll() ([]model.Receipt, error) {
	query := `SELECT id, name, description, filename, date FROM receipt`

	result, err := s.db.Query(query)
	if err != nil {
		slog.Error("Cannot get all receipts from store", err)
		return nil, err
	}

	defer result.Close()

	receipts := []model.Receipt{}

	for result.Next() {
		receipt := model.Receipt{}
		result.Scan(&receipt.ID, &receipt.Name, &receipt.Description, &receipt.Filename, &receipt.Date)
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

func (rr *ReceiptService) ReadByText(text string) ([]model.Receipt, error) {
	query := "SELECT id, name, description, filename, date FROM receipt WHERE name like '%" + text + "?%' OR description like '%" + text + "%' OR filename like '%" + text + "%' OR date like '%" + text + "%'"

	result, err := rr.db.Query(query)
	if err != nil {
		slog.Error("Cannot get receipts matching "+text+" from store", err)
		return nil, err
	}

	defer result.Close()

	receipts := []model.Receipt{}
	for result.Next() {
		receipt := model.Receipt{}
		result.Scan(&receipt.ID, &receipt.Name, &receipt.Description, &receipt.Filename, &receipt.Date)
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}
