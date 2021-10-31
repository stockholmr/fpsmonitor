package computer

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/guregu/null.v3"
)

var (
	dbCtx context.Context
)

func Setup() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	dbCtx = context.Background()
	return db, nil
}

func TestInstall(t *testing.T) {
	db, err := Setup()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewComputerRepository(db)
	err = repo.Install(dbCtx)
	if err != nil {
		t.Fatal(err)
	}

	var data null.String
	row := db.QueryRowContext(dbCtx, "SELECT sql FROM sqlite_master WHERE name='computers'")
	err = row.Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			t.Fatal("missing table schema")
		}
		t.Fatal(err)
		return
	}

	schemaPattern := `CREATE TABLE computers`

	ok, err := regexp.MatchString(schemaPattern, data.String)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("invalid table schema")
	}
}

func TestCreate(t *testing.T) {
	db, err := Setup()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewComputerRepository(db)
	err = repo.Install(dbCtx)
	if err != nil {
		t.Fatal(err)
	}

	computer := &Computer{
		Name: null.NewString("Test Computer", true),
	}

	id, err := repo.Create(dbCtx, computer)
	if err != nil {
		t.Fatal(err)
	}

	if id != 1 {
		t.Fatal("failed to create computer record")
	}
}

func TestSelect(t *testing.T) {
	db, err := Setup()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewComputerRepository(db)
	err = repo.Install(dbCtx)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		computer := &Computer{
			Name: null.NewString(fmt.Sprintf("Test Computer %d", i), true),
		}

		_, err := repo.Create(dbCtx, computer)
		if err != nil {
			t.Fatal(err)
		}
	}

	comp, err := repo.Select(dbCtx, "Test Computer 3")
	if err != nil {
		t.Fatal(err)
	}

	if comp.Name.String != "Test Computer 3" {
		t.Fatal("failed to retrieve correct record")
	}
}

func TestUpdate(t *testing.T) {
	db, err := Setup()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewComputerRepository(db)
	err = repo.Install(dbCtx)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		computer := &Computer{
			Name: null.NewString(fmt.Sprintf("Test Computer %d", i), true),
		}

		_, err := repo.Create(dbCtx, computer)
		if err != nil {
			t.Fatal(err)
		}
	}

	computer := &Computer{
		ID:   null.IntFrom(3),
		Name: null.NewString("Test Computer 33", true),
	}

	err = repo.Update(dbCtx, computer)
	if err != nil {
		t.Fatal(err)
	}

	comp, err := repo.Select(dbCtx, "Test Computer 33")
	if err != nil {
		t.Fatal(err)
	}

	if comp.Name.String != "Test Computer 33" {
		t.Fatal("failed to update record")
	}
}

func TestDelete(t *testing.T) {
	db, err := Setup()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewComputerRepository(db)
	err = repo.Install(dbCtx)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		computer := &Computer{
			Name: null.NewString(fmt.Sprintf("Test Computer %d", i), true),
		}

		_, err := repo.Create(dbCtx, computer)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = repo.Delete(dbCtx, 3)
	if err != nil {
		t.Fatal(err)
	}

	comp, err := repo.Select(dbCtx, "Test Computer 3")
	if err != nil {
		t.Fatal(err)
	}

	if comp.Deleted.Valid {
		t.Fatal("failed to delete record")
	}
}

func TestList(t *testing.T) {
	db, err := Setup()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewComputerRepository(db)
	err = repo.Install(dbCtx)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		computer := &Computer{
			Name: null.NewString(fmt.Sprintf("Test Computer %d", i), true),
		}

		_, err := repo.Create(dbCtx, computer)
		if err != nil {
			t.Fatal(err)
		}
	}

	computers, err := repo.List(dbCtx, 1, 2)
	if err != nil {
		t.Fatal(err)
	}

	if len(computers) != 2 {
		t.Fatal("expected 2 records")
	}

	if computers[0].ID.Int64 != 2 {
		t.Fatal("expected record with id of 2")
	}

	if computers[1].ID.Int64 != 3 {
		t.Fatal("expected record with id of 3")
	}
}
