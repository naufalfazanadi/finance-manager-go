package usecases

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/repositories"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/internal/infrastructure/auth"
	"github.com/naufalfazanadi/finance-manager-go/pkg/encryption"
	"github.com/naufalfazanadi/finance-manager-go/pkg/helpers"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/naufalfazanadi/finance-manager-go/pkg/mail"
	"github.com/sirupsen/logrus"
)

type AuthUseCaseInterface interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) error
	ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error
}

type AuthUseCase struct {
	userRepo repositories.UserRepository
}

func NewAuthUseCase(userRepo repositories.UserRepository) AuthUseCaseInterface {
	return &AuthUseCase{
		userRepo: userRepo,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	funcCtx := "Register"

	// Validate password strength
	if err := auth.ValidatePasswordStrength(req.Password); err != nil {
		logger.LogError(funcCtx, "password validation failed", err, logrus.Fields{"email": req.Email})
		return nil, helpers.NewBadRequestError("password validation failed", err.Error())
	}

	// Hash the email to check for existing user
	hashResult := encryption.HashSHA256(req.Email)
	if hashResult.Error != nil {
		logger.LogError(funcCtx, "failed to hash email", hashResult.Error, logrus.Fields{"email": req.Email})
		return nil, helpers.NewInternalError("failed to process email", hashResult.Error.Error())
	}

	emailHash := hashResult.Data.(string)

	// Check if user already exists using email hash
	existingUser, errUser := uc.userRepo.GetByEmailHash(ctx, emailHash)
	if existingUser != nil {
		logger.LogError(funcCtx, "user already exists during registration", nil, logrus.Fields{
			"email_hash": emailHash,
		})
		return nil, helpers.NewConflictError("user with this email already exists", "")
	}
	fmt.Println(errUser)

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.LogError(funcCtx, "failed to hash password", err, logrus.Fields{
			"email_hash": emailHash,
		})
		return nil, helpers.NewInternalError("failed to hash password", err.Error())
	}

	// Parse birth date from string to time.Time
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		logger.LogError(funcCtx, "failed to parse birth date", err, logrus.Fields{
			"birth_date": req.BirthDate,
		})
		return nil, helpers.NewBadRequestError("invalid birth date format", "birth date must be in YYYY-MM-DD format")
	}

	// Set role to default user role (no longer accepting from request)
	role := entities.UserRoleUser

	// Create user entity
	user := &entities.User{
		Email:     req.Email,  // Set the plain email - it will be encrypted in BeforeCreate hook
		BirthDate: &birthDate, // Set the parsed birth date - it will be encrypted in BeforeCreate hook
		Name:      req.Name,
		Password:  hashedPassword,
		Role:      role,
	}

	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		logger.LogError(funcCtx, "failed to create user", err, logrus.Fields{
			"email_hash": emailHash,
		})
		return nil, helpers.NewInternalError("failed to create user", err.Error())
	}

	// Generate JWT token (user.Email should be decrypted by AfterFind hook)
	token, err := auth.GenerateToken(user)
	if err != nil {
		logger.LogError(funcCtx, "failed to generate token", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return nil, helpers.NewInternalError("failed to generate token", err.Error())
	}

	return &dto.AuthResponse{
		UserResponse: *dto.MapToUserResponse(user),
		Token:        token,
	}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	funcCtx := "Login"

	// Hash the email to lookup user
	hashResult := encryption.HashSHA256(req.Email)
	if hashResult.Error != nil {
		logger.LogError(funcCtx, "failed to hash email", hashResult.Error, logrus.Fields{"email": req.Email})
		return nil, helpers.NewUnauthorizedError("invalid email or password", "")
	}

	emailHash := hashResult.Data.(string)

	// Get user by email hash
	user, err := uc.userRepo.GetByEmailHash(ctx, emailHash)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user by email hash", err, logrus.Fields{
			"email_hash": emailHash,
		})
		return nil, helpers.NewUnauthorizedError("invalid email or password", "")
	}

	// Check password
	if err := auth.CheckPassword(user.Password, req.Password); err != nil {
		logger.LogError(funcCtx, "login attempt with invalid password", nil, logrus.Fields{
			"email_hash": emailHash,
		})
		return nil, helpers.NewUnauthorizedError("invalid email or password", "")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user)
	if err != nil {
		logger.LogError(funcCtx, "failed to generate token", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return nil, helpers.NewInternalError("failed to generate token", err.Error())
	}

	return &dto.LoginResponse{
		UserResponse: *dto.MapToUserResponse(user),
		Token:        token,
	}, nil
}

