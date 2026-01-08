BEGIN; 

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TYPE account_role AS ENUM ('USER', 'WORKER');

CREATE TABLE accounts (
    account_id    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email         TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role          account_role NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_accounts_email ON accounts (lower(email));
CREATE INDEX idx_accounts_role ON accounts (role);
CREATE INDEX idx_accounts_role_active ON accounts (role, is_active);

CREATE TABLE zones (
    zone_id   SMALLINT PRIMARY KEY,
    title     TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE addresses (
    address_id        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    zone_id           SMALLINT NOT NULL REFERENCES zones(zone_id) ON DELETE RESTRICT,
    street            TEXT NOT NULL,
    house_number      TEXT NOT NULL,
    entrance          TEXT,
    floor_number      INTEGER NOT NULL,
    apartment_number  INTEGER NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_addresses_zone_id ON addresses(zone_id);

CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id  UUID NOT NULL UNIQUE REFERENCES accounts(account_id) ON DELETE CASCADE,

    first_name TEXT NOT NULL,
    surname TEXT NOT NULL,
    last_name TEXT,

    address_id UUID NOT NULL REFERENCES addresses(address_id) ON DELETE RESTRICT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_created_at ON users (created_at);
CREATE INDEX idx_users_address_id ON users(address_id);

CREATE TABLE workers (
    worker_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id  UUID NOT NULL UNIQUE REFERENCES accounts(account_id) ON DELETE CASCADE,

    zone_id    SMALLINT NOT NULL REFERENCES zones(zone_id) ON DELETE RESTRICT,

    first_name TEXT NOT NULL, 
    surname TEXT NOT NULL,
    last_name TEXT,
    
    is_active BOOLEAN NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_workers_zone_active ON workers(zone_id, is_active);
CREATE INDEX idx_workers_created_at ON workers (created_at);
CREATE INDEX idx_workers_updated_at ON workers (updated_at);

CREATE TYPE task_status AS ENUM ('OPEN', 'IN_PROGRESS', 'DONE', 'CANCELED');

CREATE TABLE tasks (
    task_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    client_id UUID NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,

    address_id  UUID NOT NULL REFERENCES addresses(address_id) ON DELETE RESTRICT,
    worker_id   UUID REFERENCES workers(worker_id) ON DELETE SET NULL,

    status task_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
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

CREATE INDEX idx_tasks_open_created_at
  ON tasks(created_at)
  WHERE status = 'OPEN';

COMMIT;