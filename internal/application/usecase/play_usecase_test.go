package usecase

import (
	"context"
	"math/rand"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
	repository_in_memory "slot-machine/internal/infrastructure/repository/in_memory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayUseCase(t *testing.T) {
	ctx := context.Background()

	playerRepo := repository_in_memory.NewInMemoryPlayerRepository()
	slotRepo := repository_in_memory.NewInMemorySlotMachineRepository()

	playUC := NewPlayUseCase(playerRepo, slotRepo)

	// Cria um RNG com seed fixa para testes
	fixedSeed := int64(42)
	playUC.rng = rand.New(rand.NewSource(fixedSeed))

	player := &model.Player{
		ID:      "player1",
		Balance: 1000,
	}
	err := playerRepo.CreatePlayer(ctx, player)
	assert.NoError(t, err, "Erro ao criar jogador para testes")

	machine := &model.SlotMachine{
		ID:           "machine1",
		MultipleGain: 2,
		Balance:      5000,
		Permutations: [][3]string{
			{"A", "B", "C"}, // sem vitória
			{"A", "A", "A"}, // vitória
		},
	}
	err = slotRepo.CreateSlotMachine(ctx, machine)
	assert.NoError(t, err, "Erro ao criar máquina de slot para testes")

	t.Run("Execute_Success_Win", func(t *testing.T) {
		// Ajusta o RNG para escolher a segunda permutação (vitória)
		playUC.rng = rand.New(rand.NewSource(1)) // Escolhe index 1

		req := &PlayRequest{
			PlayerID:  "player1",
			MachineID: "machine1",
			AmountBet: 100,
		}

		resp, err := playUC.Execute(ctx, req)

		assert.NoError(t, err, "Esperava-se nenhuma erro em jogada com vitória")
		assert.NotNil(t, resp, "Esperava-se uma resposta")

		assert.True(t, resp.Win, "Esperava-se que o jogador ganhasse")
		assert.Equal(t, [3]string{"A", "A", "A"}, resp.Result, "Esperava-se que todos os símbolos fossem 'A'")
		assert.Equal(t, 1000+100*2, resp.PlayerBalance, "Saldo do jogador deveria ter sido incrementado corretamente")
		assert.Equal(t, 5000-100*2, resp.SlotMachineBalance, "Saldo da máquina de slot deveria ter sido decrementado corretamente")

		// Verifica se os saldos foram atualizados nos repositórios
		updatedPlayer, err := playerRepo.GetPlayer(ctx, "player1")
		assert.NoError(t, err, "Esperava-se encontrar o jogador após atualização")
		assert.Equal(t, 1200, updatedPlayer.Balance, "Saldo do jogador deveria ser 1200")

		updatedMachine, err := slotRepo.GetSlotMachine(ctx, "machine1")
		assert.NoError(t, err, "Esperava-se encontrar a máquina de slot após atualização")
		assert.Equal(t, 4800, updatedMachine.Balance, "Saldo da máquina de slot deveria ser 4800")
	})

	t.Run("Execute_Success_Lose", func(t *testing.T) {
		player.Balance = 1000
		machine.Balance = 5000
		err := playerRepo.UpdatePlayer(ctx, player)
		assert.NoError(t, err, "Erro ao resetar saldo do jogador")
		err = slotRepo.UpdateSlotMachine(ctx, machine)
		assert.NoError(t, err, "Erro ao resetar saldo da máquina de slot")

		// Ajusta o RNG para escolher a primeira permutação (sem vitória)
		playUC.rng = rand.New(rand.NewSource(0)) // Escolhe index 0

		req := &PlayRequest{
			PlayerID:  "player1",
			MachineID: "machine1",
			AmountBet: 100,
		}

		resp, err := playUC.Execute(ctx, req)

		assert.NoError(t, err, "Esperava-se nenhuma erro em jogada sem vitória")
		assert.NotNil(t, resp, "Esperava-se uma resposta")

		assert.False(t, resp.Win, "Esperava-se que o jogador perdesse")
		assert.Equal(t, [3]string{"A", "B", "C"}, resp.Result, "Esperava-se os símbolos 'A', 'B', 'C'")
		assert.Equal(t, 900, resp.PlayerBalance, "Saldo do jogador deveria ter sido decrementado corretamente")
		assert.Equal(t, 5100, resp.SlotMachineBalance, "Saldo da máquina de slot deveria ter sido incrementado corretamente")

		// Verifica se os saldos foram atualizados nos repositórios
		updatedPlayer, err := playerRepo.GetPlayer(ctx, "player1")
		assert.NoError(t, err, "Esperava-se encontrar o jogador após atualização")
		assert.Equal(t, 900, updatedPlayer.Balance, "Saldo do jogador deveria ser 900")

		updatedMachine, err := slotRepo.GetSlotMachine(ctx, "machine1")
		assert.NoError(t, err, "Esperava-se encontrar a máquina de slot após atualização")
		assert.Equal(t, 5100, updatedMachine.Balance, "Saldo da máquina de slot deveria ser 5100")
	})

	t.Run("Execute_InsufficientBalance", func(t *testing.T) {
		playerInsufficient := &model.Player{
			ID:      "player2",
			Balance: 50,
		}
		err := playerRepo.CreatePlayer(ctx, playerInsufficient)
		assert.NoError(t, err, "Erro ao criar jogador com saldo insuficiente")

		req := &PlayRequest{
			PlayerID:  "player2",
			MachineID: "machine1",
			AmountBet: 100,
		}

		resp, err := playUC.Execute(ctx, req)

		assert.Error(t, err, "Esperava-se erro devido ao saldo insuficiente")
		assert.Equal(t, ErrInsufficientBalance, err, "Esperava-se o erro ErrInsufficientBalance")
		assert.Nil(t, resp, "Esperava-se nenhuma resposta quando há erro")

		updatedPlayer, err := playerRepo.GetPlayer(ctx, "player2")
		assert.NoError(t, err, "Esperava-se encontrar o jogador após tentativa de jogada")
		assert.Equal(t, 50, updatedPlayer.Balance, "Saldo do jogador deveria permanecer inalterado")
	})

	t.Run("Execute_SlotMachineNotFound", func(t *testing.T) {
		req := &PlayRequest{
			PlayerID:  "player1",
			MachineID: "nonexistent_machine",
			AmountBet: 100,
		}

		resp, err := playUC.Execute(ctx, req)

		assert.Error(t, err, "Esperava-se erro quando a máquina de slot não é encontrada")
		assert.Equal(t, ErrSlotMachineNotFound, err, "Esperava-se o erro ErrSlotMachineNotFound")
		assert.Nil(t, resp, "Esperava-se nenhuma resposta quando há erro")
	})

	t.Run("Execute_PlayerNotFound", func(t *testing.T) {
		req := &PlayRequest{
			PlayerID:  "nonexistent_player",
			MachineID: "machine1",
			AmountBet: 100,
		}

		resp, err := playUC.Execute(ctx, req)

		assert.Error(t, err, "Esperava-se erro quando o jogador não é encontrado")
		assert.Equal(t, repository.ErrPlayerNotFound, err, "Esperava-se o erro ErrPlayerNotFound")
		assert.Nil(t, resp, "Esperava-se nenhuma resposta quando há erro")
	})
}