func (uc *AuthUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	funcCtx := "GetProfile"

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user", err, logrus.Fields{
			"user_id": userID.String(),
		})
		return nil, helpers.NewNotFoundError("user not found", "")
	}

	return dto.MapToUserResponse(user), nil
}

func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) error {
	funcCtx := "ChangePassword"

	// Get user by ID
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.LogError(funcCtx, "failed to get user", err, logrus.Fields{
			"user_id": userID.String(),
		})
		return helpers.NewNotFoundError("user not found", "")
	}

	// Verify old password
	if err := auth.CheckPassword(user.Password, req.OldPassword); err != nil {
		logger.LogError(funcCtx, "invalid old password", nil, logrus.Fields{
			"user_id": userID.String(),
		})
		return helpers.NewBadRequestError("invalid old password", "")
	}

	// Validate new password strength
	if err := auth.ValidatePasswordStrength(req.NewPassword); err != nil {
		logger.LogError(funcCtx, "new password validation failed", err, logrus.Fields{
			"user_id": userID.String(),
		})
		return helpers.NewBadRequestError("new password validation failed", err.Error())
	}

	// Check if new password is different from old password
	if err := auth.CheckPassword(user.Password, req.NewPassword); err == nil {
		logger.LogError(funcCtx, "new password same as old password", nil, logrus.Fields{
			"user_id": userID.String(),
		})
		return helpers.NewBadRequestError("new password must be different from current password", "")
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		logger.LogError(funcCtx, "failed to hash new password", err, logrus.Fields{
			"user_id": userID.String(),
		})
		return helpers.NewInternalError("failed to hash new password", err.Error())
	}

	// Update password
	user.Password = hashedPassword
	if err := uc.userRepo.Update(ctx, user); err != nil {
		logger.LogError(funcCtx, "failed to update user password", err, logrus.Fields{
			"user_id": userID.String(),
		})
		return helpers.NewInternalError("failed to update password", err.Error())
	}

	logger.LogSuccess(funcCtx, "password changed successfully", logrus.Fields{
		"user_id": userID.String(),
	})

	return nil
}

