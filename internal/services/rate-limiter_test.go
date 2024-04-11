package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockRateLimiterRepository struct {
	mock.Mock
}

func (m *MockRateLimiterRepository) GetIpConnections(ctx context.Context, ip string) (int, error) {
	args := m.Called(ctx, ip)
	return args.Int(0), args.Error(1)
}

func (m *MockRateLimiterRepository) IncrementIpConnection(ctx context.Context, ip string) error {
	args := m.Called(ctx, ip)
	return args.Error(0)
}

func (m *MockRateLimiterRepository) GetApiKeyConnections(ctx context.Context, apiKey string) (int, error) {
	args := m.Called(ctx, apiKey)
	return args.Int(0), args.Error(1)
}

func (m *MockRateLimiterRepository) IncrementApiKeyConnection(ctx context.Context, apiKey string) error {
	args := m.Called(ctx, apiKey)
	return args.Error(0)
}

type RateLimiterServiceTestSuit struct {
	suite.Suite
	service RateLimiterChecker
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *RateLimiterServiceTestSuit) SetupTest() {
	testObj := new(MockRateLimiterRepository)
	testObj.On("GetIpConnections", context.Background(), "127.0.0.1").Return(2, nil)
	testObj.On("GetIpConnections", context.Background(), "192.168.0.1").Return(0, nil)
	testObj.On("GetApiKeyConnections", context.Background(), "token").Return(5, nil)
	testObj.On("GetApiKeyConnections", context.Background(), "allowed_token").Return(0, nil)
	testObj.On("IncrementIpConnection", context.Background(), "192.168.0.1").Return(nil)
	testObj.On("IncrementApiKeyConnection", context.Background(), "allowed_token").Return(nil)
	suite.service = NewRateLimiterRedisChecker(testObj, 1, 1)
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *RateLimiterServiceTestSuit) Test_ShouldReturnIPNotFoundError() {
	status, err := suite.service.CheckIpLimit("")
	suite.Equal(status, IP_REQUEST_LIMIT)
	suite.Error(err)
	suite.ErrorIs(err, ErrIpNotFound)
}

func (suite *RateLimiterServiceTestSuit) Test_ShouldReturnIpRequestLimitStatus() {
	status, err := suite.service.CheckIpLimit("127.0.0.1")
	suite.Equal(status, IP_REQUEST_LIMIT)
	suite.NoError(err)
}

func (suite *RateLimiterServiceTestSuit) Test_ShouldReturnAllowIpRequest() {
	status, err := suite.service.CheckIpLimit("192.168.0.1")
	suite.Equal(status, REQUEST_ALLOWED)
	suite.NoError(err)
}

func (suite *RateLimiterServiceTestSuit) Test_ShouldReturnErrApiKeyNotFoundError() {
	status, err := suite.service.CheckApiKeyLimit("")
	suite.Equal(status, TOKEN_REQUEST_LIMIT)
	suite.Error(err)
	suite.ErrorIs(err, ErrApiKeyNotFound)
}

func (suite *RateLimiterServiceTestSuit) Test_ShouldReturnApiKeyRequestLimitStatus() {
	status, err := suite.service.CheckApiKeyLimit("token")
	suite.Equal(status, TOKEN_REQUEST_LIMIT)
	suite.NoError(err)
}

func (suite *RateLimiterServiceTestSuit) Test_ShouldReturnAllowApiKeyRequest() {
	status, err := suite.service.CheckApiKeyLimit("allowed_token")
	suite.Equal(status, REQUEST_ALLOWED)
	suite.NoError(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterServiceTestSuit))
}
