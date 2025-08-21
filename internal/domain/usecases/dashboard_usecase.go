package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/sirupsen/logrus"
)

type DashboardUseCaseInterface interface {
	GetMonthlySumByUser(ctx context.Context, id uuid.UUID) (interface{}, error)
}

type DashboardUseCase struct {
	dashboardRepo repositories.DashboardRepository
}

func NewDashboardUseCase(dashboardRepo repositories.DashboardRepository) DashboardUseCaseInterface {
	return &DashboardUseCase{
		dashboardRepo: dashboardRepo,
	}
}

func (uc *DashboardUseCase) GetMonthlySumByUser(ctx context.Context, id uuid.UUID) (interface{}, error) {
	funcCtx := "GetMonthlySumByUser"

	dashboardData, err := uc.dashboardRepo.GetMonthlySumByUser(ctx, id)
	if err != nil {
		logger.LogError(funcCtx, "failed to get dashboard data", err, logrus.Fields{
			"dashboard_id": id.String(),
		})
		return nil, helpers.NewNotFoundError("dashboard data not found", "")
	}

	return dashboardData, nil
}
