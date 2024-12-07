package usecase

import (
	"context"
	"slot-machine/internal/domain/contextkeys"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
)

type GetSlotMachineBalanceUseCase struct {
	SlotMachineRepo repository.SlotMachineRepository
}

type GetSlotMachineBalanceRequest struct {
	MachineID string `json:"machine_id"`
}

type GetSlotMachineBalanceResponse struct {
	Machine model.SlotMachine `json:"machine"`
}

func NewGetSlotMachineBalanceUseCase(smr repository.SlotMachineRepository) *GetSlotMachineBalanceUseCase {
	return &GetSlotMachineBalanceUseCase{
		SlotMachineRepo: smr,
	}
}

func (uc *GetSlotMachineBalanceUseCase) Execute(ctx context.Context, req *GetSlotMachineBalanceRequest) (*GetSlotMachineBalanceResponse, error) {
	isAdmin, ok := ctx.Value(contextkeys.ContextKeyIsAdmin).(bool)
	if !ok || !isAdmin {
		return nil, ErrUnauthorized
	}

	machine, err := uc.SlotMachineRepo.GetSlotMachine(ctx, req.MachineID)
	if err != nil {
		return nil, err
	}

	return &GetSlotMachineBalanceResponse{
		Machine: *machine,
	}, nil
}
