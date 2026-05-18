# Merge repository suite into integration suite

## Context

Today there are two test directories doing the same thing — running a real
Postgres testcontainer, applying migrations, exercising the broker — but
split arbitrarily:

- `tests/integration/` — HTTP-level tests via Fiber `.Test()`. Asserts
  status code + JSON body only.
- `tests/repository/` — direct repository calls. Asserts via `Suite`
  `Assert*` methods that run SELECTs against the same DB.

The split means HTTP tests cannot prove the request *actually persisted
anything*, and the assertion vocabulary (`AssertContract`,
`AssertContractVersion`, …) is invisible to them. Merging the two into one
suite — in `tests/integration/` — lets every HTTP test follow `Request()`
with the same `Assert*` calls, and lets a single test mix HTTP-driven and
repository-driven verifications.

Three additional rules the user gave:
1. `suite_test.go` is **config only** (Suite struct, `SetupTest`,
   `TearDownTest`, `StartPostgressContainer`, `Request` plumbing).
2. The `Assert*` methods move to their own file.
3. `create_contract` HTTP tests use those `Assert*` methods.

Controller wiring for `Update` stays out of scope — duplicate POST still
returns "already uploaded"; the repository-level `TestUpdate_*` test
continues to cover the diff-persistence path.

## Approach

Consolidate everything under `tests/integration/` (package
`integration_test`). The existing repository suite gets deleted; its
tests, fixtures, and assertion methods move over.

## Files to create

### `tests/integration/assertions_test.go`

All `Assert*` methods on `*Suite`. Moved from the current
`tests/repository/suite_test.go` plus a couple of additions needed by the
HTTP tests:

- `AssertContract(name, owner string) int64`
- `AssertContractVersion(contractID int64, version int, checksum string) int64`
- `AssertContractVersionCount(contractID int64, expected int)`
- `AssertResource(contractID int64, direction, kind, endpoint, method, statusCode string) int64`
- `AssertResourceCount(contractID int64, expected int)` *(new — used by the rich create test)*
- `AssertPropertyCount(resourceID int64, expected int)`
- `AssertPropertyVersionCount(contractVersionID int64, expected int)`
- `AssertPropertyVersionChangeCounts(resourceID, contractVersionID int64, added, nonAdded int)`
- `AssertSinglePropertyVersion(contractVersionID int64, path, change string)`
- `AssertContractTreeCounts(name string, contracts, versions, resources, properties int)`
- `AssertNoContracts()` *(new — used by the missing-name test)*

Plus the `nullable()` helper used by `AssertResource`.

All methods use `suite.Require()` (halt-on-failure) so a missing row stops
the cascade.

### `tests/integration/contract_repository_test.go`

Moved verbatim from `tests/repository/contract_repository_test.go`. Same
five tests, same shape (each calls `Repo.Save` / `Repo.Update` /
`Repo.FindByName` directly, then `Assert*`):

1. `TestSave_PersistsContractTree`
2. `TestSave_RollsBackOnDuplicateInsert`
3. `TestUpdate_PersistsNewFieldAsAddedInNextVersion`
4. `TestFindByName_ReturnsNilWhenMissing`
5. `TestFindByName_ReturnsContractWhenSaved`

Only changes: package becomes `integration_test`.

## Files to modify

### `tests/integration/suite_test.go`

Combine the two suites' setup into one — config only:

```go
type Suite struct {
    suite.Suite
    Pool              *pgxpool.Pool
    Repo              *repository.ContractRepository
    Components        *components.Components
    PostgresContainer *postgres.PostgresContainer
}
```

`SetupTest` wires `Pool = Components.Pool` and
`Repo = repository.NewContractRepository(Pool)`. Request/Response types
and `Request()` helper unchanged. Package becomes `integration_test`. No
`Assert*` methods here.

### `tests/integration/create_contract_test.go`

Keep the `itemsContractPayload` constant. Enhance the three existing HTTP
tests with DB-state assertions, and add two new tests that focus on parts
of the persisted tree:

1. **`TestCreateContract`** — keep the existing HTTP assertions
   (status 200, JSON body). Then assert contract / version (checksum
   computed by parsing the same payload through `dsl.Contract` +
   `ToContractModel().Checksum()`) / one provided REST resource /
   property-version change counts.

2. **`TestCreateContract_MissingName_Returns400`** — keep 400 + JSON.
   Add `suite.AssertNoContracts()` to prove no partial write happened.

3. **`TestCreateContract_DuplicateName_ReturnsAlreadyUploaded`** — keep
   the two HTTP assertions. Add `AssertContract` + `AssertContractVersionCount(.., 1)`
   to prove the short-circuit did not create a second version.

4. **`TestCreateContract_PersistsAllResources`** *(new)* — POST the rich
   payload, then `AssertResourceCount(contractID, 8)`.

5. **`TestCreateContract_PersistsConsumedBillingInvoicesPost`** *(new)* —
   POST the rich payload, then `AssertResource(contractID, "consumes",
   "rest_response", "/invoices", "POST", "201")` — exercises the
   consumed-side persistence (provider/status_code columns).

## Files to delete

- `tests/repository/suite_test.go`
- `tests/repository/contract_repository_test.go`
- `tests/repository/` (directory)

## Verification

From `/Users/alefcastelo/workspace/contracttests/broker`:

```bash
go build ./...
go test ./tests/integration/... -v
go test ./... -v
```

Expected:
- 10 tests under `tests/integration/`: 5 repository-direct + 3 enhanced
  HTTP + 2 new HTTP.
- `tests/repository/` no longer exists.
- Pre-existing non-test packages unchanged.
