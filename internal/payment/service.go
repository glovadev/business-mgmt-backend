package payment

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"business-mgmt-backend/internal/models"
	"github.com/google/uuid"
)

var (
	ErrPaymentNotFound    = errors.New("payment not found")
	ErrInvalidInput       = errors.New("invalid input data")
	ErrInvalidAllocation  = errors.New("allocation exceeds payment amount or invoice outstanding amount")
	ErrInvoiceNotFound    = errors.New("invoice not found for allocation")
	ErrPaymentNotPosted   = errors.New("payment is not in POSTED state")
)

type Service interface {
	CollectCustomerPayment(ctx context.Context, req models.CreatePaymentRequest, idempotencyKey string, currentUserID uuid.UUID) (*models.PaymentModel, error)
	ListPayments(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListPaymentsResponse, error)
	GetPayment(ctx context.Context, id uuid.UUID) (*models.PaymentModel, error)
	ReversePayment(ctx context.Context, id uuid.UUID, reason, idempotencyKey string, currentUserID uuid.UUID) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func parseStrAmt(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func (s *serviceImpl) CollectCustomerPayment(ctx context.Context, req models.CreatePaymentRequest, idempotencyKey string, currentUserID uuid.UUID) (*models.PaymentModel, error) {
	if req.CustomerID == "" || req.Amount == "" {
		return nil, ErrInvalidInput
	}
	
	amount := parseStrAmt(req.Amount)
	if amount <= 0 {
		return nil, ErrInvalidInput
	}

	customerID, _ := uuid.Parse(req.CustomerID)
	var collectionStoreID *uuid.UUID
	if req.CollectionStoreID != nil && *req.CollectionStoreID != "" {
		if id, err := uuid.Parse(*req.CollectionStoreID); err == nil {
			collectionStoreID = &id
		}
	}

	payment := &models.PaymentModel{
		ID:                   uuid.New(),
		PaymentNumber:        fmt.Sprintf("PAY-%d", time.Now().Unix()),
		CustomerID:           customerID,
		Amount:               formatFloat(amount),
		PaymentMethod:        req.PaymentMethod,
		TransactionReference: req.TransactionReference,
		CollectionStoreID:    collectionStoreID,
		CollectedBy:          currentUserID,
		CollectedAt:          time.Now(),
		Status:               "POSTED",
		Notes:                req.Notes,
		CreatedAt:            time.Now(),
	}

	err := s.repo.ExecutePaymentTransaction(ctx, payment, req.AllocationStrategy, req.Allocations, func(ctx context.Context, locker InvoiceLocker) ([]models.PaymentAllocationModel, error) {
		allocations := make([]models.PaymentAllocationModel, 0)
		remainingAmount := amount

		if req.AllocationStrategy == "EXPLICIT" {
			for _, allocReq := range req.Allocations {
				invID, err := uuid.Parse(allocReq.InvoiceID)
				if err != nil { return nil, ErrInvalidInput }
				
				allocAmt := parseStrAmt(allocReq.Amount)
				if allocAmt <= 0 { continue }
				
				invoice, err := locker.GetInvoiceByID(ctx, invID)
				if err != nil || invoice == nil {
					return nil, ErrInvoiceNotFound
				}
				
				if invoice.CustomerID != customerID {
					return nil, errors.New("invoice does not belong to the specified customer")
				}

				outstanding := parseStrAmt(invoice.OutstandingAmount)
				if allocAmt > outstanding {
					return nil, ErrInvalidAllocation
				}
				
				if allocAmt > remainingAmount {
					return nil, ErrInvalidAllocation
				}

				allocations = append(allocations, models.PaymentAllocationModel{
					ID:        uuid.New(),
					PaymentID: payment.ID,
					InvoiceID: invID,
					Amount:    formatFloat(allocAmt),
					CreatedAt: time.Now(),
				})
				
				remainingAmount -= allocAmt
				
				// In a real app, update invoice.PaidAmount, OutstandingAmount, Status here.
				// This is represented by the atomic transaction wrapper in the database.
			}
		} else {
			// OLDEST_FIRST
			invoices, err := locker.GetOutstandingInvoicesByCustomer(ctx, customerID)
			if err != nil { return nil, err }
			
			for _, invoice := range invoices {
				if remainingAmount <= 0 { break }
				
				outstanding := parseStrAmt(invoice.OutstandingAmount)
				allocAmt := outstanding
				if remainingAmount < outstanding {
					allocAmt = remainingAmount
				}

				allocations = append(allocations, models.PaymentAllocationModel{
					ID:        uuid.New(),
					PaymentID: payment.ID,
					InvoiceID: invoice.ID,
					Amount:    formatFloat(allocAmt),
					CreatedAt: time.Now(),
				})
				
				remainingAmount -= allocAmt
			}
		}

		// It is acceptable for remainingAmount > 0 (creates an unallocated customer credit), 
		// but for this V1 model we'll assume it's stored on the payment level implicitly.
		
		return allocations, nil
	})

	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *serviceImpl) ListPayments(ctx context.Context, page, pageSize int, filters map[string]string) (*models.ListPaymentsResponse, error) {
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	payments, total, err := s.repo.ListPayments(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	return &models.ListPaymentsResponse{
		Payments:   payments,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *serviceImpl) GetPayment(ctx context.Context, id uuid.UUID) (*models.PaymentModel, error) {
	return s.repo.GetPaymentByID(ctx, id)
}

func (s *serviceImpl) ReversePayment(ctx context.Context, id uuid.UUID, reason, idempotencyKey string, currentUserID uuid.UUID) error {
	payment, err := s.repo.GetPaymentByID(ctx, id)
	if err != nil {
		return ErrPaymentNotFound
	}

	if payment.Status != "POSTED" {
		return ErrPaymentNotPosted
	}

	return s.repo.ExecutePaymentReversalTransaction(ctx, id, reason, currentUserID)
}