func (uc *AuthUseCase) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error {
	funcCtx := "ForgotPassword"

	// Hash the email to lookup user
	hashResult := encryption.HashSHA256(req.Email)
	if hashResult.Error != nil {
		logger.LogError(funcCtx, "failed to hash email", hashResult.Error, logrus.Fields{"email": req.Email})
		return helpers.NewInternalError("failed to process email", hashResult.Error.Error())
	}

	emailHash := hashResult.Data.(string)

	// Get user by email hash
	user, err := uc.userRepo.GetByEmailHash(ctx, emailHash)
	if err != nil {
		logger.LogError(funcCtx, "user not found", err, logrus.Fields{
			"email_hash": emailHash,
		})
		return helpers.NewNotFoundError("user not found", "")
	}

	// Check if user already has a forgot password token and if it's still in cooldown period
	if user.ForgotPasswordToken != "" {
		_, timestamp, err := encryption.DecryptResetToken(user.ForgotPasswordToken)
		if err == nil {
			// Check cooldown (3 minutes = 180000 milliseconds)
			if err := encryption.CheckResetTokenCooldown(timestamp, 180000); err != nil {
				logger.LogError(funcCtx, "forgot password cooldown active", nil, logrus.Fields{
					"user_id": user.ID.String(),
				})
				return helpers.NewConflictError("password reset request too frequent", err.Error())
			}
		}
	}

	// Generate random string for token
	randomString := encryption.GenerateRandomString(32)

	// Encrypt the token with timestamp
	forgotPasswordToken, err := encryption.EncryptResetToken(randomString)
	if err != nil {
		logger.LogError(funcCtx, "failed to encrypt token", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return helpers.NewInternalError("failed to generate reset token", err.Error())
	}

	// Update user with forgot password token
	if err := uc.userRepo.UpdateForgotPasswordToken(ctx, user.ID, forgotPasswordToken); err != nil {
		logger.LogError(funcCtx, "failed to update forgot password token", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return helpers.NewInternalError("failed to update user", err.Error())
	}

	// Generate reset URL
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, forgotPasswordToken)

	// Prepare template data
	templateData := mail.EmailTemplateData{
		Name:     user.Name,
		ResetURL: resetURL,
	}

	// Render email template using LoadTemplate
	htmlBody, err := mail.LoadTemplate("forgot_password.html", templateData)
	if err != nil {
		logger.LogError(funcCtx, "failed to render email template", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return helpers.NewInternalError("failed to prepare email", err.Error())
	}

	// Send reset password email with rendered template
	subject := "Reset Your Password - Finance Manager"
	if err := mail.SendEmailWithTemplate(user.Email, subject, htmlBody); err != nil {
		logger.LogError(funcCtx, "failed to send reset password email", err, logrus.Fields{
			"user_id": user.ID.String(),
			"email":   user.Email,
		})
		// Don't return error to user for email sending failure for security reasons
		// Just log it and return success
	}

	logger.LogSuccess(funcCtx, "forgot password token generated and email sent", logrus.Fields{
		"user_id": user.ID.String(),
	})

	return nil
}

func (uc *AuthUseCase) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	funcCtx := "ResetPassword"

	// URL decode the token (replace spaces with + if needed)
	token := strings.ReplaceAll(req.Token, " ", "+")

	// Get user by forgot password token
	user, err := uc.userRepo.GetByForgotPasswordToken(ctx, token)
	if err != nil {
		logger.LogError(funcCtx, "invalid or expired reset token", err, logrus.Fields{
			"token": token,
		})
		return helpers.NewNotFoundError("invalid or expired reset token", "")
	}

	// Validate new password strength
	if err := auth.ValidatePasswordStrength(req.NewPassword); err != nil {
		logger.LogError(funcCtx, "new password validation failed", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return helpers.NewBadRequestError("password validation failed", err.Error())
	}

	// Decrypt and validate token
	_, timestamp, err := encryption.DecryptResetToken(token)
	if err != nil {
		logger.LogError(funcCtx, "failed to decrypt token", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return helpers.NewBadRequestError("invalid reset token", "")
	}

	// Check if token has expired (24 hours = 86400000 milliseconds)
	if err := encryption.ValidateResetTokenExpiry(timestamp, 86400000); err != nil {
		// Clear the expired token
		uc.userRepo.ClearForgotPasswordToken(ctx, user.ID)

		logger.LogError(funcCtx, "reset token expired", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return helpers.NewConflictError("reset token has expired", "")
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		logger.LogError(funcCtx, "failed to hash new password", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return helpers.NewInternalError("failed to hash new password", err.Error())
	}

	// Update user password and clear forgot password token
	user.Password = hashedPassword
	user.ForgotPasswordToken = ""

	if err := uc.userRepo.Update(ctx, user); err != nil {
		logger.LogError(funcCtx, "failed to update user password", err, logrus.Fields{
			"user_id": user.ID.String(),
		})
		return helpers.NewInternalError("failed to update password", err.Error())
	}

	logger.LogSuccess(funcCtx, "password reset successfully", logrus.Fields{
		"user_id": user.ID.String(),
	})

	return nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
