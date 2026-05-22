package integration_test

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBAssertions struct {
	suite *Suite
}

func newDBAssertions(s *Suite) *DBAssertions {
	return &DBAssertions{suite: s}
}

func (d *DBAssertions) pool() *pgxpool.Pool {
	return d.suite.Pool
}

func (d *DBAssertions) AssertContract(name, owner string) int64 {
	ctx := context.Background()
	var id int64
	var actualOwner string
	err := d.pool().QueryRow(ctx,
		`SELECT id, owner FROM contracts WHERE name = $1`, name,
	).Scan(&id, &actualOwner)
	d.suite.Require().NoError(err)
	d.suite.Require().Equal(owner, actualOwner)
	return id
}

func (d *DBAssertions) AssertNoContracts() {
	ctx := context.Background()
	var count int
	err := d.pool().QueryRow(ctx, `SELECT COUNT(*) FROM contracts`).Scan(&count)
	d.suite.Require().NoError(err)
	d.suite.Require().Equal(0, count)
}

func (d *DBAssertions) AssertContractNotPersisted(name string) {
	ctx := context.Background()
	var exists bool
	err := d.pool().QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM contracts WHERE name = $1)`, name,
	).Scan(&exists)
	d.suite.Require().NoError(err)
	d.suite.Require().Falsef(exists, "contract %q should not be persisted", name)
}

func (d *DBAssertions) AssertContractVersion(contractID int64, version int, checksum string) int64 {
	ctx := context.Background()
	var id int64
	var actualChecksum string
	var rawPayloadLen int
	err := d.pool().QueryRow(ctx,
		`SELECT id, checksum, LENGTH(raw_payload::text)
		 FROM contract_versions WHERE contract_id = $1 AND version = $2`,
		contractID, version,
	).Scan(&id, &actualChecksum, &rawPayloadLen)
	d.suite.Require().NoError(err)
	d.suite.Require().Equal(checksum, actualChecksum)
	d.suite.Require().Positive(rawPayloadLen)
	return id
}

func (d *DBAssertions) AssertContractVersionCount(contractID int64, expected int) {
	ctx := context.Background()
	var count int
	err := d.pool().QueryRow(ctx,
		`SELECT COUNT(*) FROM contract_versions WHERE contract_id = $1`, contractID,
	).Scan(&count)
	d.suite.Require().NoError(err)
	d.suite.Require().Equal(expected, count)
}

func (d *DBAssertions) AssertResource(contractID int64, direction, kind, endpoint, method, statusCode string) int64 {
	ctx := context.Background()
	var id int64
	err := d.pool().QueryRow(ctx,
		`SELECT id FROM resources
		 WHERE contract_id = $1 AND direction = $2 AND kind = $3
		   AND endpoint = $4 AND method = $5 AND status_code IS NOT DISTINCT FROM $6`,
		contractID, direction, kind, endpoint, method, nullable(statusCode),
	).Scan(&id)
	d.suite.Require().NoError(err)
	return id
}

func (d *DBAssertions) AssertResourceCount(contractID int64, expected int) {
	ctx := context.Background()
	var count int
	err := d.pool().QueryRow(ctx,
		`SELECT COUNT(*) FROM resources WHERE contract_id = $1`, contractID,
	).Scan(&count)
	d.suite.Require().NoError(err)
	d.suite.Require().Equal(expected, count)
}

func (d *DBAssertions) AssertPropertyCount(resourceID int64, expected int) {
	ctx := context.Background()
	var count int
	err := d.pool().QueryRow(ctx,
		`SELECT COUNT(*) FROM properties WHERE resource_id = $1`, resourceID,
	).Scan(&count)
	d.suite.Require().NoError(err)
	d.suite.Require().Equal(expected, count)
}

func (d *DBAssertions) AssertPropertyVersionChangeCounts(resourceID, contractVersionID int64, added, nonAdded int) {
	ctx := context.Background()
	var actualAdded, actualNonAdded int
	err := d.pool().QueryRow(ctx,
		`SELECT
			COUNT(*) FILTER (WHERE pv.change = 'added'),
			COUNT(*) FILTER (WHERE pv.change != 'added')
		 FROM property_versions pv
		 JOIN properties p ON p.id = pv.property_id
		 WHERE p.resource_id = $1 AND pv.contract_version_id = $2`,
		resourceID, contractVersionID,
	).Scan(&actualAdded, &actualNonAdded)
	d.suite.Require().NoError(err)
	d.suite.Require().Equal(added, actualAdded)
	d.suite.Require().Equal(nonAdded, actualNonAdded)
}

func (d *DBAssertions) AssertPropertyVersionCount(contractVersionID int64, expected int) {
	ctx := context.Background()
	var count int
	err := d.pool().QueryRow(ctx,
		`SELECT COUNT(*) FROM property_versions WHERE contract_version_id = $1`, contractVersionID,
	).Scan(&count)
	d.suite.Require().NoError(err)
	d.suite.Require().Equal(expected, count)
}

func (d *DBAssertions) AssertSinglePropertyVersion(contractVersionID int64, path, change string) {
	ctx := context.Background()
	var actualPath, actualChange string
	err := d.pool().QueryRow(ctx,
		`SELECT p.path, pv.change
		 FROM property_versions pv
		 JOIN properties p ON p.id = pv.property_id
		 WHERE pv.contract_version_id = $1`,
		contractVersionID,
	).Scan(&actualPath, &actualChange)
	d.suite.Require().NoError(err)
	d.suite.Require().Equal(path, actualPath)
	d.suite.Require().Equal(change, actualChange)
}

func (d *DBAssertions) AssertContractTreeCounts(name string, contracts, versions, resources, properties int) {
	ctx := context.Background()
	var c, v, r, p int
	err := d.pool().QueryRow(ctx,
		`SELECT
			(SELECT COUNT(*) FROM contracts WHERE name = $1),
			(SELECT COUNT(*) FROM contract_versions cv JOIN contracts c ON c.id = cv.contract_id WHERE c.name = $1),
			(SELECT COUNT(*) FROM resources r JOIN contracts c ON c.id = r.contract_id WHERE c.name = $1),
			(SELECT COUNT(*) FROM properties p JOIN resources r ON r.id = p.resource_id JOIN contracts c ON c.id = r.contract_id WHERE c.name = $1)`,
		name,
	).Scan(&c, &v, &r, &p)
	d.suite.Require().NoError(err)
	d.suite.Require().Equal(contracts, c)
	d.suite.Require().Equal(versions, v)
	d.suite.Require().Equal(resources, r)
	d.suite.Require().Equal(properties, p)
}

func nullable(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
