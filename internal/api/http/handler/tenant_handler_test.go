package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"tenant/internal/api/http/handler"
	"tenant/internal/container"
	"tenant/internal/model"
	"tenant/internal/test"

	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func NewValidator() *requestValidator {
	return &requestValidator{}
}

type requestValidator struct{}

func (rv *requestValidator) Validate(i interface{}) (err error) {
	_, err = govalidator.ValidateStruct(i)
	return
}

func TestCreateTenantHandler(t *testing.T) {
	e := echo.New()
	e.Validator = &requestValidator{}
	mockComponent := test.InitMockComponent(t)

	hc := &container.HandlerComponent{
		TenantUsecase: mockComponent.TenantUsecase,
	}

	h := handler.NewTenantHandler(hc)

	var testCases = []struct {
		caseName     string
		requestBody  string
		mockSetup    func()
		expectedCode int
		expectedBody string
	}{
		{
			caseName:    "CreateTenant_Success",
			requestBody: `{"name":"Test Tenant"}`,
			mockSetup: func() {
				mockComponent.TenantUsecase.On("CreateTenant", mock.Anything, "Test Tenant").Return(&model.Tenant{
					ClientID: "test-client",
					Name:     "Test Tenant",
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"message":"Tenant created successfully"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			tc.mockSetup()

			req := httptest.NewRequest(http.MethodPost, "/tenants", strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.CreateTenant(c)
			if err != nil {
				t.Errorf("Handler returned error: %v", err)
			}

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedBody)
		})
	}
}

func TestDeleteTenantHandler(t *testing.T) {
	e := echo.New()
	e.Validator = &requestValidator{}
	mockComponent := test.InitMockComponent(t)

	hc := &container.HandlerComponent{
		TenantUsecase: mockComponent.TenantUsecase,
	}

	h := handler.NewTenantHandler(hc)

	var testCases = []struct {
		caseName     string
		clientID     string
		mockSetup    func()
		expectedCode int
		expectedBody string
	}{
		{
			caseName: "DeleteTenant_Success",
			clientID: "test-client",
			mockSetup: func() {
				mockComponent.TenantUsecase.On("DeleteTenant", mock.Anything, "test-client").Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"message":"Delete Tenant successfully"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			tc.mockSetup()

			path := fmt.Sprintf("/tenants/%s", tc.clientID)
			req := httptest.NewRequest(http.MethodDelete, path, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("clientID")
			c.SetParamValues(tc.clientID)

			err := h.DeleteTenant(c)
			if err != nil {
				t.Errorf("Handler returned error: %v", err)
			}

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestProcessTenantHandler(t *testing.T) {
	e := echo.New()
	e.Validator = &requestValidator{}
	mockComponent := test.InitMockComponent(t)

	hc := &container.HandlerComponent{
		TenantUsecase: mockComponent.TenantUsecase,
	}

	h := handler.NewTenantHandler(hc)

	var testCases = []struct {
		caseName     string
		requestBody  string
		clientID     string
		mockSetup    func()
		expectedCode int
		expectedBody string
	}{
		{
			caseName:    "ProcessTenant_Success",
			requestBody: `{"payload":"Test Tenant"}`,
			clientID:    "test-client",
			mockSetup: func() {
				mockComponent.TenantUsecase.On("ProcessPayload", mock.Anything, "test-client", "Test Tenant").Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"message":"Process Tenant successfully"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.caseName, func(t *testing.T) {
			tc.mockSetup()

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/tenants/%s/process", tc.clientID), strings.NewReader(tc.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("clientID")
			c.SetParamValues(tc.clientID)

			err := h.ProcessPayload(c)
			if err != nil {
				t.Errorf("Handler returned error: %v", err)
			}

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
