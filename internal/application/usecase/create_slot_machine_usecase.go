package usecase

import (
	"context"
	"errors"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"

	"github.com/google/uuid"
)

var (
	ErrSlotMachineAlreadyExists = errors.New("slot machine already exists")
)

type CreateSlotMachineUseCase struct {
	SlotMachineRepo repository.SlotMachineRepository
}

type CreateSlotMachineRequest struct {
	Level        int    `json:"level"`
	Balance      int    `json:"balance"`
	MultipleGain int    `json:"multiple_gain"`
	Description  string `json:"description"`
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
	id := uuid.New().String()

	if req.Level == 0 || req.MultipleGain == 0 || req.Description == "" {
		return nil, ErrValidate
	}

	machine := model.NewSlotMachine(id, req.Level, req.Balance, req.MultipleGain, req.Description)

	if err := uc.SlotMachineRepo.CreateSlotMachine(ctx, machine); err != nil {
		return nil, err
	}

	return &CreateSlotMachineResponse{
		Machine: *machine,
	}, nil
}
