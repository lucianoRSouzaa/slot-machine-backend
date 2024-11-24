package usecase

import (
	"context"
	"errors"
	"math/rand"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"
	"time"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrSlotMachineNotFound = errors.New("slot machine not found")
)

type PlayUseCase struct {
	PlayerRepo      repository.PlayerRepository
	SlotMachineRepo repository.SlotMachineRepository
	rng             *rand.Rand
}

type PlayRequest struct {
	PlayerID  string `json:"player_id"`
	MachineID string `json:"machine_id"`
	AmountBet int    `json:"amount_bet"`
}

type PlayResponse struct {
	Result             [3]string `json:"result"`
	Win                bool      `json:"win"`
	PlayerBalance      int       `json:"player_balance"`
	SlotMachineBalance int       `json:"slot_machine_balance"`
}

func NewPlayUseCase(playerRepo repository.PlayerRepository, slotRepo repository.SlotMachineRepository) *PlayUseCase {
	return &PlayUseCase{
		PlayerRepo:      playerRepo,
		SlotMachineRepo: slotRepo,
		rng:             rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (uc *PlayUseCase) Execute(ctx context.Context, req *PlayRequest) (*PlayResponse, error) {
	player, err := uc.PlayerRepo.GetPlayer(ctx, req.PlayerID)
	if err != nil {
		return nil, err
	}

	machine, err := uc.SlotMachineRepo.GetSlotMachine(ctx, req.MachineID)
	if err != nil {
		return nil, err
	}

	if player.Balance < req.AmountBet {
		return nil, ErrInsufficientBalance
	}

	result := uc.generateFinalResult(machine)

	win := uc.checkResultUser(result)

	if win {
		player.Balance += req.AmountBet * machine.MultipleGain
		machine.Balance -= req.AmountBet * machine.MultipleGain
	} else {
		player.Balance -= req.AmountBet
		machine.Balance += req.AmountBet
	}

	if err := uc.PlayerRepo.UpdatePlayer(ctx, player); err != nil {
		return nil, err
	}
	if err := uc.SlotMachineRepo.UpdateSlotMachine(ctx, machine); err != nil {
		return nil, err
	}

	return &PlayResponse{
		Result:             result,
		Win:                win,
		PlayerBalance:      player.Balance,
		SlotMachineBalance: machine.Balance,
	}, nil
}

func (uc *PlayUseCase) generateFinalResult(machine *model.SlotMachine) [3]string {
	permutation := machine.Permutations[uc.rng.Intn(len(machine.Permutations))]
	result := permutation

	if len(unique(result[:])) == 3 && uc.rng.Intn(6) >= 2 {
		if uc.rng.Intn(2) == 0 {
			result[1] = result[0]
		} else {
			result[2] = result[1]
		}
	}
	return result
}

func unique(slice []string) []string {
	uniqueMap := make(map[string]struct{})
	for _, item := range slice {
		uniqueMap[item] = struct{}{}
	}
	keys := make([]string, 0, len(uniqueMap))
	for k := range uniqueMap {
		keys = append(keys, k)
	}
	return keys
}

func (uc *PlayUseCase) checkResultUser(result [3]string) bool {
	first := result[0]
	for _, sym := range result[1:] {
		if sym != first {
			return false
		}
	}
	return true
}
