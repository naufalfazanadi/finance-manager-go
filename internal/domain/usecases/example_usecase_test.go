package usecases

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Test Suite
type ExampleUseCaseTestSuite struct {
	suite.Suite
	useCase ExampleUseCaseInterface
	ctx     context.Context
}

func (suite *ExampleUseCaseTestSuite) SetupTest() {
	suite.useCase = NewExampleUseCase()
	suite.ctx = context.Background()
}

// Test CreateExample
func (suite *ExampleUseCaseTestSuite) TestCreateExample_Success() {
	// Arrange
	req := map[string]interface{}{
		"name":        "Test Example",
		"description": "This is a test example",
	}

	// Act
	result, err := suite.useCase.CreateExample(suite.ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

// Test GetExample
func (suite *ExampleUseCaseTestSuite) TestGetExample_Success() {
	// Arrange
	id := uuid.New()

	// Act
	result, err := suite.useCase.GetExample(suite.ctx, id)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

// Test GetExamples
func (suite *ExampleUseCaseTestSuite) TestGetExamples_Success() {
	// Arrange
	queryParams := map[string]interface{}{
		"page":  1,
		"limit": 10,
	}

	// Act
	result, err := suite.useCase.GetExamples(suite.ctx, queryParams)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

// Test UpdateExample
func (suite *ExampleUseCaseTestSuite) TestUpdateExample_Success() {
	// Arrange
	id := uuid.New()
	req := map[string]interface{}{
		"name":        "Updated Example",
		"description": "This is an updated example",
	}

	// Act
	result, err := suite.useCase.UpdateExample(suite.ctx, id, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

// Test DeleteExample
func (suite *ExampleUseCaseTestSuite) TestDeleteExample_Success() {
	// Arrange
	id := uuid.New()

	// Act
	err := suite.useCase.DeleteExample(suite.ctx, id)

	// Assert
	assert.NoError(suite.T(), err)
}

// Run the test suite
func TestExampleUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleUseCaseTestSuite))
}
