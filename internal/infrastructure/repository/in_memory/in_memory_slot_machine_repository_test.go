package repository_in_memory

import (
	"context"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemorySlotMachineRepository(t *testing.T) {
	repo := NewInMemorySlotMachineRepository()

	ctx := context.Background()

	t.Run("CreateSlotMachine_Success", func(t *testing.T) {
		machine := model.NewSlotMachine("machine1", 1, 10000)

		err := repo.CreateSlotMachine(ctx, machine)
		assert.NoError(t, err, "Expected no error on creating slot machine")

		retrievedMachine, err := repo.GetSlotMachine(ctx, "machine1")
		assert.NoError(t, err, "Expected no error on retrieving existing slot machine")
		assert.Equal(t, machine, retrievedMachine, "Retrieved slot machine should match the created slot machine")
	})

	t.Run("CreateSlotMachine_DuplicateID", func(t *testing.T) {
		machine := model.NewSlotMachine("machine1", 2, 20000)

		err := repo.CreateSlotMachine(ctx, machine)
		assert.Error(t, err, "Expected error on creating slot machine with duplicate ID")
		assert.Equal(t, repository.ErrSlotMachineExists, err, "Expected ErrSlotMachineExists error")
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

	t.Run("ConcurrentAccess", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 100
		machineIDs := make([]string, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			machineID := "concurrent_machine_" + strconv.Itoa(i)
			machineIDs[i] = machineID
			go func(id string) {
				defer wg.Done()
				machine := model.NewSlotMachine(id, 1, 1000)
				err := repo.CreateSlotMachine(ctx, machine)
				assert.NoError(t, err, "Expected no error on concurrent slot machine creation")
			}(machineID)
		}

		wg.Wait()

		for _, id := range machineIDs {
			machine, err := repo.GetSlotMachine(ctx, id)
			assert.NoError(t, err, "Expected no error on retrieving concurrently created slot machine")
			assert.Equal(t, 1000, machine.Balance, "Expected slot machine balance to be 1000")
		}
	})
}
