package handlers

import (
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/usecases"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	ut "github.com/naufalfazanadi/finance-manager-go/pkg/utils"
	"github.com/naufalfazanadi/finance-manager-go/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	dashboardUseCase usecases.DashboardUseCaseInterface
	validator        *validator.Validator
}

func NewDashboardHandler(dashboardUseCase usecases.DashboardUseCaseInterface, validator *validator.Validator) *DashboardHandler {
	return &DashboardHandler{
		dashboardUseCase: dashboardUseCase,
		validator:        validator,
	}
}

func (h *DashboardHandler) GetMonthlySumByUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrIDRequired, ut.ErrIDRequired), ut.MsgErrIDRequired)
	}

	// Parse and validate UUID format
	userID, err := uuid.Parse(id)
	if err != nil {
		return helpers.HandleErrorResponse(c, helpers.NewBadRequestError(ut.MsgErrInvalidID, ut.ErrInvalidIDFormat), ut.MsgErrInvalidID)
	}

	if userID != c.Locals("userID").(uuid.UUID) && c.Locals("userRole") != "admin" {
		return helpers.HandleErrorResponse(c, helpers.NewForbiddenError("You do not have permission", "Permission denied"), "Permission denied")
	}

	dashboardData, err := h.dashboardUseCase.GetMonthlySumByUser(c.Context(), userID)
	if err != nil {
		return helpers.HandleErrorResponse(c, err, ut.FailedGetMsg("Dashboard"))
	}

	return helpers.SuccessResponse(c, ut.SuccessRetrieveMsg("Dashboard"), dashboardData)
}
