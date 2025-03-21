package services

import (
	"context"
	"fmt"

	"github.com/SemanticSugar/rightsizer/clients"
	"github.com/SemanticSugar/rightsizer/models"
)

type AllocationService interface {
	// GetAllocation returns the allocation of a service
	GetAllocation(ctx context.Context, input *GetAllocationInput) (*models.ServiceAllocation, error)
}

// GetAllocationInput is the input for GetAllocation
type GetAllocationInput struct {
	// ClusterName is the name of the cluster to get the allocation of
	ClusterName string
	// ServiceName is the name of the service to get the allocation of
	ServiceName string
}

type allocationService struct {
	ecsClient clients.ECSClient
}

func NewAllocationService(ecsClient clients.ECSClient) AllocationService {
	return &allocationService{ecsClient: ecsClient}
}

func (s *allocationService) GetAllocation(ctx context.Context, input *GetAllocationInput) (*models.ServiceAllocation, error) {
	service, err := s.ecsClient.GetService(ctx, &clients.GetServiceInput{
		Cluster: input.ClusterName,
		Service: input.ServiceName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}
	taskDefinition, err := s.ecsClient.GetTaskDefinition(ctx, &clients.GetTaskDefinitionInput{
		TaskDefinition: *service.TaskDefinition,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get task definition: %w", err)
	}
	if taskDefinition.ContainerDefinitions == nil {
		return nil, fmt.Errorf("task definition %s has no container definitions", *service.TaskDefinition)
	}
	containerAllocations := make(map[string]*models.ContainerAllocation)
	for _, cd := range taskDefinition.ContainerDefinitions {
		containerAllocations[*cd.Name] = &models.ContainerAllocation{
			CPU:               cd.Cpu,
			Memory:            cd.Memory,
			MemoryReservation: cd.MemoryReservation,
		}
	}
	return &models.ServiceAllocation{
		CPU:                  taskDefinition.Cpu,
		Memory:               taskDefinition.Memory,
		ContainerAllocations: containerAllocations,
	}, nil
}
