package clients

//go:generate mockgen -destination=mocks/ecs.go . ECSClient

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ECSClient interface {
	// GetService returns the service with the given name in the given cluster
	GetService(ctx context.Context, input *GetServiceInput) (*ecsTypes.Service, error)
	// GetTaskDefinition returns the task definition with the given name
	GetTaskDefinition(ctx context.Context, input *GetTaskDefinitionInput) (*ecsTypes.TaskDefinition, error)
}

// GetServiceInput is the input for GetService
type GetServiceInput struct {
	// Cluster is the name of the ECS cluster
	Cluster string
	// Service is the name of the ECS service
	Service string
}

type GetTaskDefinitionInput struct {
	// TaskDefinition is the name of the ECS task definition
	TaskDefinition string
}

type ecsClient struct {
	client *ecs.Client
}

func NewECSClient(client *ecs.Client) ECSClient {
	return &ecsClient{client: client}
}

func (c *ecsClient) GetService(ctx context.Context, input *GetServiceInput) (*ecsTypes.Service, error) {
	output, err := c.client.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  &input.Cluster,
		Services: []string{input.Service},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe service: %w", err)
	}
	if len(output.Services) == 0 {
		return nil, errors.New("service not found")
	}
	return &output.Services[0], nil
}

func (c *ecsClient) GetTaskDefinition(ctx context.Context, input *GetTaskDefinitionInput) (*ecsTypes.TaskDefinition, error) {
	output, err := c.client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &input.TaskDefinition,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe task definition: %w", err)
	}
	return output.TaskDefinition, nil
}
