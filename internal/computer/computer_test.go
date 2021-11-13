package computer

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/guregu/null.v3"
)

var (
	dbCtx context.Context
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: "+msg+"\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: unexpected error: %s\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func dbSetup() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	dbCtx = context.Background()
	return db, nil
}

func TestComputerRepositoryInstall(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewComputerRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	var data null.String
	row := db.QueryRowContext(dbCtx, "SELECT sql FROM sqlite_master WHERE name='computers'")
	err = row.Scan(&data)
	ok(t, err)

	schemaPattern := `CREATE TABLE computers`

	matched, err := regexp.MatchString(schemaPattern, data.String)
	ok(t, err)
	assert(t, matched, "invalid table schema", nil)
}

func TestComputerRepositoryCreate(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewComputerRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	computer := &Computer{
		Name: null.NewString("Test Computer", true),
	}

	id, err := repo.Create(dbCtx, computer)
	ok(t, err)

	equals(t, int64(1), id)
}

func TestComputerRepositorySelect(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewComputerRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		computer := &Computer{
			Name: null.NewString(fmt.Sprintf("Test Computer %d", i), true),
		}

		_, err := repo.Create(dbCtx, computer)
		ok(t, err)
	}

	comp, err := repo.Select(dbCtx, "Test Computer 3")
	ok(t, err)
	equals(t, "Test Computer 3", comp.Name.String)
}

func TestComputerRepositoryUpdate(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewComputerRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		computer := &Computer{
			Name: null.NewString(fmt.Sprintf("Test Computer %d", i), true),
		}

		_, err := repo.Create(dbCtx, computer)
		ok(t, err)
	}

	computer := &Computer{
		ID:   null.IntFrom(3),
		Name: null.NewString("Test Computer 33", true),
	}

	err = repo.Update(dbCtx, computer)
	ok(t, err)

	comp, err := repo.Select(dbCtx, "Test Computer 33")
	ok(t, err)
	equals(t, "Test Computer 33", comp.Name.String)
}

func TestComputerRepositoryDelete(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewComputerRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		computer := &Computer{
			Name: null.NewString(fmt.Sprintf("Test Computer %d", i), true),
		}

		_, err := repo.Create(dbCtx, computer)
		ok(t, err)
	}

	err = repo.Delete(dbCtx, 4)
	ok(t, err)

	comp, err := repo.Select(dbCtx, "Test Computer 3")
	ok(t, err)
	equals(t, true, comp.Deleted.Valid)
}

func TestComputerRepositoryList(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewComputerRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		computer := &Computer{
			Name: null.NewString(fmt.Sprintf("Test Computer %d", i), true),
		}

		_, err := repo.Create(dbCtx, computer)
		ok(t, err)
	}

	computers, err := repo.List(dbCtx, 1, 2)
	ok(t, err)

	equals(t, 2, len(computers))
	equals(t, int64(2), computers[0].ID.Int64)
	equals(t, int64(3), computers[1].ID.Int64)
}

func TestNetworkAdapterRepositoryInstall(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewNetworkAdapterRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	var data null.String
	row := db.QueryRowContext(dbCtx, "SELECT sql FROM sqlite_master WHERE name='computer_network_adapters'")
	err = row.Scan(&data)
	ok(t, err)

	schemaPattern := `CREATE TABLE computer_network_adapters`

	matched, err := regexp.MatchString(schemaPattern, data.String)
	ok(t, err)
	assert(t, matched, "invalid table schema", nil)
}

func TestNetworkAdapterRepositoryCreate(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewNetworkAdapterRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	na := &NetworkAdapter{
		ComputerID: null.IntFrom(1),
		Name:       null.NewString("Local Network", true),
		MacAddress: null.NewString("00:00:00:00:00", true),
		IPAddress:  null.NewString("192.168.1.1", true),
	}

	id, err := repo.Create(dbCtx, na)
	ok(t, err)

	equals(t, int64(1), id)
}

func TestNetworkAdapterRepositorySelect(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewNetworkAdapterRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		na := &NetworkAdapter{
			ComputerID: null.IntFrom(1),
			Name:       null.NewString(fmt.Sprintf("Network %d", i), true),
			MacAddress: null.NewString("00:00:00:00:00", true),
			IPAddress:  null.NewString("192.168.1.1", true),
		}

		_, err := repo.Create(dbCtx, na)
		ok(t, err)
	}

	comp, err := repo.Select(dbCtx, 4)
	ok(t, err)
	equals(t, "Network 3", comp.Name.String)
}

func TestNetworkAdapterRepositoryUpdate(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewNetworkAdapterRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		na := &NetworkAdapter{
			ComputerID: null.IntFrom(1),
			Name:       null.NewString(fmt.Sprintf("Network %d", i), true),
			MacAddress: null.NewString("00:00:00:00:00", true),
			IPAddress:  null.NewString("192.168.1.1", true),
		}

		_, err := repo.Create(dbCtx, na)
		ok(t, err)
	}

	na := &NetworkAdapter{
		ID:   null.IntFrom(4),
		Name: null.NewString("Network 33", true),
	}

	err = repo.Update(dbCtx, na)
	ok(t, err)

	comp, err := repo.Select(dbCtx, 4)
	ok(t, err)
	equals(t, "Network 33", comp.Name.String)
}

