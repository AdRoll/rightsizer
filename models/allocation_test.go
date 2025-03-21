package models_test

import (
	"testing"

	"github.com/SemanticSugar/rightsizer/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func Test_ContainerAllocation_Fix_SameUsageAndTarget(t *testing.T) {
	c := models.ContainerAllocation{
		CPU:               100,
		Memory:            aws.Int32(200),
		MemoryReservation: aws.Int32(100),
	}
	usage := &models.Usage{
		CPU:    1,
		Memory: 1,
	}
	target := &models.Usage{
		CPU:    1,
		Memory: 1,
	}
	expected := models.ContainerAllocation{
		CPU:               100,
		Memory:            aws.Int32(200),
		MemoryReservation: aws.Int32(100),
	}
	actual := c.Fix(usage, target)
	assert.Equal(t, expected, *actual)
}

func Test_ContainerAllocation_Fix_TargetIsTwiceTheUsage(t *testing.T) {
	c := models.ContainerAllocation{
		CPU:               100,
		Memory:            aws.Int32(200),
		MemoryReservation: aws.Int32(100),
	}
	usage := &models.Usage{
		CPU:    45,
		Memory: 45,
	}
	target := &models.Usage{
		CPU:    90,
		Memory: 90,
	}
	expected := models.ContainerAllocation{
		CPU:               50,
		Memory:            aws.Int32(200),
		MemoryReservation: aws.Int32(50),
	}
	actual := c.Fix(usage, target)
	assert.Equal(t, expected, *actual)
}

func Test_ContainerAllocation_Fix_TargetIsHalfTheUsage(t *testing.T) {
	c := models.ContainerAllocation{
		CPU:               100,
		Memory:            aws.Int32(200),
		MemoryReservation: aws.Int32(100),
	}
	usage := &models.Usage{
		CPU:    90,
		Memory: 90,
	}
	target := &models.Usage{
		CPU:    45,
		Memory: 45,
	}
	expected := models.ContainerAllocation{
		CPU:               200,
		Memory:            aws.Int32(250),
		MemoryReservation: aws.Int32(200),
	}
	actual := c.Fix(usage, target)
	assert.Equal(t, expected, *actual)
}

func Test_ContainerAllocation_Fix_ValuesWontGoBelowThresholds(t *testing.T) {
	c := models.ContainerAllocation{
		CPU:               10,
		Memory:            aws.Int32(200),
		MemoryReservation: aws.Int32(10),
	}
	usage := &models.Usage{
		CPU:    1,
		Memory: 1,
	}
	target := &models.Usage{
		CPU:    90,
		Memory: 90,
	}
	expected := models.ContainerAllocation{
		CPU:               1,
		Memory:            aws.Int32(200),
		MemoryReservation: aws.Int32(6),
	}
	actual := c.Fix(usage, target)
	assert.Equal(t, expected, *actual)
}

func Test_ContainerAllocation_Fix_MemoryReservationIsNil(t *testing.T) {
	c := models.ContainerAllocation{
		CPU:    10,
		Memory: aws.Int32(200),
	}
	usage := &models.Usage{
		CPU:    1,
		Memory: 1,
	}
	target := &models.Usage{
		CPU:    90,
		Memory: 90,
	}
	expected := models.ContainerAllocation{
		CPU:    1,
		Memory: aws.Int32(200),
	}
	actual := c.Fix(usage, target)
	assert.Equal(t, expected, *actual)
}

func Test_ContainerAllocation_Fix_MemoryIsNil(t *testing.T) {
	c := models.ContainerAllocation{
		CPU:               10,
		MemoryReservation: aws.Int32(10),
	}
	usage := &models.Usage{
		CPU:    1,
		Memory: 1,
	}
	target := &models.Usage{
		CPU:    90,
		Memory: 90,
	}
	expected := models.ContainerAllocation{
		CPU:               1,
		MemoryReservation: aws.Int32(6),
	}
	actual := c.Fix(usage, target)
	assert.Equal(t, expected, *actual)
}

func Test_ServiceAllocation_Fix_SameUsageAndTarget(t *testing.T) {
	c := models.ServiceAllocation{
		CPU:    aws.String("100"),
		Memory: aws.String("200"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               100,
				Memory:            aws.Int32(200),
				MemoryReservation: aws.Int32(100),
			},
		},
	}
	usage := &models.Usage{
		CPU:    1,
		Memory: 1,
	}
	target := &models.Usage{
		CPU:    1,
		Memory: 1,
	}
	expected := models.ServiceAllocation{
		CPU:    aws.String("100"),
		Memory: aws.String("200"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               100,
				Memory:            aws.Int32(200),
				MemoryReservation: aws.Int32(100),
			},
		},
	}
	actual := c.Fix(usage, target)
	assert.Equal(t, expected, *actual)
}

