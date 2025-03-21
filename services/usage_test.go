package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/SemanticSugar/rightsizer/clients"
	mockClients "github.com/SemanticSugar/rightsizer/clients/mocks"
	"github.com/SemanticSugar/rightsizer/models"
	"github.com/SemanticSugar/rightsizer/services"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_UsageService_GetUsage_FailsToGetCPUAverage(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mockClients.NewMockCloudwatchClient(ctrl)
	usageService := services.NewUsageService(mockClient)
	ctx := context.Background()
	input := &services.GetUsageInput{
		ClusterName: "cluster",
		ServiceName: "service",
		TimeFrame:   5 * time.Hour,
	}
	mockClient.EXPECT().GetAverage(ctx, &clients.GetAverageInput{
		MetricName:  "CPUUtilization",
		ClusterName: input.ClusterName,
		ServiceName: input.ServiceName,
		TimeFrame:   input.TimeFrame,
	}).Return(nil, assert.AnError)
	usage, err := usageService.GetUsage(ctx, input)
	assert.Nil(t, usage)
	assert.Error(t, err)
	assert.Equal(t, "failed to get CPU average: assert.AnError general error for testing", err.Error())
}

func Test_UsageService_GetUsage_FailsToGetMemoryAverage(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mockClients.NewMockCloudwatchClient(ctrl)
	usageService := services.NewUsageService(mockClient)
	ctx := context.Background()
	input := &services.GetUsageInput{
		ClusterName: "cluster",
		ServiceName: "service",
		TimeFrame:   5 * time.Hour,
	}
	mockClient.EXPECT().GetAverage(ctx, &clients.GetAverageInput{
		MetricName:  "CPUUtilization",
		ClusterName: input.ClusterName,
		ServiceName: input.ServiceName,
		TimeFrame:   input.TimeFrame,
	}).Return(aws.Float64(50), nil)
	mockClient.EXPECT().GetAverage(ctx, &clients.GetAverageInput{
		MetricName:  "MemoryUtilization",
		ClusterName: input.ClusterName,
		ServiceName: input.ServiceName,
		TimeFrame:   input.TimeFrame,
	}).Return(nil, assert.AnError)
	usage, err := usageService.GetUsage(ctx, input)
	assert.Nil(t, usage)
	assert.Error(t, err)
	assert.Equal(t, "failed to get memory average: assert.AnError general error for testing", err.Error())
}

func Test_UsageService_GetUsage_Works(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mockClients.NewMockCloudwatchClient(ctrl)
	usageService := services.NewUsageService(mockClient)
	ctx := context.Background()
	input := &services.GetUsageInput{
		ClusterName: "cluster",
		ServiceName: "service",
		TimeFrame:   5 * time.Hour,
	}
	mockClient.EXPECT().GetAverage(ctx, &clients.GetAverageInput{
		MetricName:  "CPUUtilization",
		ClusterName: input.ClusterName,
		ServiceName: input.ServiceName,
		TimeFrame:   input.TimeFrame,
	}).Return(aws.Float64(50), nil)
	mockClient.EXPECT().GetAverage(ctx, &clients.GetAverageInput{
		MetricName:  "MemoryUtilization",
		ClusterName: input.ClusterName,
		ServiceName: input.ServiceName,
		TimeFrame:   input.TimeFrame,
	}).Return(aws.Float64(60), nil)
	usage, err := usageService.GetUsage(ctx, input)
	assert.NoError(t, err)
	assert.Equal(t, &models.Usage{
		CPU:    50,
		Memory: 60,
	}, usage)
}
