package integration_test

import (
	"context"
)

func (suite *Suite) AssertContract(name, owner string) int64 {
	ctx := context.Background()
	var id int64
	var actualOwner string
	err := suite.Pool.QueryRow(ctx,
		`SELECT id, owner FROM contracts WHERE name = $1`, name,
	).Scan(&id, &actualOwner)
	suite.Require().NoError(err)
	suite.Require().Equal(owner, actualOwner)
	return id
}

func (suite *Suite) AssertNoContracts() {
	ctx := context.Background()
	var count int
	err := suite.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM contracts`).Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(0, count)
}

func (suite *Suite) AssertContractVersion(contractID int64, version int, checksum string) int64 {
	ctx := context.Background()
	var id int64
	var actualChecksum string
	var rawPayloadLen int
	err := suite.Pool.QueryRow(ctx,
		`SELECT id, checksum, LENGTH(raw_payload::text)
		 FROM contract_versions WHERE contract_id = $1 AND version = $2`,
		contractID, version,
	).Scan(&id, &actualChecksum, &rawPayloadLen)
	suite.Require().NoError(err)
	suite.Require().Equal(checksum, actualChecksum)
	suite.Require().Positive(rawPayloadLen)
	return id
}

func (suite *Suite) AssertContractVersionCount(contractID int64, expected int) {
	ctx := context.Background()
	var count int
	err := suite.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM contract_versions WHERE contract_id = $1`, contractID,
	).Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(expected, count)
}

func (suite *Suite) AssertResource(contractID int64, direction, kind, endpoint, method, statusCode string) int64 {
	ctx := context.Background()
	var id int64
	err := suite.Pool.QueryRow(ctx,
		`SELECT id FROM resources
		 WHERE contract_id = $1 AND direction = $2 AND kind = $3
		   AND endpoint = $4 AND method = $5 AND status_code IS NOT DISTINCT FROM $6`,
		contractID, direction, kind, endpoint, method, nullable(statusCode),
	).Scan(&id)
	suite.Require().NoError(err)
	return id
}

func (suite *Suite) AssertResourceCount(contractID int64, expected int) {
	ctx := context.Background()
	var count int
	err := suite.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM resources WHERE contract_id = $1`, contractID,
	).Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(expected, count)
}

func (suite *Suite) AssertPropertyCount(resourceID int64, expected int) {
	ctx := context.Background()
	var count int
	err := suite.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM properties WHERE resource_id = $1`, resourceID,
	).Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(expected, count)
}

func (suite *Suite) AssertPropertyVersionChangeCounts(resourceID, contractVersionID int64, added, nonAdded int) {
	ctx := context.Background()
	var actualAdded, actualNonAdded int
	err := suite.Pool.QueryRow(ctx,
		`SELECT
			COUNT(*) FILTER (WHERE pv.change = 'added'),
			COUNT(*) FILTER (WHERE pv.change != 'added')
		 FROM property_versions pv
		 JOIN properties p ON p.id = pv.property_id
		 WHERE p.resource_id = $1 AND pv.contract_version_id = $2`,
		resourceID, contractVersionID,
	).Scan(&actualAdded, &actualNonAdded)
	suite.Require().NoError(err)
	suite.Require().Equal(added, actualAdded)
	suite.Require().Equal(nonAdded, actualNonAdded)
}

func (suite *Suite) AssertPropertyVersionCount(contractVersionID int64, expected int) {
	ctx := context.Background()
	var count int
	err := suite.Pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM property_versions WHERE contract_version_id = $1`, contractVersionID,
	).Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(expected, count)
}

func (suite *Suite) AssertSinglePropertyVersion(contractVersionID int64, path, change string) {
	ctx := context.Background()
	var actualPath, actualChange string
	err := suite.Pool.QueryRow(ctx,
		`SELECT p.path, pv.change
		 FROM property_versions pv
		 JOIN properties p ON p.id = pv.property_id
		 WHERE pv.contract_version_id = $1`,
		contractVersionID,
	).Scan(&actualPath, &actualChange)
	suite.Require().NoError(err)
	suite.Require().Equal(path, actualPath)
	suite.Require().Equal(change, actualChange)
}

func (suite *Suite) AssertContractTreeCounts(name string, contracts, versions, resources, properties int) {
	ctx := context.Background()
	var c, v, r, p int
	err := suite.Pool.QueryRow(ctx,
		`SELECT
			(SELECT COUNT(*) FROM contracts WHERE name = $1),
			(SELECT COUNT(*) FROM contract_versions cv JOIN contracts c ON c.id = cv.contract_id WHERE c.name = $1),
			(SELECT COUNT(*) FROM resources r JOIN contracts c ON c.id = r.contract_id WHERE c.name = $1),
			(SELECT COUNT(*) FROM properties p JOIN resources r ON r.id = p.resource_id JOIN contracts c ON c.id = r.contract_id WHERE c.name = $1)`,
		name,
	).Scan(&c, &v, &r, &p)
	suite.Require().NoError(err)
	suite.Require().Equal(contracts, c)
	suite.Require().Equal(versions, v)
	suite.Require().Equal(resources, r)
	suite.Require().Equal(properties, p)
}

func nullable(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
