BEGIN; 

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name TEXT NOT NULL,
    surname TEXT NOT NULL,
    last_name TEXT,
    address_id UUID NOT NULL REFERENCES addresses(address_id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE addresses(
    address_id        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    street TEXT NOT NULL,
    house_number TEXT NOT NULL,
    entrance TEXT,
    floor_number INTEGER NOT NULL,
    apartment_number INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_addresses_user_id ON addresses (user_id);
CREATE INDEX idx_users_created_at ON users (created_at);

CREATE TABLE workers (
    worker_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name TEXT NOT NULL, 
    surname TEXT NOT NULL,
    last_name TEXT,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_workers_is_active ON workers (is_active);
CREATE INDEX idx_workers_created_at ON workers (created_at);
CREATE INDEX idx_workers_updated_at ON workers (updated_at);

CREATE TYPE task_status AS ENUM ('OPEN', 'IN PROGRESS', 'DONE', 'CANCELED');

CREATE TABLE tasks (
    task_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id UUID NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    address_id  UUID NOT NULL REFERENCES addresses(address_id) ON DELETE RESTRICT,
    worker_id   UUID REFERENCES workers(worker_id) ON DELETE SET NULL,
    status task_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    closed_at TIMESTAMPTZ
);

CREATE INDEX idx_tasks_client_id ON tasks (client_id);
CREATE INDEX idx_tasks_worker_id ON tasks (worker_id);

CREATE INDEX idx_tasks_worker_open
  ON tasks (worker_id, created_at)
  WHERE closed_at IS NULL;

CREATE INDEX idx_tasks_open
  ON tasks (created_at)
  WHERE closed_at IS NULL;

COMMIT;