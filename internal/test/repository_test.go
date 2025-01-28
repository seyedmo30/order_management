package test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/seyedmo30/order_management/internal/config"
	"github.com/seyedmo30/order_management/internal/dto"
	"github.com/seyedmo30/order_management/internal/interfaces"
	"github.com/seyedmo30/order_management/internal/repository"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

type RepositoryTestSuit struct {
	ctx               context.Context
	repositoryService interfaces.OrderRepository

	suite.Suite
}

func (p *RepositoryTestSuit) SetupSuite() {
	sm, err := GetSampleData()
	p.NoError(err)

	err = SetEnvFromStruct(&sm.Config.DatabaseConfig)
	p.NoError(err)

	var cfg config.DatabaseConfig

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	repositoryService := repository.NewOrderManagementRepository(cfg)
	p.repositoryService = repositoryService
	p.ctx = context.Background()

}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuit))
}

func (p *RepositoryTestSuit) Test0CreateInvoicePayment() {
	sm, err := GetSampleData()
	p.NoError(err)

	request := sm.RepositoryData.CreatOrderRepositoryRequest

	err = p.repositoryService.CreateOrder(p.ctx, request)

	p.NoError(err)

}

func (p *RepositoryTestSuit) Test1GetOrderByID() {
	sm, err := GetSampleData()
	p.NoError(err)

	request := sm.RepositoryData.GetOrderByIDRepositoryResponse

	res, err := p.repositoryService.GetOrderByID(p.ctx, request)

	p.NoError(err)

	fmt.Println(res)

}

func (p *RepositoryTestSuit) Test2GetNextHighPriorityReadyOrderAndLock() {

	res, err := p.repositoryService.GetNextHighPriorityReadyOrder(p.ctx)

	p.NoError(err)

	lockOrderOptimisticRepositoryRequest := dto.LockOrderOptimisticRepositoryRequest{BaseOrder: res.BaseOrder}

	err = p.repositoryService.LockOrderOptimistic(p.ctx, lockOrderOptimisticRepositoryRequest)
	p.NoError(err)

}

func (p *RepositoryTestSuit) Test3UpdateOrderByID() {
	sm, err := GetSampleData()
	p.NoError(err)

	request := sm.RepositoryData.UpdateOrderByIDRepositoryRequest

	err = p.repositoryService.UpdateOrderByID(p.ctx, request)

	p.NoError(err)

}

func (p *RepositoryTestSuit) Test4ListAggregateOrderReport() {

	res, err := p.repositoryService.ListAggregateOrderReport(p.ctx)

	p.NoError(err)

	fmt.Printf("%+v \n", res)
}

func (p *RepositoryTestSuit) TestFinalGetOrderByID() {
	sm, err := GetSampleData()
	p.NoError(err)

	request := sm.RepositoryData.GetOrderByIDRepositoryResponse

	res, err := p.repositoryService.GetOrderByID(p.ctx, request)

	p.NoError(err)

	fmt.Println(res)

}
