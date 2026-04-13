BEGIN;

CREATE TYPE order_status AS ENUM ('created', 'assigned', 'in_progress', 'completed', 'canceled');

CREATE TABLE orders (
    id             UUID PRIMARY KEY,
    user_id        UUID NOT NULL,
    worker_id      UUID,
    address        TEXT NOT NULL,
    description    TEXT,
    preferred_time TIMESTAMPTZ,

    status         order_status NOT NULL DEFAULT 'created',

    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    assigned_at    TIMESTAMPTZ,
    started_at     TIMESTAMPTZ,
    completed_at   TIMESTAMPTZ,
    canceled_at    TIMESTAMPTZ
);

CREATE INDEX idx_orders_user_id ON orders (user_id);
CREATE INDEX idx_orders_worker_id ON orders (worker_id);
CREATE INDEX idx_orders_status ON orders (status);
CREATE INDEX idx_orders_created_at ON orders (created_at DESC);

COMMIT;
