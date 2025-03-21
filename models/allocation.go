package models

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	cpuMinimum           = 1
	memoryMinimum        = 6
	defaultLeewayPercent = 25
)

// ContainerAllocation Type container allocation
type ContainerAllocation struct {
	// CPU is the CPU allocation
	CPU int32 `yaml:"cpu,omitempty"`
	// Memory is the hard limit for memory
	//
	// If present, this value should be larger than MemoryReservation
	Memory *int32 `yaml:"memory,omitempty"`
	// MemoryReservation is the soft limit for memory
	MemoryReservation *int32 `yaml:"memoryReservation,omitempty"`
}

// Fix fixes the container allocation
func (allocation *ContainerAllocation) Fix(usage *Usage, target *Usage) *ContainerAllocation {
	return allocation.
		fixCPU(float64(usage.CPU) / float64(target.CPU)).
		fixMemoryReservation(float64(usage.Memory) / float64(target.Memory)).
		makeConsistent(defaultLeewayPercent)
}

func (allocation *ContainerAllocation) fixCPU(ratio float64) *ContainerAllocation {
	currentCpu := float64(allocation.CPU)
	return &ContainerAllocation{
		CPU:               int32(max(currentCpu*ratio, cpuMinimum)),
		Memory:            allocation.Memory,
		MemoryReservation: allocation.MemoryReservation,
	}
}

func (allocation *ContainerAllocation) fixMemoryReservation(ratio float64) *ContainerAllocation {
	if allocation.MemoryReservation == nil {
		return allocation
	}
	currentMemory := float64(allocation.memoryReserved())
	return &ContainerAllocation{
		CPU:               allocation.CPU,
		Memory:            allocation.Memory,
		MemoryReservation: aws.Int32(int32(max(currentMemory*ratio, memoryMinimum))),
	}
}

func (allocation *ContainerAllocation) makeConsistent(leewayPercent float64) *ContainerAllocation {
	if allocation.MemoryReservation == nil || allocation.Memory == nil || *allocation.Memory > *allocation.MemoryReservation {
		return allocation
	}
	newMemory := float64(*allocation.MemoryReservation) * max(1, 1+leewayPercent/100)
	return &ContainerAllocation{
		CPU:               allocation.CPU,
		Memory:            aws.Int32(int32(max(1, newMemory))),
		MemoryReservation: allocation.MemoryReservation,
	}
}

func (allocation *ContainerAllocation) memoryReserved() int32 {
	if allocation.MemoryReservation != nil {
		return *allocation.MemoryReservation
	}
	if allocation.Memory != nil {
		return *allocation.Memory
	}
	return 1
}

type ServiceAllocation struct {
	CPU                  *string                         `yaml:"cpu,omitempty"`
	Memory               *string                         `yaml:"memory,omitempty"`
	ContainerAllocations map[string]*ContainerAllocation `yaml:"containerDefinitions,omitempty"`
}

func (allocation *ServiceAllocation) Fix(usage *Usage, target *Usage) *ServiceAllocation {
	return allocation.
		fixCPU(float64(usage.CPU)/float64(target.CPU)).
		fixMemory(float64(usage.Memory)/float64(target.Memory)).
		fixContainerAllocations(usage, target).
		makeConsistent()
}

func (allocation *ServiceAllocation) fixCPU(perShare float64) *ServiceAllocation {
	if allocation.CPU == nil {
		return allocation
	}
	currentCPU, _ := strconv.ParseFloat(*allocation.CPU, 64)
	return &ServiceAllocation{
		CPU:                  aws.String(strconv.FormatFloat(max(currentCPU*perShare, cpuMinimum), 'f', 0, 64)),
		Memory:               allocation.Memory,
		ContainerAllocations: allocation.ContainerAllocations,
	}
}

func (allocation *ServiceAllocation) fixMemory(perShare float64) *ServiceAllocation {
	if allocation.Memory == nil {
		return allocation
	}
	currentMemory, _ := strconv.ParseFloat(*allocation.Memory, 64)
	return &ServiceAllocation{
		CPU:                  allocation.CPU,
		Memory:               aws.String(strconv.FormatFloat(max(currentMemory*perShare, memoryMinimum), 'f', 0, 64)),
		ContainerAllocations: allocation.ContainerAllocations,
	}
}

func (allocation *ServiceAllocation) fixContainerAllocations(usage *Usage, target *Usage) *ServiceAllocation {
	newAllocations := make(map[string]*ContainerAllocation)
	for name, allocation := range allocation.ContainerAllocations {
		newAllocations[name] = allocation.Fix(usage, target)
	}
	return &ServiceAllocation{
		CPU:                  allocation.CPU,
		Memory:               allocation.Memory,
		ContainerAllocations: newAllocations,
	}
}

func (allocation *ServiceAllocation) makeConsistent() *ServiceAllocation {
	if allocation.Memory == nil {
		return allocation
	}
	requiredMemory := int64(0)
	for _, allocation := range allocation.ContainerAllocations {
		requiredMemory += int64(allocation.memoryReserved())
	}
	currentMemory, _ := strconv.ParseInt(*allocation.Memory, 10, 32)
	if requiredMemory < currentMemory {
		return allocation
	}
	return &ServiceAllocation{
		CPU:                  allocation.CPU,
		Memory:               aws.String(strconv.FormatInt(requiredMemory, 10)),
		ContainerAllocations: allocation.ContainerAllocations,
	}
}
