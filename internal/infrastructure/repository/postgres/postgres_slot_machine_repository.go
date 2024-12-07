package repository_postgres

import (
	"context"
	"slot-machine/internal/domain/model"
	"slot-machine/internal/domain/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSlotMachineRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresSlotMachineRepository(pool *pgxpool.Pool) *PostgresSlotMachineRepository {
	return &PostgresSlotMachineRepository{
		pool: pool,
	}
}

func (r *PostgresSlotMachineRepository) GetSlotMachine(ctx context.Context, id string) (*model.SlotMachine, error) {
	var (
		slotID         string
		level          int
		balance        int
		initialBalance int
		multipleGain   int
		description    string
	)

	row := r.pool.QueryRow(ctx, `
		SELECT id, level, balance, initial_balance, multiple_gain, description
		FROM slot_machines
		WHERE id = $1
	`, id)

	err := row.Scan(&slotID, &level, &balance, &initialBalance, &multipleGain, &description)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrSlotMachineNotFound
		}
		return nil, err
	}

	sm := model.NewSlotMachine(
		slotID,
		level,
		balance,
		multipleGain,
		description,
	)

	return sm, nil
}

func (r *PostgresSlotMachineRepository) UpdateSlotMachine(ctx context.Context, machine *model.SlotMachine) error {
	commandTag, err := r.pool.Exec(ctx, `
		UPDATE slot_machines
		SET level = $1, balance = $2, initial_balance = $3, multiple_gain = $4, description = $5
		WHERE id = $6
	`, machine.Level, machine.Balance, machine.InitialBalance, machine.MultipleGain, machine.Description, machine.ID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return repository.ErrSlotMachineNotFound
	}
	return nil
}

func (r *PostgresSlotMachineRepository) CreateSlotMachine(ctx context.Context, machine *model.SlotMachine) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO slot_machines (id, level, balance, initial_balance, multiple_gain, description)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, machine.ID, machine.Level, machine.Balance, machine.InitialBalance, machine.MultipleGain, machine.Description)
	if err != nil {
		if err.Error() == "duplicate key value violates unique constraint" {
			return repository.ErrSlotMachineExists
		}
		return err
	}
	return nil
}
