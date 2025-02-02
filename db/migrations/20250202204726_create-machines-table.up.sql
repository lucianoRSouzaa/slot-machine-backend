CREATE TABLE IF NOT EXISTS slot_machines (
    id VARCHAR(36) PRIMARY KEY,
    level INTEGER NOT NULL,
    balance INTEGER NOT NULL,
    initial_balance INTEGER NOT NULL,
    multiple_gain INTEGER NOT NULL,
    description TEXT NOT NULL
);
