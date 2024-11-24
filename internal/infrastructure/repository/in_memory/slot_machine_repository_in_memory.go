package repository_in_memory

import (
	"context"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
	"sync"
)

type InMemorySlotMachineRepository struct {
	machines map[string]*model.SlotMachine
	mu       sync.RWMutex
}

func NewInMemorySlotMachineRepository() *InMemorySlotMachineRepository {
	return &InMemorySlotMachineRepository{
		machines: make(map[string]*model.SlotMachine),
	}
}

func (r *InMemorySlotMachineRepository) GetSlotMachine(ctx context.Context, id string) (*model.SlotMachine, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	machine, exists := r.machines[id]
	if !exists {
		return nil, repository.ErrSlotMachineNotFound
	}
	return machine, nil
}

func (r *InMemorySlotMachineRepository) UpdateSlotMachine(ctx context.Context, machine *model.SlotMachine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, exists := r.machines[machine.ID]
	if !exists {
		return repository.ErrSlotMachineNotFound
	}
	r.machines[machine.ID] = machine
	return nil
}

func (r *InMemorySlotMachineRepository) CreateSlotMachine(ctx context.Context, machine *model.SlotMachine) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.machines[machine.ID]; exists {
		return repository.ErrSlotMachineExists
	}
	r.machines[machine.ID] = machine
	return nil
}
