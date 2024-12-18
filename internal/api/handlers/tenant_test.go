package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"queue-manager/internal/models"
)

type MockTenantRepository struct {
	mock.Mock
}

func (m *MockTenantRepository) Create(ctx context.Context, tenant *models.CreateTenantRequest) (*models.Tenant, error) {
	args := m.Called(ctx, tenant)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tenant), args.Error(1)
}

type MockRabbitMQ struct {
	mock.Mock
}

func (m *MockRabbitMQ) CreateQueue(clientID string) error {
	args := m.Called(clientID)
	return args.Error(0)
}

func (m *MockRabbitMQ) StartConsumer(clientID string) error {
	args := m.Called(clientID)
	return args.Error(0)
}

func TestTenantHandler_CreateTenant(t *testing.T) {
	e := echo.New()
	mockRepo := new(MockTenantRepository)
	mockRMQ := new(MockRabbitMQ)
	logger := logrus.New()
	
	handler := NewTenantHandler(mockRepo, mockRMQ, logger)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func()
		expectedStatus int
	}{
		{
			name: "Valid tenant creation",
			requestBody: models.CreateTenantRequest{
				ClientID: "test-client",
				Name:     "Test Tenant",
			},
			setupMocks: func() {
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(&models.Tenant{
					ID:       1,
					ClientID: "test-client",
					Name:     "Test Tenant",
				}, nil)
				mockRMQ.On("CreateQueue", "test-client").Return(nil)
				mockRMQ.On("StartConsumer", "test-client").Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid request body",
			requestBody: map[string]string{
				"invalid": "data",
			},
			setupMocks: func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewBuffer(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.CreateTenant(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockRepo.AssertExpectations(t)
			mockRMQ.AssertExpectations(t)
		})
	}
} 