package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"tenant/internal/model"
	"tenant/internal/repository"
	"tenant/internal/service/messaging"
	"tenant/pkg/derrors"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
)

type TenantUsecase interface {
	CreateTenant(ctx context.Context, name string) (tenant *model.Tenant, err error)
	DeleteTenant(ctx context.Context, clientID string) error
	ProcessPayload(ctx context.Context, clientID string, payload interface{}) error
	GetTenant(ctx context.Context, clientID string) (tenant *model.Tenant, err error)
}

type tenantUsecase struct {
	repo repository.TenantRepository
	mq   messaging.Messagging
}

// NewTenantUsecase initializes a new tenant usecase
func NewTenantUsecase(repo repository.TenantRepository, mq messaging.Messagging) TenantUsecase {
	return &tenantUsecase{repo: repo, mq: mq}
}

// CreateTenant creates a new tenant and its associated RabbitMQ queue
func (s *tenantUsecase) CreateTenant(ctx context.Context, name string) (tenant *model.Tenant, err error) {
	defer derrors.Wrap(&err, "CreateTenant(%q)", name)

	// Generate a unique Client ID
	clientID := ksuid.New().String()

	// Prepare tenant model
	tenant = &model.Tenant{
		ClientID: clientID,
		Name:     name,
	}

	// Save tenant to the repository
	err = s.repo.CreateTenant(ctx, tenant)
	if err != nil {
		return nil, err
	}

	// Define queue name
	queueName := fmt.Sprintf("%s.process", clientID)

	// Create RabbitMQ queue for the tenant
	err = s.mq.CreateQueue(queueName)
	if err != nil {
		// Rollback tenant creation if queue creation fails
		err = s.repo.SoftDeleteTenant(ctx, clientID)
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	// Start RabbitMQ consumer for the queue
	err = s.mq.StartQueue(ctx, queueName, func(message string) {
		// Log or process the message here
		logrus.Infof("Processing message for tenant %s: %s\n", clientID, message)
	})
	if err != nil {
		// Rollback tenant creation if starting consumer fails
		err = s.repo.SoftDeleteTenant(ctx, clientID)
		if err != nil {
			return nil, err
		}

		err = s.mq.DeleteQueue(queueName)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	return tenant, nil
}

// DeleteTenant deletes a tenant and its associated RabbitMQ queue
func (s *tenantUsecase) DeleteTenant(ctx context.Context, clientID string) (err error) {
	defer derrors.Wrap(&err, "DeleteTenant(%q)", clientID)

	// Soft delete tenant
	err = s.repo.SoftDeleteTenant(ctx, clientID)
	if err != nil {
		return err
	}

	// Define queue name
	queueName := fmt.Sprintf("%s.process", clientID)

	// Delete RabbitMQ queue
	err = s.mq.DeleteQueue(queueName)
	if err != nil {
		return err
	}

	return
}

// ProcessPayload publishes a payload to the RabbitMQ queue of a specific tenant
func (s *tenantUsecase) ProcessPayload(ctx context.Context, clientID string, payload interface{}) (err error) {
	defer derrors.Wrap(&err, "ProcessPayload(%q)", clientID)

	// Validate tenant existence
	tenant, err := s.repo.GetTenantByClientID(ctx, clientID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if tenant == nil {
		return derrors.New(derrors.NotFound, "tenant not found")
	}

	// Define queue name
	queueName := fmt.Sprintf("%s.process", clientID)

	// Publish payload to the queue
	err = s.mq.Publish(ctx, queueName, payload)
	if err != nil {
		return derrors.New(derrors.Unknown, "failed to publish payload to queue %s for tenant %s: %v", queueName, tenant.Name, err)
	}

	logrus.Infof("Payload successfully published to queue %s for tenant %s\n", queueName, tenant.Name)
	return nil
}

// GetTenant retrieves a tenant by Client ID
func (s *tenantUsecase) GetTenant(ctx context.Context, clientID string) (tenant *model.Tenant, err error) {
	tenant, err = s.repo.GetTenantByClientID(ctx, clientID)
	if err != nil {
		return
	}
	return
}