func TestNetworkAdapterRepositoryDelete(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewNetworkAdapterRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		na := &NetworkAdapter{
			ComputerID: null.IntFrom(1),
			Name:       null.NewString(fmt.Sprintf("Network %d", i), true),
			MacAddress: null.NewString("00:00:00:00:00", true),
			IPAddress:  null.NewString("192.168.1.1", true),
		}

		_, err := repo.Create(dbCtx, na)
		ok(t, err)
	}

	err = repo.Delete(dbCtx, 4)
	ok(t, err)

	comp, err := repo.Select(dbCtx, 4)
	ok(t, err)
	equals(t, true, comp.Deleted.Valid)
}

func TestNetworkAdapterRepositoryList(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewNetworkAdapterRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		na := &NetworkAdapter{
			ComputerID: null.IntFrom(1),
			Name:       null.NewString(fmt.Sprintf("Network %d", i), true),
			MacAddress: null.NewString("00:00:00:00:00", true),
			IPAddress:  null.NewString("192.168.1.1", true),
		}

		_, err := repo.Create(dbCtx, na)
		ok(t, err)
	}

	nas, err := repo.List(dbCtx, 1, 2)
	ok(t, err)

	equals(t, 2, len(nas))
	equals(t, int64(2), nas[0].ID.Int64)
	equals(t, int64(3), nas[1].ID.Int64)
}

func TestNetworkAdapterRepositoryListNoRows(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewNetworkAdapterRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	nas, err := repo.List(dbCtx, 1, 2)
	ok(t, err)

	equals(t, 0, len(nas))
}

func TestUserRepositoryInstall(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	var data null.String
	row := db.QueryRowContext(dbCtx, "SELECT sql FROM sqlite_master WHERE name='computer_users'")
	err = row.Scan(&data)
	ok(t, err)

	schemaPattern := `CREATE TABLE computer_users`

	matched, err := regexp.MatchString(schemaPattern, data.String)
	ok(t, err)
	assert(t, matched, "invalid table schema", nil)
}

func TestUserRepositoryCreate(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	user := &User{
		ComputerID: null.IntFrom(1),
		Username:   null.StringFrom("Test"),
	}

	id, err := repo.Create(dbCtx, user)
	ok(t, err)

	equals(t, int64(1), id)
}

func TestUserRepositorySelect(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		user := &User{
			ComputerID: null.IntFrom(int64(i)),
			Username:   null.StringFrom(fmt.Sprintf("Test User %d", i)),
		}

		_, err := repo.Create(dbCtx, user)
		ok(t, err)
	}

	comp, err := repo.Select(dbCtx, 4)
	ok(t, err)
	equals(t, "Test User 3", comp.Username.String)
}

func TestUserRepositorySelectWithUsername(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		user := &User{
			ComputerID: null.IntFrom(int64(i)),
			Username:   null.StringFrom(fmt.Sprintf("Test User %d", i)),
		}

		_, err := repo.Create(dbCtx, user)
		ok(t, err)
	}

	comp, err := repo.SelectWithUsername(dbCtx, "Test User 3")
	ok(t, err)
	equals(t, int64(4), comp.ID.Int64)
}

func TestUserRepositorySelectWithUsernameAndComputerID(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		user := &User{
			ComputerID: null.IntFrom(int64(i)),
			Username:   null.StringFrom(fmt.Sprintf("Test User %d", i)),
		}

		_, err := repo.Create(dbCtx, user)
		ok(t, err)
	}

	comp, err := repo.SelectWithUsernameAndComputerID(dbCtx, 4, "Test User 4")
	ok(t, err)
	equals(t, int64(5), comp.ID.Int64)
}

func TestUserRepositoryUpdate(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		user := &User{
			ComputerID: null.IntFrom(int64(i)),
			Username:   null.StringFrom(fmt.Sprintf("Test User %d", i)),
		}

		_, err := repo.Create(dbCtx, user)
		ok(t, err)
	}

	user := &User{
		ID:       null.IntFrom(3),
		Username: null.NewString("Test User 33", true),
	}

	err = repo.Update(dbCtx, user)
	ok(t, err)

	comp, err := repo.SelectWithUsername(dbCtx, "Test User 33")
	ok(t, err)
	equals(t, "Test User 33", comp.Username.String)
}

func TestUserRepositoryDelete(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		user := &User{
			ComputerID: null.IntFrom(int64(i)),
			Username:   null.StringFrom(fmt.Sprintf("Test User %d", i)),
		}

		_, err := repo.Create(dbCtx, user)
		ok(t, err)
	}

	err = repo.Delete(dbCtx, 4)
	ok(t, err)

	comp, err := repo.SelectWithUsername(dbCtx, "Test User 3")
	ok(t, err)
	equals(t, true, comp.Deleted.Valid)
}

func TestUserRepositoryList(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	err = repo.Install(dbCtx)
	ok(t, err)

	for i := 0; i < 10; i++ {
		user := &User{
			ComputerID: null.IntFrom(int64(i)),
			Username:   null.StringFrom(fmt.Sprintf("Test User %d", i)),
		}

		_, err := repo.Create(dbCtx, user)
		ok(t, err)
	}

	users, err := repo.List(dbCtx, 1, 2)
	ok(t, err)

	equals(t, 2, len(users))
	equals(t, int64(2), users[0].ID.Int64)
	equals(t, int64(3), users[1].ID.Int64)
}
