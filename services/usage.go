package services

import (
	"context"
	"fmt"
	"time"

	"github.com/SemanticSugar/rightsizer/clients"
	"github.com/SemanticSugar/rightsizer/models"
)

type UsageService interface {
	// GetUsage returns the usage of a service
	GetUsage(ctx context.Context, input *GetUsageInput) (*models.Usage, error)
}

// GetUsageInput is the input for GetUsage
type GetUsageInput struct {
	// ClusterName is the name of the cluster to get the usage of
	ClusterName string
	// ServiceName is the name of the service to get the usage of
	ServiceName string
	// TimeFrame is the period of time to get the usage over
	TimeFrame time.Duration
}

type usageService struct {
	cloudwatchClient clients.CloudwatchClient
}

func NewUsageService(cloudwatchClient clients.CloudwatchClient) UsageService {
	return &usageService{cloudwatchClient: cloudwatchClient}
}

func (s *usageService) GetUsage(ctx context.Context, input *GetUsageInput) (*models.Usage, error) {
	cpuAverage, err := s.cloudwatchClient.GetAverage(ctx, &clients.GetAverageInput{
		MetricName:  "CPUUtilization",
		ClusterName: input.ClusterName,
		ServiceName: input.ServiceName,
		TimeFrame:   input.TimeFrame,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU average: %w", err)
	}
	memoryAverage, err := s.cloudwatchClient.GetAverage(ctx, &clients.GetAverageInput{
		MetricName:  "MemoryUtilization",
		ClusterName: input.ClusterName,
		ServiceName: input.ServiceName,
		TimeFrame:   input.TimeFrame,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get memory average: %w", err)
	}
	return &models.Usage{
		CPU:    *cpuAverage,
		Memory: *memoryAverage,
	}, nil
}
