package usecase

import (
	"context"
	"slot-machine/internal/domain/model"
	repository_in_memory "slot-machine/internal/infrastructure/repository/in_memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSlotMachineUseCase(t *testing.T) {
	slotRepo := repository_in_memory.NewInMemorySlotMachineRepository()

	createSlotMachineUC := NewCreateSlotMachineUseCase(slotRepo)

	ctx := context.Background()

	t.Run("Execute_Success", func(t *testing.T) {
		req := &CreateSlotMachineRequest{
			ID:      "machine1",
			Level:   1,
			Balance: 10000,
		}

		resp, err := createSlotMachineUC.Execute(ctx, req)

		assert.NoError(t, err, "Expected no error when creating a new slot machine")

		assert.NotNil(t, resp, "Expected a response")
		assert.Equal(t, req.ID, resp.Machine.ID, "Slot machine ID should match the request")
		assert.Equal(t, req.Level, resp.Machine.Level, "Slot machine level should match the request")
		assert.Equal(t, req.Balance, resp.Machine.Balance, "Slot machine balance should match the request")

		storedMachine, err := slotRepo.GetSlotMachine(ctx, "machine1")
		assert.NoError(t, err, "Expected no error when retrieving the created slot machine")
		assert.Equal(t, resp.Machine, *storedMachine, "Stored slot machine should match the response")
	})

	t.Run("Execute_SlotMachineAlreadyExists", func(t *testing.T) {
		initialMachine := model.NewSlotMachine("machine2", 2, 20000)
		err := slotRepo.CreateSlotMachine(ctx, initialMachine)
		assert.NoError(t, err, "Expected no error when initially creating a slot machine")

		req := &CreateSlotMachineRequest{
			ID:      "machine2",
			Level:   3,
			Balance: 30000,
		}

		resp, err := createSlotMachineUC.Execute(ctx, req)

		assert.Error(t, err, "Expected an error when creating a slot machine that already exists")
		assert.Equal(t, ErrSlotMachineAlreadyExists, err, "Expected ErrSlotMachineAlreadyExists error")
		assert.Nil(t, resp, "Expected no response when there is an error")

		storedMachine, err := slotRepo.GetSlotMachine(ctx, "machine2")
		assert.NoError(t, err, "Expected no error when retrieving the existing slot machine")
		assert.Equal(t, initialMachine.Balance, storedMachine.Balance, "Slot machine balance should remain unchanged")
		assert.Equal(t, initialMachine.Level, storedMachine.Level, "Slot machine level should remain unchanged")
	})
}
