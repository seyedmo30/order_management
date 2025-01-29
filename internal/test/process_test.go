package test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/seyedmo30/order_management/internal/config"
	"github.com/seyedmo30/order_management/internal/interfaces"
	"github.com/seyedmo30/order_management/internal/process"
	"github.com/seyedmo30/order_management/pkg"
	"github.com/stretchr/testify/suite"
)

type ProcessOrderTestSuite struct {
	ctx        context.Context
	processSvc interfaces.Process

	suite.Suite
}

// MockConfig is a mock type for the config.App interface

func (p *ProcessOrderTestSuite) SetupSuite() {
	sm, err := GetSampleData()
	p.NoError(err)

	err = SetEnvFromStruct(&sm.Config.ServiceConfig)
	p.NoError(err)

	var cfg config.App

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the service with mock config
	p.processSvc = process.NewProcessUseCase(cfg)
	p.ctx = context.Background()
}


func TestProcessOrder(t *testing.T) {
	suite.Run(t, new(ProcessOrderTestSuite))
}


func (p *ProcessOrderTestSuite) TestProcessOrder_Success() {
	// Mock the timeout to 5 seconds

	// Call ProcessOrder with a processing time of 3 seconds
	status, err := p.processSvc.ProcessOrder(p.ctx, 1)

	// Assert no error and the status is as expected
	p.NoError(err)
	p.Equal(pkg.StatusOrderManagementProcessed, status)

	// Ensure that expectations on the mock were met

}

func (p *ProcessOrderTestSuite) TestProcessOrder_Failure() {
	// Mock the timeout to 5 seconds

	// Create a context with a short timeout of 1 second
	ctx, cancel := context.WithTimeout(p.ctx, time.Second*10)
	defer cancel()

	// Call ProcessOrder with a processing time of 10 seconds
	status, err := p.processSvc.ProcessOrder(ctx, 10)

	// Assert no error and the status is "failed" due to context timeout
	p.NoError(err)
	p.Equal(pkg.StatusOrderManagementFailed, status)

	// Ensure that expectations on the mock were met

}
