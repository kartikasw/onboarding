CREATE TABLE IF NOT EXISTS users (
  id bigserial NOT NULL,
  uuid uuid NOT NULL UNIQUE DEFAULT gen_random_uuid(),
  email varchar(50) NOT NULL UNIQUE,
  password varchar NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now()),
  updated_at timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT user__pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS user__uuid__idx ON users USING BTREE (uuid);

CREATE INDEX IF NOT EXISTS user__email__idx ON users USING BTREE (email);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();