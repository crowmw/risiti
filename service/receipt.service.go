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
	ReadByName(name string) (model.Receipt, error)
}

type ReceiptService struct {
	DB *sql.DB
}

func NewReceiptService(db *sql.DB) *ReceiptService {
	return &ReceiptService{
		DB: db,
	}
}

func (s *ReceiptService) Create(receipt model.Receipt) (model.Receipt, error) {
	stmt := `INSERT INTO receipt(name, description, filename, date) VALUES($1, $2, $3, $4)`

	_, err := s.DB.Exec(stmt, receipt.Name, receipt.Description, receipt.Filename, receipt.Date.Format(YYYYMMDD))
	if err != nil {
		slog.Error("Cannot add new receipt to store", err)
		return model.Receipt{}, err
	}

	// TODO PARSE RESULT
	return model.Receipt{}, nil
}

func (s *ReceiptService) ReadAll() ([]model.Receipt, error) {
	query := `SELECT id, name, description, filename, date FROM receipt`

	result, err := s.DB.Query(query)
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

func (s *ReceiptService) ReadByName(name string) (model.Receipt, error) {
	query := "SELECT id, name, description, filename, date FROM receipt WHERE name = '" + name + "'"

	receipt := model.Receipt{}
	err := s.DB.QueryRow(query).Scan(&receipt.ID, &receipt.Name, &receipt.Description, &receipt.Filename, &receipt.Date)
	if err != nil {
		slog.Error("Cannot get receipts matching "+name+" from store", err)
		return model.Receipt{}, err
	}

	return receipt, nil
}

func (s *ReceiptService) ReadByText(text string) ([]model.Receipt, error) {
	query := "SELECT id, name, description, filename, date FROM receipt WHERE name like '%" + text + "?%' OR description like '%" + text + "%' OR filename like '%" + text + "%' OR date like '%" + text + "%'"

	result, err := s.DB.Query(query)
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
