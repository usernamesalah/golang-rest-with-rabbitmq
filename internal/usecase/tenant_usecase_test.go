package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"tenant/internal/model"
	"tenant/internal/test/mockrepository"
	"tenant/internal/test/mockservice"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func mockInit() (*mockrepository.TenantRepository, *mockservice.Messagging) {
	return new(mockrepository.TenantRepository), new(mockservice.Messagging)
}

func TestCreateTenant(t *testing.T) {
	type params struct {
		tenantName string
	}

	ctx := context.Background()
	mockRepo, mockMQ := mockInit()
	tenantUsecase := NewTenantUsecase(mockRepo, mockMQ)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params)
		results      func(err error)
	}{
		{
			caseName: "CreateTenant_Success",
			params: params{
				tenantName: "Test Tenant",
			},
			expectations: func(params params) {
				mockRepo.On("CreateTenant", mock.Anything, mock.AnythingOfType("*model.Tenant")).Return(nil)
				mockMQ.On("CreateQueue", mock.Anything).Return(nil)
				mockMQ.On("StartQueue", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			results: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			caseName: "CreateTenant_FailToCreateQueue",
			params: params{
				tenantName: "Test Tenant",
			},
			expectations: func(params params) {
				mockRepo.On("CreateTenant", mock.Anything, mock.AnythingOfType("*model.Tenant")).Return(nil)
				mockMQ.On("CreateQueue", mock.Anything).Return(errors.New("failed to create queue"))
			},
			results: func(err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			_, err := tenantUsecase.CreateTenant(ctx, testCase.params.tenantName)
			testCase.results(err)

			mockRepo.AssertExpectations(t)
			mockMQ.AssertExpectations(t)
		})
	}
}

func TestDeleteTenant(t *testing.T) {
	type params struct {
		clientID string
	}

	ctx := context.Background()
	mockRepo, mockMQ := mockInit()
	tenantUsecase := NewTenantUsecase(mockRepo, mockMQ)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params)
		results      func(err error)
	}{
		{
			caseName: "DeleteTenant_Success",
			params: params{
				clientID: "test-client",
			},
			expectations: func(params params) {
				mockRepo.On("SoftDeleteTenant", mock.Anything, params.clientID).Return(nil)
				mockMQ.On("DeleteQueue", fmt.Sprintf("%s.process", params.clientID)).Return(nil)
			},
			results: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			caseName: "DeleteTenant_FailToSoftDelete",
			params: params{
				clientID: "test-client",
			},
			expectations: func(params params) {
				mockRepo.On("SoftDeleteTenant", mock.Anything, params.clientID).Return(errors.New("failed to soft delete"))
			},
			results: func(err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			err := tenantUsecase.DeleteTenant(ctx, testCase.params.clientID)
			testCase.results(err)

			mockRepo.AssertExpectations(t)
			mockMQ.AssertExpectations(t)
		})
	}
}

func TestProcessPayload(t *testing.T) {
	type params struct {
		clientID string
		payload  interface{}
	}

	ctx := context.Background()
	mockRepo, mockMQ := mockInit()
	tenantUsecase := NewTenantUsecase(mockRepo, mockMQ)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params params)
		results      func(err error)
	}{
		{
			caseName: "ProcessPayload_Success",
			params: params{
				clientID: "test-client",
				payload:  "test-payload",
			},
			expectations: func(params params) {
				mockRepo.On("GetTenantByClientID", mock.Anything, params.clientID).Return(&model.Tenant{ClientID: params.clientID, Name: "Test Tenant"}, nil)
				mockMQ.On("Publish", mock.Anything, fmt.Sprintf("%s.process", params.clientID), params.payload).Return(nil)
			},
			results: func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			caseName: "ProcessPayload_FailToPublish",
			params: params{
				clientID: "test-client",
				payload:  "test-payload",
			},
			expectations: func(params params) {
				mockRepo.On("GetTenantByClientID", mock.Anything, params.clientID).Return(&model.Tenant{ClientID: params.clientID, Name: "Test Tenant"}, nil)
				mockMQ.On("Publish", mock.Anything, fmt.Sprintf("%s.process", params.clientID), params.payload).Return(errors.New("failed to publish"))
			},
			results: func(err error) {
				assert.Nil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			err := tenantUsecase.ProcessPayload(ctx, testCase.params.clientID, testCase.params.payload)
			testCase.results(err)

			mockRepo.AssertExpectations(t)
			mockMQ.AssertExpectations(t)
		})
	}
}
func TestGetTenant(t *testing.T) {
	type params struct {
		clientID string
	}

	ctx := context.Background()
	mockRepo, mockMQ := mockInit()
	tenantUsecase := NewTenantUsecase(mockRepo, mockMQ)

	var testCases = []struct {
		caseName     string
		params       params
		expectations func(params params)
		results      func(tenant *model.Tenant, err error)
	}{
		{
			caseName: "GetTenant_Success",
			params: params{
				clientID: "test-client",
			},
			expectations: func(params params) {
				mockRepo.On("GetTenantByClientID", mock.Anything, params.clientID).Return(&model.Tenant{ClientID: params.clientID, Name: "Test Tenant"}, nil)
			},
			results: func(tenant *model.Tenant, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, tenant)
				assert.Equal(t, "test-client", tenant.ClientID)
				assert.Equal(t, "Test Tenant", tenant.Name)
			},
		},
		{
			caseName: "GetTenant_NotFound",
			params: params{
				clientID: "non-existent-client",
			},
			expectations: func(params params) {
				mockRepo.On("GetTenantByClientID", mock.Anything, params.clientID).Return(nil, sql.ErrNoRows)
			},
			results: func(tenant *model.Tenant, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, tenant)
				assert.Equal(t, sql.ErrNoRows, err)
			},
		},
		{
			caseName: "GetTenant_Error",
			params: params{
				clientID: "test-client",
			},
			expectations: func(params params) {
				mockRepo.On("GetTenantByClientID", mock.Anything, params.clientID).Return(nil, errors.New("unexpected error"))
			},
			results: func(tenant *model.Tenant, err error) {
				assert.NotNil(t, tenant)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.caseName, func(t *testing.T) {
			testCase.expectations(testCase.params)
			tenant, err := tenantUsecase.GetTenant(ctx, testCase.params.clientID)
			testCase.results(tenant, err)

			mockRepo.AssertExpectations(t)
		})
	}
}
