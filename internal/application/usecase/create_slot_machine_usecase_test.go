package usecase

import (
	"context"
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
			Level:        1,
			Balance:      10000,
			MultipleGain: 3,
			Description:  "teste",
		}

		resp, err := createSlotMachineUC.Execute(ctx, req)

		assert.NoError(t, err, "Expected no error when creating a new slot machine")

		assert.NotNil(t, resp, "Expected a response")
		assert.Equal(t, req.MultipleGain, resp.Machine.MultipleGain, "Slot machine multiple gain should match the request")
		assert.Equal(t, req.Level, resp.Machine.Level, "Slot machine level should match the request")
		assert.Equal(t, req.Balance, resp.Machine.Balance, "Slot machine balance should match the request")
		assert.Equal(t, req.Description, resp.Machine.Description, "Slot machine description should match the request")

		storedMachine, err := slotRepo.GetSlotMachine(ctx, resp.Machine.ID)
		assert.NoError(t, err, "Expected no error when retrieving the created slot machine")
		assert.Equal(t, resp.Machine, *storedMachine, "Stored slot machine should match the response")
	})
}
