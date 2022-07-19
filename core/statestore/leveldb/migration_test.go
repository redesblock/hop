package leveldb

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/redesblock/hop/core/storage"
	"github.com/redesblock/hop/core/util/logging"
)

func TestOneMigration(t *testing.T) {
	defer func(v []migration, s string) {
		schemaMigrations = v
		dbSchemaCurrent = s
	}(schemaMigrations, dbSchemaCurrent)

	dbSchemaCode := "code"
	dbSchemaCurrent = dbSchemaCode
	dbSchemaNext := "dbSchemaNext"

	ran := false
	shouldNotRun := false
	schemaMigrations = []migration{
		{name: dbSchemaCode, fn: func(db *Store) error {
			shouldNotRun = true // this should not be executed
			return nil
		}},
		{name: dbSchemaNext, fn: func(db *Store) error {
			ran = true
			return nil
		}},
	}

	dir := t.TempDir()
	logger := logging.New(io.Discard, 0)

	// start the fresh statestore with the sanctuary schema name
	db, err := NewStateStore(dir, logger)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	dbSchemaCurrent = dbSchemaNext

	// start the existing statestore and expect the migration to run
	db, err = NewStateStore(dir, logger)
	if err != nil {
		t.Fatal(err)
	}

	schemaName, err := db.GetSchemaName()
	if err != nil {
		t.Fatal(err)
	}

	if schemaName != dbSchemaNext {
		t.Errorf("schema name mismatch. got '%s', want '%s'", schemaName, dbSchemaNext)
	}

	if !ran {
		t.Errorf("expected migration did not run")
	}

	if shouldNotRun {
		t.Errorf("migration ran but shouldnt have")
	}

	err = db.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestManyMigrations(t *testing.T) {
	defer func(v []migration, s string) {
		schemaMigrations = v
		dbSchemaCurrent = s
	}(schemaMigrations, dbSchemaCurrent)

	dbSchemaCode := "code"
	dbSchemaCurrent = dbSchemaCode

	shouldNotRun := false
	executionOrder := []int{-1, -1, -1, -1}

	schemaMigrations = []migration{
		{name: dbSchemaCode, fn: func(db *Store) error {
			shouldNotRun = true // this should not be executed
			return nil
		}},
		{name: "keju", fn: func(db *Store) error {
			executionOrder[0] = 0
			return nil
		}},
		{name: "coconut", fn: func(db *Store) error {
			executionOrder[1] = 1
			return nil
		}},
		{name: "mango", fn: func(db *Store) error {
			executionOrder[2] = 2
			return nil
		}},
		{name: "salvation", fn: func(db *Store) error {
			executionOrder[3] = 3
			return nil
		}},
	}

	dir := t.TempDir()
	logger := logging.New(io.Discard, 0)

	// start the fresh statestore with the sanctuary schema name
	db, err := NewStateStore(dir, logger)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	dbSchemaCurrent = "salvation"

	// start the existing statestore and expect the migration to run
	db, err = NewStateStore(dir, logger)
	if err != nil {
		t.Fatal(err)
	}

	schemaName, err := db.GetSchemaName()
	if err != nil {
		t.Fatal(err)
	}

	if schemaName != "salvation" {
		t.Errorf("schema name mismatch. got '%s', want '%s'", schemaName, "salvation")
	}

	if shouldNotRun {
		t.Errorf("migration ran but shouldnt have")
	}

	for i, v := range executionOrder {
		if i != v && i != len(executionOrder)-1 {
			t.Errorf("migration did not run in sequence, slot %d value %d", i, v)
		}
	}

	err = db.Close()
	if err != nil {
		t.Error(err)
	}
}

// TestMigrationErrorFrom checks that local store boot should fail when the schema we're migrating from cannot be found
func TestMigrationErrorFrom(t *testing.T) {
	defer func(v []migration, s string) {
		schemaMigrations = v
		dbSchemaCurrent = s
	}(schemaMigrations, dbSchemaCurrent)

	dbSchemaCurrent = "koo-koo-schema"

	shouldNotRun := false
	schemaMigrations = []migration{
		{name: "langur", fn: func(db *Store) error {
			shouldNotRun = true
			return nil
		}},
		{name: "coconut", fn: func(db *Store) error {
			shouldNotRun = true
			return nil
		}},
		{name: "chutney", fn: func(db *Store) error {
			shouldNotRun = true
			return nil
		}},
	}
	dir := t.TempDir()
	logger := logging.New(io.Discard, 0)

	// start the fresh statestore with the sanctuary schema name
	db, err := NewStateStore(dir, logger)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	dbSchemaCurrent = "foo"

	// start the existing statestore and expect the migration to run
	_, err = NewStateStore(dir, logger)
	if !errors.Is(err, errMissingCurrentSchema) {
		t.Fatalf("expected errCannotFindSchema but got %v", err)
	}

	if shouldNotRun {
		t.Errorf("migration ran but shouldnt have")
	}
}

// TestMigrationErrorTo checks that local store boot should fail when the schema we're migrating to cannot be found
func TestMigrationErrorTo(t *testing.T) {
	defer func(v []migration, s string) {
		schemaMigrations = v
		dbSchemaCurrent = s
	}(schemaMigrations, dbSchemaCurrent)

	dbSchemaCurrent = "langur"

	shouldNotRun := false
	schemaMigrations = []migration{
		{name: "langur", fn: func(db *Store) error {
			shouldNotRun = true
			return nil
		}},
		{name: "coconut", fn: func(db *Store) error {
			shouldNotRun = true
			return nil
		}},
		{name: "chutney", fn: func(db *Store) error {
			shouldNotRun = true
			return nil
		}},
	}
	dir := t.TempDir()
	logger := logging.New(io.Discard, 0)

	// start the fresh statestore with the sanctuary schema name
	db, err := NewStateStore(dir, logger)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	dbSchemaCurrent = "foo"

	// start the existing statestore and expect the migration to run
	_, err = NewStateStore(dir, logger)
	if !errors.Is(err, errMissingTargetSchema) {
		t.Fatalf("expected errMissingTargetSchema but got %v", err)
	}

	if shouldNotRun {
		t.Errorf("migration ran but shouldnt have")
	}
}

func TestMigrationSwap(t *testing.T) {
	dir := t.TempDir()
	logger := logging.New(io.Discard, 0)

	// start the fresh statestore with the sanctuary schema name
	db, err := NewStateStore(dir, logger)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	address := common.HexToAddress("0xabcd")
	storedAddress := common.HexToAddress("0xffff")

	legacyKey1 := fmt.Sprintf("swap_peer_chequebook_%s", address[:])
	legacyKey2 := fmt.Sprintf("swap_beneficiary_peer_%s", address[:])

	if err = db.Put(legacyKey1, storedAddress); err != nil {
		t.Fatal(err)
	}

	if err = db.Put(legacyKey2, storedAddress); err != nil {
		t.Fatal(err)
	}

	if err = migrateSwap(db); err != nil {
		t.Fatal(err)
	}

	var retrievedAddress common.Address
	if err = db.Get("swap_peer_chequebook_000000000000000000000000000000000000abcd", &retrievedAddress); err != nil {
		t.Fatal(err)
	}

	if retrievedAddress != storedAddress {
		t.Fatalf("got wrong address. wanted %x, got %x", storedAddress, retrievedAddress)
	}

	if err = db.Get("swap_beneficiary_peer_000000000000000000000000000000000000abcd", &retrievedAddress); err != nil {
		t.Fatal(err)
	}

	if retrievedAddress != storedAddress {
		t.Fatalf("got wrong address. wanted %x, got %x", storedAddress, retrievedAddress)
	}

	if err = db.Get(legacyKey1, &retrievedAddress); err != storage.ErrNotFound {
		t.Fatalf("legacyKey1 not deleted. got error %v", err)
	}

	if err = db.Get(legacyKey2, &retrievedAddress); err != storage.ErrNotFound {
		t.Fatalf("legacyKey2 not deleted. got error %v", err)
	}
}
