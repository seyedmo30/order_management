package dto

type CreatOrderRepositoryRequest struct {
	BaseOrder
}

type UpdateOrderByIDRepositoryRequest struct {
	BaseOrder
}

type GetOrderByIDRepositoryResponse struct {
	BaseOrder
}

type GetNextHighPriorityReadyOrderRepositoryResponse struct {
	BaseOrder
}

type PatchTransferByIdRepositoryRequest struct {
	BaseOrder
}
