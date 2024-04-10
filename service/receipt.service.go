package service

import (
	"errors"
	"strings"

	"github.com/crowmw/risiti/model"
	"gorm.io/gorm"
)

const (
	YYYYMMDD = "2006-01-02"
)

type ReceiptService interface {
	Create(receipt model.Receipt) (*model.Receipt, error)
	GetAll() (*[]model.Receipt, error)
	GetByText(text string) (*[]model.Receipt, error)
	GetByName(name string) (*model.Receipt, error)
}

type receiptService struct {
	db *gorm.DB
}

func DefaultReceiptService(db *gorm.DB) ReceiptService {
	return &receiptService{
		db,
	}
}

func (s *receiptService) Create(receipt model.Receipt) (*model.Receipt, error) {
	result := s.db.Create(&receipt)
	if result.Error != nil {
		return &model.Receipt{}, result.Error
	}
	return &receipt, nil
}

func (s *receiptService) GetAll() (*[]model.Receipt, error) {
	receipts := []model.Receipt{}
	result := s.db.Find(&receipts)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return &receipts, result.Error
	}
	return &receipts, nil
}

func (s *receiptService) GetByName(name string) (*model.Receipt, error) {
	receipt := model.Receipt{}

	result := s.db.Find(&receipt, "name = ?", name)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return &model.Receipt{}, result.Error
	}

	return &receipt, nil
}

func (s *receiptService) GetByText(text string) (*[]model.Receipt, error) {
	receipts := []model.Receipt{}
	searchQuery := strings.TrimSpace(text)
	result := s.db.Find(&receipts, "name like ? OR description like ? OR filename like ? OR date like ?", searchQuery)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return &receipts, result.Error
	}

	return &receipts, nil
}
