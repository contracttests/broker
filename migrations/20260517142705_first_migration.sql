CREATE TABLE participants (
  id          BIGSERIAL PRIMARY KEY,
  name        text NOT NULL UNIQUE,
  created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE environments (
  id          BIGSERIAL PRIMARY KEY,
  name        text NOT NULL UNIQUE,
  created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE contracts (
  id              BIGSERIAL PRIMARY KEY,
  participant_id  BIGINT NOT NULL REFERENCES participants(id),
  version         text NOT NULL,
  checksum        text NOT NULL,
  raw_payload     text NOT NULL,
  created_at      timestamptz NOT NULL DEFAULT now(),

  UNIQUE (participant_id, version),
  UNIQUE (participant_id, checksum)
);

CREATE INDEX ON contracts (participant_id, version DESC);

CREATE TABLE resources (
  id              BIGSERIAL PRIMARY KEY,
  participant_id  BIGINT NOT NULL REFERENCES participants(id),
  direction       text NOT NULL,
  kind            text NOT NULL,
  provider        text,
  endpoint        text NOT NULL,
  method          text NOT NULL,
  status_code     text,
  provider_hash   text NOT NULL,
  consumer_hash   text,
  created_at      timestamptz NOT NULL DEFAULT now(),

  CHECK ((direction = 'provides') = (consumer_hash IS NULL))
);

CREATE INDEX ON resources (participant_id);
CREATE UNIQUE INDEX ON resources (provider_hash) WHERE direction = 'provides';
CREATE UNIQUE INDEX ON resources (consumer_hash) WHERE direction = 'consumes';
CREATE INDEX ON resources (provider_hash) WHERE direction = 'consumes';

CREATE TABLE properties (
  id           BIGSERIAL PRIMARY KEY,
  resource_id  BIGINT NOT NULL REFERENCES resources(id),
  path         text NOT NULL,
  created_at   timestamptz NOT NULL DEFAULT now(),

  UNIQUE (resource_id, path)
);

CREATE INDEX ON properties (resource_id);

CREATE TABLE property_versions (
  id           BIGSERIAL PRIMARY KEY,
  property_id  BIGINT NOT NULL REFERENCES properties(id),
  contract_id  BIGINT NOT NULL REFERENCES contracts(id),
  type         text,
  optional     boolean,
  change       text NOT NULL CHECK (change IN ('added', 'modified', 'removed')),
  created_at   timestamptz NOT NULL DEFAULT now(),

  UNIQUE (property_id, contract_id)
);

CREATE INDEX ON property_versions (property_id, contract_id);
CREATE INDEX ON property_versions (contract_id);

CREATE TABLE resource_versions (
  id           BIGSERIAL PRIMARY KEY,
  resource_id  BIGINT NOT NULL REFERENCES resources(id),
  contract_id  BIGINT NOT NULL REFERENCES contracts(id),
  change       text NOT NULL CHECK (change IN ('added', 'removed')),
  created_at   timestamptz NOT NULL DEFAULT now(),

  UNIQUE (resource_id, contract_id)
);

CREATE INDEX ON resource_versions (resource_id, contract_id);
CREATE INDEX ON resource_versions (contract_id);
