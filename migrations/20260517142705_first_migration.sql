CREATE TABLE contracts (
  id          BIGSERIAL PRIMARY KEY,
  uuid        uuid NOT NULL UNIQUE,
  name        text NOT NULL UNIQUE,
  owner       text NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE contract_versions (
  id            BIGSERIAL PRIMARY KEY,
  uuid          uuid NOT NULL UNIQUE,
  contract_id   BIGINT NOT NULL REFERENCES contracts(id),
  version       int  NOT NULL,
  checksum      text NOT NULL,
  raw_payload   jsonb NOT NULL,
  created_at    timestamptz NOT NULL DEFAULT now(),

  UNIQUE (contract_id, version),
  UNIQUE (contract_id, checksum)
);

CREATE INDEX ON contract_versions (contract_id, version DESC);

CREATE TABLE resources (
  id            BIGSERIAL PRIMARY KEY,
  uuid          uuid NOT NULL UNIQUE,
  contract_id   BIGINT NOT NULL REFERENCES contracts(id),
  direction     text NOT NULL,
  kind          text NOT NULL,
  provider      text,
  endpoint      text NOT NULL,
  method        text NOT NULL,
  status_code   text,
  provider_hash text NOT NULL,
  consumer_hash text,
  created_at    timestamptz NOT NULL DEFAULT now(),

  CHECK ((direction = 'provides') = (consumer_hash IS NULL))
);

CREATE INDEX ON resources (contract_id);
CREATE UNIQUE INDEX ON resources (provider_hash) WHERE direction = 'provides';
CREATE UNIQUE INDEX ON resources (consumer_hash) WHERE direction = 'consumes';
CREATE INDEX ON resources (provider_hash) WHERE direction = 'consumes';

CREATE TABLE properties (
  id           BIGSERIAL PRIMARY KEY,
  uuid         uuid NOT NULL UNIQUE,
  resource_id  BIGINT NOT NULL REFERENCES resources(id),
  path         text NOT NULL,
  created_at   timestamptz NOT NULL DEFAULT now(),

  UNIQUE (resource_id, path)
);

CREATE INDEX ON properties (resource_id);

CREATE TABLE property_versions (
  id                   BIGSERIAL PRIMARY KEY,
  uuid                 uuid NOT NULL UNIQUE,
  property_id          BIGINT NOT NULL REFERENCES properties(id),
  contract_version_id  BIGINT NOT NULL REFERENCES contract_versions(id),
  type                 text,
  optional             boolean,
  change               text NOT NULL CHECK (change IN ('added', 'modified', 'removed')),
  created_at           timestamptz NOT NULL DEFAULT now(),

  UNIQUE (property_id, contract_version_id)
);

CREATE INDEX ON property_versions (property_id, contract_version_id);
CREATE INDEX ON property_versions (contract_version_id);
