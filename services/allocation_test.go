package services_test

import (
	"context"
	"testing"

	"github.com/SemanticSugar/rightsizer/clients"
	mockClients "github.com/SemanticSugar/rightsizer/clients/mocks"
	"github.com/SemanticSugar/rightsizer/models"
	"github.com/SemanticSugar/rightsizer/services"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_AllocationService_GetAllocation_FailsToGetService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mockClients.NewMockECSClient(ctrl)
	allocationService := services.NewAllocationService(mockClient)
	ctx := context.Background()
	input := &services.GetAllocationInput{
		ClusterName: "cluster",
		ServiceName: "service",
	}

	mockClient.EXPECT().GetService(ctx, gomock.Any()).Return(nil, assert.AnError)

	allocation, err := allocationService.GetAllocation(context.Background(), input)
	assert.Nil(t, allocation)
	assert.Error(t, err)
	assert.Equal(t, "failed to get service: assert.AnError general error for testing", err.Error())
}

func Test_AllocationService_GetAllocation_FailsToGetTaskDefinition(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mockClients.NewMockECSClient(ctrl)
	allocationService := services.NewAllocationService(mockClient)
	ctx := context.Background()
	input := &services.GetAllocationInput{
		ClusterName: "cluster",
		ServiceName: "service",
	}
	taskDefinitionName := "taskDefinition:1"
	mockClient.EXPECT().GetService(ctx, &clients.GetServiceInput{
		Cluster: input.ClusterName,
		Service: input.ServiceName,
	}).Return(&ecsTypes.Service{
		TaskDefinition: &taskDefinitionName,
	}, nil)
	mockClient.EXPECT().GetTaskDefinition(ctx, &clients.GetTaskDefinitionInput{
		TaskDefinition: taskDefinitionName,
	}).Return(nil, assert.AnError)

	allocation, err := allocationService.GetAllocation(context.Background(), input)
	assert.Nil(t, allocation)
	assert.Error(t, err)
	assert.Equal(t, "failed to get task definition: assert.AnError general error for testing", err.Error())
}

func Test_AllocationService_GetAllocation_TaskDefinitionHasNoContainerDefinitions(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mockClients.NewMockECSClient(ctrl)
	allocationService := services.NewAllocationService(mockClient)
	ctx := context.Background()
	input := &services.GetAllocationInput{
		ClusterName: "cluster",
		ServiceName: "service",
	}
	taskDefinitionName := "taskDefinition:1"
	mockClient.EXPECT().GetService(ctx, &clients.GetServiceInput{
		Cluster: input.ClusterName,
		Service: input.ServiceName,
	}).Return(&ecsTypes.Service{
		TaskDefinition: &taskDefinitionName,
	}, nil)
	mockClient.EXPECT().GetTaskDefinition(ctx, &clients.GetTaskDefinitionInput{
		TaskDefinition: taskDefinitionName,
	}).Return(&ecsTypes.TaskDefinition{}, nil)

	allocation, err := allocationService.GetAllocation(context.Background(), input)
	assert.Nil(t, allocation)
	assert.Error(t, err)
	assert.Equal(t, "task definition taskDefinition:1 has no container definitions", err.Error())
}

func Test_AllocationService_GetAllocation_Works(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mockClients.NewMockECSClient(ctrl)
	allocationService := services.NewAllocationService(mockClient)
	ctx := context.Background()
	input := &services.GetAllocationInput{
		ClusterName: "cluster",
		ServiceName: "service",
	}
	taskDefinitionName := "taskDefinition:1"
	mockClient.EXPECT().GetService(ctx, &clients.GetServiceInput{
		Cluster: input.ClusterName,
		Service: input.ServiceName,
	}).Return(&ecsTypes.Service{
		TaskDefinition: &taskDefinitionName,
	}, nil)
	mockClient.EXPECT().GetTaskDefinition(ctx, &clients.GetTaskDefinitionInput{
		TaskDefinition: taskDefinitionName,
	}).Return(&ecsTypes.TaskDefinition{
		Cpu:    aws.String("512"),
		Memory: aws.String("1024"),
		ContainerDefinitions: []ecsTypes.ContainerDefinition{
			{
				Name:              aws.String("container1"),
				Cpu:               512,
				Memory:            aws.Int32(1024),
				MemoryReservation: aws.Int32(512),
			},
		},
	}, nil)
	allocation, err := allocationService.GetAllocation(ctx, input)
	assert.NoError(t, err)
	assert.Equal(t, &models.ServiceAllocation{
		CPU:    aws.String("512"),
		Memory: aws.String("1024"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               512,
				Memory:            aws.Int32(1024),
				MemoryReservation: aws.Int32(512),
			},
		},
	}, allocation)

}
