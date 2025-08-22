package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/naufalfazanadi/finance-manager-go/internal/domain/entities"
	"github.com/naufalfazanadi/finance-manager-go/internal/dto"
	"github.com/naufalfazanadi/finance-manager-go/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock interfaces
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) Create(ctx context.Context, wallet *entities.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Wallet, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Wallet), args.Error(1)
}

func (m *MockWalletRepository) GetOne(ctx context.Context, filter map[string]interface{}) (*entities.Wallet, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Wallet), args.Error(1)
}

func (m *MockWalletRepository) GetAll(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.Wallet, error) {
	args := m.Called(ctx, queryParams)
	return args.Get(0).([]*entities.Wallet), args.Error(1)
}

func (m *MockWalletRepository) Update(ctx context.Context, wallet *entities.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWalletRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockWalletRepository) CountWithFilters(ctx context.Context, queryParams *dto.QueryParams) (int64, error) {
	args := m.Called(ctx, queryParams)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockWalletRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWalletRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWalletRepository) Restore(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

// Add other required methods for UserRepository interface
func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByIDWithPreload(ctx context.Context, id uuid.UUID, preloadRelations []string) (*entities.User, error) {
	args := m.Called(ctx, id, preloadRelations)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmailHash(ctx context.Context, emailHash string) (*entities.User, error) {
	args := m.Called(ctx, emailHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByForgotPasswordToken(ctx context.Context, token string) (*entities.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetOne(ctx context.Context, filter map[string]interface{}) (*entities.User, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.User, error) {
	args := m.Called(ctx, queryParams)
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) CountWithFilters(ctx context.Context, queryParams *dto.QueryParams) (int64, error) {
	args := m.Called(ctx, queryParams)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Restore(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateForgotPasswordToken(ctx context.Context, userID uuid.UUID, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockUserRepository) ClearForgotPasswordToken(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetWithDeleted(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.User, error) {
	args := m.Called(ctx, queryParams)
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetOnlyDeleted(ctx context.Context, queryParams *dto.QueryParams) ([]*entities.User, error) {
	args := m.Called(ctx, queryParams)
	return args.Get(0).([]*entities.User), args.Error(1)
}

// Test Suite
type WalletUseCaseTestSuite struct {
	suite.Suite
	useCase    WalletUseCaseInterface
	walletRepo *MockWalletRepository
	userRepo   *MockUserRepository
	ctx        context.Context
}

func (suite *WalletUseCaseTestSuite) SetupTest() {
	// Initialize logger for tests
	logger.Init("info")

	suite.walletRepo = new(MockWalletRepository)
	suite.userRepo = new(MockUserRepository)
	suite.useCase = NewWalletUseCase(suite.walletRepo, suite.userRepo)
	suite.ctx = context.Background()
}

func (suite *WalletUseCaseTestSuite) TearDownTest() {
	suite.walletRepo.AssertExpectations(suite.T())
	suite.userRepo.AssertExpectations(suite.T())
}

// Test CreateWallet
func (suite *WalletUseCaseTestSuite) TestCreateWallet_Success() {
	// Arrange
	userID := uuid.New()
	req := &dto.CreateWalletRequest{
		Name:     "Test Wallet",
		Type:     "personal",
		Category: "income",
		Balance:  1000.0,
		Currency: "IDR",
		UserID:   userID,
	}

	// Mock: wallet doesn't exist
	suite.walletRepo.On("GetOne", suite.ctx, map[string]interface{}{
		"name":    req.Name,
		"user_id": req.UserID,
	}).Return((*entities.Wallet)(nil), errors.New("wallet not found"))

	// Mock: create wallet succeeds
	suite.walletRepo.On("Create", suite.ctx, mock.MatchedBy(func(wallet *entities.Wallet) bool {
		return wallet.Name == req.Name &&
			wallet.Type == req.Type &&
			wallet.Category == req.Category &&
			wallet.Balance == req.Balance &&
			wallet.Currency == req.Currency &&
			wallet.UserID == req.UserID
	})).Return(nil)

	// Act
	result, err := suite.useCase.CreateWallet(suite.ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *WalletUseCaseTestSuite) TestCreateWallet_WalletExists() {
	// Arrange
	userID := uuid.New()
	req := &dto.CreateWalletRequest{
		Name:     "Test Wallet",
		Type:     "personal",
		Category: "income",
		Balance:  1000.0,
		Currency: "IDR",
		UserID:   userID,
	}

	existingWallet := &entities.Wallet{
		ID:       uuid.New(),
		Name:     req.Name,
		UserID:   req.UserID,
		Type:     req.Type,
		Category: req.Category,
	}

	// Mock: wallet already exists
	suite.walletRepo.On("GetOne", suite.ctx, map[string]interface{}{
		"name":    req.Name,
		"user_id": req.UserID,
	}).Return(existingWallet, nil)

	// Act
	result, err := suite.useCase.CreateWallet(suite.ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *WalletUseCaseTestSuite) TestCreateWallet_CreateFails() {
	// Arrange
	userID := uuid.New()
	req := &dto.CreateWalletRequest{
		Name:     "Test Wallet",
		Type:     "personal",
		Category: "income",
		Balance:  1000.0,
		Currency: "IDR",
		UserID:   userID,
	}

	// Mock: wallet doesn't exist
	suite.walletRepo.On("GetOne", suite.ctx, map[string]interface{}{
		"name":    req.Name,
		"user_id": req.UserID,
	}).Return((*entities.Wallet)(nil), errors.New("wallet not found"))

	// Mock: create wallet fails
	suite.walletRepo.On("Create", suite.ctx, mock.AnythingOfType("*entities.Wallet")).
		Return(errors.New("database error"))

	// Act
	result, err := suite.useCase.CreateWallet(suite.ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// Test GetWallet
func (suite *WalletUseCaseTestSuite) TestGetWallet_Success() {
	// Arrange
	walletID := uuid.New()
	userID := uuid.New()
	wallet := &entities.Wallet{
		ID:       walletID,
		Name:     "Test Wallet",
		Type:     "personal",
		Category: "income",
		Balance:  1000.0,
		Currency: "IDR",
		UserID:   userID,
		User: entities.User{
			ID: userID,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock: get wallet succeeds
	suite.walletRepo.On("GetByID", suite.ctx, walletID).Return(wallet, nil)

	// Act
	result, err := suite.useCase.GetWallet(suite.ctx, walletID, userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *WalletUseCaseTestSuite) TestGetWallet_NotFound() {
	// Arrange
	walletID := uuid.New()
	userID := uuid.New()

	// Mock: wallet not found
	suite.walletRepo.On("GetByID", suite.ctx, walletID).
		Return((*entities.Wallet)(nil), errors.New("wallet not found"))

	// Act
	result, err := suite.useCase.GetWallet(suite.ctx, walletID, userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *WalletUseCaseTestSuite) TestGetWallet_UnauthorizedAccess() {
	// Arrange
	walletID := uuid.New()
	walletOwnerID := uuid.New()
	loggedUserID := uuid.New() // Different user
	wallet := &entities.Wallet{
		ID:     walletID,
		UserID: walletOwnerID,
		Name:   "Test Wallet",
	}

	// Mock: get wallet succeeds but user is different
	suite.walletRepo.On("GetByID", suite.ctx, walletID).Return(wallet, nil)

	// Act
	result, err := suite.useCase.GetWallet(suite.ctx, walletID, loggedUserID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// Test GetWallets
func (suite *WalletUseCaseTestSuite) TestGetWallets_Success() {
	// Arrange
	queryParams := &dto.QueryParams{
		PaginationQuery: &dto.PaginationQuery{
			Page:  1,
			Limit: 10,
		},
		FilterQuery:  &dto.FilterQuery{},
		LoggedUserID: uuid.New(),
	}

	wallets := []*entities.Wallet{
		{
			ID:        uuid.New(),
			Name:      "Wallet 1",
			Type:      "personal",
			Category:  "income",
			Balance:   1000.0,
			Currency:  "IDR",
			UserID:    queryParams.LoggedUserID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Wallet 2",
			Type:      "business",
			Category:  "expense",
			Balance:   500.0,
			Currency:  "USD",
			UserID:    queryParams.LoggedUserID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mock: get wallets succeeds
	suite.walletRepo.On("GetAll", suite.ctx, queryParams).Return(wallets, nil)
	suite.walletRepo.On("CountWithFilters", suite.ctx, queryParams).Return(int64(2), nil)

	// Act
	result, err := suite.useCase.GetWallets(suite.ctx, queryParams)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.NotEmpty(suite.T(), result.Data)
	assert.NotNil(suite.T(), result.Meta)
}

func (suite *WalletUseCaseTestSuite) TestGetWallets_GetAllFails() {
	// Arrange
	queryParams := &dto.QueryParams{
		PaginationQuery: &dto.PaginationQuery{Page: 1, Limit: 10},
		FilterQuery:     &dto.FilterQuery{},
		LoggedUserID:    uuid.New(),
	}

	// Mock: get wallets fails
	suite.walletRepo.On("GetAll", suite.ctx, queryParams).
		Return(([]*entities.Wallet)(nil), errors.New("database error"))

	// Act
	result, err := suite.useCase.GetWallets(suite.ctx, queryParams)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *WalletUseCaseTestSuite) TestGetWallets_CountFails() {
	// Arrange
	queryParams := &dto.QueryParams{
		PaginationQuery: &dto.PaginationQuery{Page: 1, Limit: 10},
		FilterQuery:     &dto.FilterQuery{},
		LoggedUserID:    uuid.New(),
	}

	wallets := []*entities.Wallet{}

	// Mock: get wallets succeeds, count fails
	suite.walletRepo.On("GetAll", suite.ctx, queryParams).Return(wallets, nil)
	suite.walletRepo.On("CountWithFilters", suite.ctx, queryParams).
		Return(int64(0), errors.New("database error"))

	// Act
	result, err := suite.useCase.GetWallets(suite.ctx, queryParams)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// Test UpdateWallet
func (suite *WalletUseCaseTestSuite) TestUpdateWallet_Success() {
	// Arrange
	walletID := uuid.New()
	userID := uuid.New()
	newUserID := uuid.New()

	existingWallet := &entities.Wallet{
		ID:       walletID,
		Name:     "Old Wallet",
		Type:     "personal",
		Category: "income",
		Balance:  1000.0,
		Currency: "IDR",
		UserID:   userID,
	}

	req := &dto.UpdateWalletRequest{
		Name:     "Updated Wallet",
		Type:     "business",
		Category: "expense",
		Balance:  2000.0,
		Currency: "USD",
		UserID:   newUserID,
	}

	newUser := &entities.User{
		ID: newUserID,
	}

	// Mock: get existing wallet
	suite.walletRepo.On("GetByID", suite.ctx, walletID).Return(existingWallet, nil)

	// Mock: get new user
	suite.userRepo.On("GetByID", suite.ctx, newUserID).Return(newUser, nil)

	// Mock: update wallet
	suite.walletRepo.On("Update", suite.ctx, mock.MatchedBy(func(wallet *entities.Wallet) bool {
		return wallet.ID == walletID &&
			wallet.Name == req.Name &&
			wallet.Type == req.Type &&
			wallet.Category == req.Category &&
			wallet.Balance == req.Balance &&
			wallet.Currency == req.Currency &&
			wallet.UserID == req.UserID
	})).Return(nil)

	// Act
	result, err := suite.useCase.UpdateWallet(suite.ctx, walletID, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *WalletUseCaseTestSuite) TestUpdateWallet_PartialUpdate() {
	// Arrange
	walletID := uuid.New()
	existingWallet := &entities.Wallet{
		ID:       walletID,
		Name:     "Old Wallet",
		Type:     "personal",
		Category: "income",
		Balance:  1000.0,
		Currency: "IDR",
		UserID:   uuid.New(),
	}

	req := &dto.UpdateWalletRequest{
		Name: "Updated Wallet", // Only update name
	}

	// Mock: get existing wallet
	suite.walletRepo.On("GetByID", suite.ctx, walletID).Return(existingWallet, nil)

	// Mock: update wallet
	suite.walletRepo.On("Update", suite.ctx, mock.MatchedBy(func(wallet *entities.Wallet) bool {
		return wallet.ID == walletID &&
			wallet.Name == req.Name &&
			wallet.Type == existingWallet.Type && // Should remain unchanged
			wallet.Category == existingWallet.Category &&
			wallet.Balance == existingWallet.Balance &&
			wallet.Currency == existingWallet.Currency &&
			wallet.UserID == existingWallet.UserID
	})).Return(nil)

	// Act
	result, err := suite.useCase.UpdateWallet(suite.ctx, walletID, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func (suite *WalletUseCaseTestSuite) TestUpdateWallet_WalletNotFound() {
	// Arrange
	walletID := uuid.New()
	req := &dto.UpdateWalletRequest{
		Name: "Updated Wallet",
	}

	// Mock: wallet not found
	suite.walletRepo.On("GetByID", suite.ctx, walletID).
		Return((*entities.Wallet)(nil), errors.New("wallet not found"))

	// Act
	result, err := suite.useCase.UpdateWallet(suite.ctx, walletID, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *WalletUseCaseTestSuite) TestUpdateWallet_UserNotFound() {
	// Arrange
	walletID := uuid.New()
	newUserID := uuid.New()
	existingWallet := &entities.Wallet{
		ID:     walletID,
		Name:   "Old Wallet",
		UserID: uuid.New(),
	}

	req := &dto.UpdateWalletRequest{
		UserID: newUserID,
	}

	// Mock: get existing wallet
	suite.walletRepo.On("GetByID", suite.ctx, walletID).Return(existingWallet, nil)

	// Mock: user not found
	suite.userRepo.On("GetByID", suite.ctx, newUserID).
		Return((*entities.User)(nil), errors.New("user not found"))

	// Act
	result, err := suite.useCase.UpdateWallet(suite.ctx, walletID, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *WalletUseCaseTestSuite) TestUpdateWallet_UpdateFails() {
	// Arrange
	walletID := uuid.New()
	existingWallet := &entities.Wallet{
		ID:     walletID,
		Name:   "Old Wallet",
		UserID: uuid.New(),
	}

	req := &dto.UpdateWalletRequest{
		Name: "Updated Wallet",
	}

	// Mock: get existing wallet
	suite.walletRepo.On("GetByID", suite.ctx, walletID).Return(existingWallet, nil)

	// Mock: update fails
	suite.walletRepo.On("Update", suite.ctx, mock.AnythingOfType("*entities.Wallet")).
		Return(errors.New("database error"))

	// Act
	result, err := suite.useCase.UpdateWallet(suite.ctx, walletID, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
}

// Test DeleteWallet
func (suite *WalletUseCaseTestSuite) TestDeleteWallet_Success() {
	// Arrange
	walletID := uuid.New()
	wallet := &entities.Wallet{
		ID:   walletID,
		Name: "Test Wallet",
	}

	// Mock: get wallet succeeds
	suite.walletRepo.On("GetByID", suite.ctx, walletID).Return(wallet, nil)

	// Mock: soft delete succeeds
	suite.walletRepo.On("SoftDelete", suite.ctx, walletID).Return(nil)

	// Act
	err := suite.useCase.DeleteWallet(suite.ctx, walletID)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *WalletUseCaseTestSuite) TestDeleteWallet_WalletNotFound() {
	// Arrange
	walletID := uuid.New()

	// Mock: wallet not found
	suite.walletRepo.On("GetByID", suite.ctx, walletID).
		Return((*entities.Wallet)(nil), errors.New("wallet not found"))

	// Act
	err := suite.useCase.DeleteWallet(suite.ctx, walletID)

	// Assert
	assert.Error(suite.T(), err)
}

func (suite *WalletUseCaseTestSuite) TestDeleteWallet_SoftDeleteFails() {
	// Arrange
	walletID := uuid.New()
	wallet := &entities.Wallet{
		ID:   walletID,
		Name: "Test Wallet",
	}

	// Mock: get wallet succeeds
	suite.walletRepo.On("GetByID", suite.ctx, walletID).Return(wallet, nil)

	// Mock: soft delete fails
	suite.walletRepo.On("SoftDelete", suite.ctx, walletID).
		Return(errors.New("database error"))

	// Act
	err := suite.useCase.DeleteWallet(suite.ctx, walletID)

	// Assert
	assert.Error(suite.T(), err)
}

// Test UUID.Nil handling in UpdateWallet
func (suite *WalletUseCaseTestSuite) TestUpdateWallet_UserIDIsNil() {
	// Arrange
	walletID := uuid.New()
	existingWallet := &entities.Wallet{
		ID:       walletID,
		Name:     "Old Wallet",
		Type:     "personal",
		Category: "income",
		Balance:  1000.0,
		Currency: "IDR",
		UserID:   uuid.New(),
	}

	req := &dto.UpdateWalletRequest{
		Name:   "Updated Wallet",
		UserID: uuid.Nil, // Should not update UserID
	}

	// Mock: get existing wallet
	suite.walletRepo.On("GetByID", suite.ctx, walletID).Return(existingWallet, nil)

	// Mock: update wallet (UserID should remain unchanged)
	suite.walletRepo.On("Update", suite.ctx, mock.MatchedBy(func(wallet *entities.Wallet) bool {
		return wallet.ID == walletID &&
			wallet.Name == req.Name &&
			wallet.UserID == existingWallet.UserID // Should remain unchanged
	})).Return(nil)

	// Act
	result, err := suite.useCase.UpdateWallet(suite.ctx, walletID, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

// Run the test suite
func TestWalletUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(WalletUseCaseTestSuite))
}
