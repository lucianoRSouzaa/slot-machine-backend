package usecase

import (
	"context"
	"errors"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
)

var (
	ErrSlotMachineAlreadyExists = errors.New("slot machine already exists")
)

type CreateSlotMachineUseCase struct {
	SlotMachineRepo repository.SlotMachineRepository
}

type CreateSlotMachineRequest struct {
	ID      string `json:"id"`
	Level   int    `json:"level"`
	Balance int    `json:"balance"`
}

type CreateSlotMachineResponse struct {
	Machine model.SlotMachine `json:"machine"`
}

func NewCreateSlotMachineUseCase(smr repository.SlotMachineRepository) *CreateSlotMachineUseCase {
	return &CreateSlotMachineUseCase{
		SlotMachineRepo: smr,
	}
}

func (uc *CreateSlotMachineUseCase) Execute(ctx context.Context, req *CreateSlotMachineRequest) (*CreateSlotMachineResponse, error) {
	_, err := uc.SlotMachineRepo.GetSlotMachine(ctx, req.ID)
	if err == nil {
		return nil, ErrSlotMachineAlreadyExists
	}
	if err != repository.ErrSlotMachineNotFound {
		return nil, err
	}

	machine := model.NewSlotMachine(req.ID, req.Level, req.Balance)

	if err := uc.SlotMachineRepo.CreateSlotMachine(ctx, machine); err != nil {
		return nil, err
	}

	return &CreateSlotMachineResponse{
		Machine: *machine,
	}, nil
}
