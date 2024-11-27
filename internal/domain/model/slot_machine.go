package model

type SlotMachine struct {
	ID             string            `json:"id"`
	Level          int               `json:"level"`
	Balance        int               `json:"balance"`
	InitialBalance int               `json:"initial_balance"`
	Symbols        map[string]string `json:"symbols"`
	Permutations   [][3]string       `json:"permutations"`
	MultipleGain   int               `json:"multiple_gain"`
	Description    string            `json:"description"`
}

func NewSlotMachine(id string, level, balance int, multipleGain int, description string) *SlotMachine {
	sm := &SlotMachine{
		ID:             id,
		Level:          level,
		Balance:        balance,
		InitialBalance: balance,
		Symbols: map[string]string{
			"money_mouth_face": "1F911",
			"cold_face":        "1F976",
			"alien":            "1F47D",
			"heart_on_fire":    "2764",
			"collision":        "1F4A5",
		},
		MultipleGain: multipleGain,
		Description:  description,
	}
	sm.GeneratePermutations()
	return sm
}

func (sm *SlotMachine) GeneratePermutations() {
	perms := [][3]string{}
	symbolKeys := make([]string, 0, len(sm.Symbols))
	for k := range sm.Symbols {
		symbolKeys = append(symbolKeys, k)
	}

	for _, a := range symbolKeys {
		for _, b := range symbolKeys {
			for _, c := range symbolKeys {
				perms = append(perms, [3]string{a, b, c})
			}
		}
	}

	for j := 0; j < sm.Level; j++ {
		for _, sym := range symbolKeys {
			perms = append(perms, [3]string{sym, sym, sym})
		}
	}

	sm.Permutations = perms
}
