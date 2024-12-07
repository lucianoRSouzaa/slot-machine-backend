package usecase

import (
	"context"
	"slot-machine/internal/domain/contextkeys"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
	repository_in_memory "slot-machine/internal/infrastructure/repository/in_memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSlotMachineBalanceUseCase(t *testing.T) {
	slotRepo := repository_in_memory.NewInMemorySlotMachineRepository()

	getSlotMachineBalanceUC := NewGetSlotMachineBalanceUseCase(slotRepo)

	ctx := context.WithValue(context.Background(), contextkeys.ContextKeyUserID, "admin")
	ctx = context.WithValue(ctx, contextkeys.ContextKeyIsAdmin, true)

	t.Run("Execute_Success", func(t *testing.T) {
		machine := model.NewSlotMachine("machine1", 1, 10000, 2, "teste")
		err := slotRepo.CreateSlotMachine(ctx, machine)
		assert.NoError(t, err, "Expected no error when creating a slot machine")

		req := &GetSlotMachineBalanceRequest{
			MachineID: "machine1",
		}

		resp, err := getSlotMachineBalanceUC.Execute(ctx, req)

		assert.NoError(t, err, "Expected no error when getting slot machine balance")

		assert.NotNil(t, resp, "Expected a response")
		assert.Equal(t, machine.ID, resp.Machine.ID, "Slot machine ID should match the request")
		assert.Equal(t, machine.Level, resp.Machine.Level, "Slot machine level should match the stored level")
		assert.Equal(t, machine.Balance, resp.Machine.Balance, "Slot machine balance should match the stored balance")
	})

	t.Run("Execute_SlotMachineNotFound", func(t *testing.T) {
		req := &GetSlotMachineBalanceRequest{
			MachineID: "nonexistent_machine",
		}

		resp, err := getSlotMachineBalanceUC.Execute(ctx, req)

		assert.Error(t, err, "Expected an error when getting balance of a non-existent slot machine")
		assert.Equal(t, repository.ErrSlotMachineNotFound, err, "Expected ErrSlotMachineNotFound error")
		assert.Nil(t, resp, "Expected no response when there is an error")
	})
}
