package repository_in_memory

import (
	"context"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemorySlotMachineRepository(t *testing.T) {
	repo := NewInMemorySlotMachineRepository()

	ctx := context.Background()

	t.Run("CreateSlotMachine_Success", func(t *testing.T) {
		machine := model.NewSlotMachine("machine1", 1, 10000, 2, "teste")

		err := repo.CreateSlotMachine(ctx, machine)
		assert.NoError(t, err, "Expected no error on creating slot machine")

		retrievedMachine, err := repo.GetSlotMachine(ctx, "machine1")
		assert.NoError(t, err, "Expected no error on retrieving existing slot machine")
		assert.Equal(t, machine, retrievedMachine, "Retrieved slot machine should match the created slot machine")
	})

	t.Run("GetSlotMachine_NotFound", func(t *testing.T) {
		_, err := repo.GetSlotMachine(ctx, "nonexistent_machine")
		assert.Error(t, err, "Expected error when retrieving non-existent slot machine")
		assert.Equal(t, repository.ErrSlotMachineNotFound, err, "Expected ErrSlotMachineNotFound error")
	})

	t.Run("UpdateSlotMachine_Success", func(t *testing.T) {
		machine := &model.SlotMachine{
			ID:             "machine1",
			Level:          3,
			Balance:        15000,
			InitialBalance: 10000,
			Symbols: map[string]string{
				"money_mouth_face": "1F911",
				"cold_face":        "1F976",
				"alien":            "1F47D",
				"heart_on_fire":    "2764",
				"collision":        "1F4A5",
			},
			MultipleGain: 5,
		}
		machine.GeneratePermutations()

		err := repo.UpdateSlotMachine(ctx, machine)
		assert.NoError(t, err, "Expected no error on updating existing slot machine")

		retrievedMachine, err := repo.GetSlotMachine(ctx, "machine1")
		assert.NoError(t, err, "Expected no error on retrieving existing slot machine")
		assert.Equal(t, machine, retrievedMachine, "Retrieved slot machine should reflect the updated data")
	})

	t.Run("UpdateSlotMachine_NotFound", func(t *testing.T) {
		machine := &model.SlotMachine{
			ID:             "nonexistent_machine",
			Level:          1,
			Balance:        5000,
			InitialBalance: 5000,
			Symbols:        map[string]string{},
			Permutations:   [][3]string{},
			MultipleGain:   2,
		}

		err := repo.UpdateSlotMachine(ctx, machine)
		assert.Error(t, err, "Expected error on updating non-existent slot machine")
		assert.Equal(t, repository.ErrSlotMachineNotFound, err, "Expected ErrSlotMachineNotFound error")
	})
}