func Test_ServiceAllocation_Fix_TargetIsTwiceTheUsage(t *testing.T) {
	c := models.ServiceAllocation{
		CPU:    aws.String("100"),
		Memory: aws.String("200"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               100,
				Memory:            aws.Int32(200),
				MemoryReservation: aws.Int32(100),
			},
		},
	}
	usage := &models.Usage{
		CPU:    45,
		Memory: 45,
	}
	target := &models.Usage{
		CPU:    90,
		Memory: 90,
	}
	expected := models.ServiceAllocation{
		CPU:    aws.String("50"),
		Memory: aws.String("100"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               50,
				Memory:            aws.Int32(200),
				MemoryReservation: aws.Int32(50),
			},
		},
	}
	actual := c.Fix(usage, target)
	assert.Equal(t, expected, *actual)
}

func Test_ServiceAllocation_Fix_TargetIsHalfTheUsage(t *testing.T) {
	c := models.ServiceAllocation{
		CPU:    aws.String("100"),
		Memory: aws.String("200"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               100,
				Memory:            aws.Int32(200),
				MemoryReservation: aws.Int32(100),
			},
		},
	}
	usage := &models.Usage{
		CPU:    90,
		Memory: 90,
	}
	target := &models.Usage{
		CPU:    45,
		Memory: 45,
	}
	expected := models.ServiceAllocation{
		CPU:    aws.String("200"),
		Memory: aws.String("400"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               200,
				Memory:            aws.Int32(250),
				MemoryReservation: aws.Int32(200),
			},
		},
	}
	actual := c.Fix(usage, target)
	assert.Equal(t, expected, *actual)
}

func Test_ServiceAllocation_Fix_ValuesWontGoBelowThresholds(t *testing.T) {
	c := models.ServiceAllocation{
		CPU:    aws.String("10"),
		Memory: aws.String("200"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               10,
				Memory:            aws.Int32(200),
				MemoryReservation: aws.Int32(10),
			},
		},
	}
	usage := &models.Usage{
		CPU:    1,
		Memory: 1,
	}
	target := &models.Usage{
		CPU:    90,
		Memory: 90,
	}
	expected := models.ServiceAllocation{
		CPU:    aws.String("1"),
		Memory: aws.String("6"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               1,
				Memory:            aws.Int32(200),
				MemoryReservation: aws.Int32(6),
			},
		},
	}
	actual := c.Fix(usage, target)
	assert.Equal(t, expected, *actual)
}

func Test_ServiceAllocation_Fix_MemoryIsNil(t *testing.T) {
	c := models.ServiceAllocation{
		CPU: aws.String("10"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               10,
				Memory:            aws.Int32(200),
				MemoryReservation: aws.Int32(10),
			},
		},
	}
	usage := &models.Usage{
		CPU:    1,
		Memory: 1,
	}
	target := &models.Usage{
		CPU:    90,
		Memory: 90,
	}
	expected := models.ServiceAllocation{
		CPU: aws.String("1"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               1,
				Memory:            aws.Int32(200),
				MemoryReservation: aws.Int32(6),
			},
		},
	}
	actual := c.Fix(usage, target)
	assert.Equal(t, expected, *actual)
}

func Test_ServiceAllocation_Fix_CPUIsNil(t *testing.T) {
	c := models.ServiceAllocation{
		Memory: aws.String("200"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               10,
				Memory:            aws.Int32(200),
				MemoryReservation: aws.Int32(10),
			},
		},
	}
	usage := &models.Usage{
		CPU:    1,
		Memory: 1,
	}
	target := &models.Usage{
		CPU:    90,
		Memory: 90,
	}
	expected := models.ServiceAllocation{
		Memory: aws.String("6"),
		ContainerAllocations: map[string]*models.ContainerAllocation{
			"container1": {
				CPU:               1,
				Memory:            aws.Int32(200),
				MemoryReservation: aws.Int32(6),
			},
		},
	}
	actual := c.Fix(usage, target)
	assert.Equal(t, expected, *actual)
}
