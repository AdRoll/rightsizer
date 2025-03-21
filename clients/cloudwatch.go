package clients

//go:generate mockgen -destination=mocks/cloudwatch.go . CloudwatchClient

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type CloudwatchClient interface {
	// GetAverage returns the average value of a metric over a period of time
	GetAverage(ctx context.Context, input *GetAverageInput) (*float64, error)
}

// GetAverageInput is the input for the GetAverage function
type GetAverageInput struct {
	// MetricName is the name of the metric to get the average value of
	MetricName string
	// ClusterName is the name of the cluster to get the average value of
	ClusterName string
	// ServiceName is the name of the service to get the average value of
	ServiceName string
	// TimeFrame is the period of time to get the average value over
	TimeFrame time.Duration
}

type cloudWatchClient struct {
	client *cloudwatch.Client
}

// NewCloudWatchClient returns a new CloudWatchClient
func NewCloudWatchClient(client *cloudwatch.Client) CloudwatchClient {
	return &cloudWatchClient{client: client}
}

func (c *cloudWatchClient) GetAverage(ctx context.Context, input *GetAverageInput) (*float64, error) {
	now := time.Now()
	period := int32(input.TimeFrame.Seconds())
	startTime := now.Add(-input.TimeFrame)
	getMetricsInput := &cloudwatch.GetMetricStatisticsInput{
		MetricName: &input.MetricName,
		Dimensions: []cloudwatchTypes.Dimension{
			{
				Name:  aws.String("ClusterName"),
				Value: &input.ClusterName,
			},
			{
				Name:  aws.String("ServiceName"),
				Value: &input.ServiceName,
			},
		},
		StartTime: &startTime,
		EndTime:   &now,
		Namespace: aws.String("AWS/ECS"),
		Period:    &period,
		Statistics: []cloudwatchTypes.Statistic{
			cloudwatchTypes.StatisticAverage,
		},
	}
	metricsResponse, err := c.client.GetMetricStatistics(ctx, getMetricsInput)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}
	if len(metricsResponse.Datapoints) != 1 {
		return nil, fmt.Errorf("failed to get %s data points from CloudWatch", input.MetricName)
	}
	return metricsResponse.Datapoints[0].Average, nil
}
