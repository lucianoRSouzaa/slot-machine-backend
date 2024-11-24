package repository

import (
	"context"
	"errors"
	"slot-machine/internal/domain/model"
)

var (
	ErrSlotMachineNotFound = errors.New("slot machine not found")
	ErrSlotMachineExists   = errors.New("slot machine already exists")
)

type SlotMachineRepository interface {
	GetSlotMachine(ctx context.Context, id string) (*model.SlotMachine, error)
	UpdateSlotMachine(ctx context.Context, machine *model.SlotMachine) error
	CreateSlotMachine(ctx context.Context, machine *model.SlotMachine) error
}
