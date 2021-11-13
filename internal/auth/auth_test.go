package auth

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
	"golang.org/x/crypto/bcrypt"
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

func TestMigration(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	err = Migrate(db)
	ok(t, err)

	var data null.String
	row := db.QueryRowContext(dbCtx, "SELECT sql FROM sqlite_master WHERE name='users'")
	err = row.Scan(&data)
	ok(t, err)

	schemaPattern := `CREATE TABLE users`

	matched, err := regexp.MatchString(schemaPattern, data.String)
	ok(t, err)
	assert(t, matched, "invalid table schema", nil)
}

func TestUserStore_Create(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	err = Migrate(db)
	ok(t, err)

	store := NewUserStore(db)

	user := &UserModel{
		Username: null.StringFrom("testUser"),
		Password: null.StringFrom("testUser"),
	}

	id, err := store.Create(dbCtx, user)
	ok(t, err)

	equals(t, int64(1), id.Int64)
}

func TestUserStore_GetByID(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	err = Migrate(db)
	ok(t, err)

	store := NewUserStore(db)

	user := &UserModel{
		Username: null.StringFrom("testUser"),
		Password: null.StringFrom("testUser"),
	}

	id, err := store.Create(dbCtx, user)
	ok(t, err)

	userData, err := store.Get(dbCtx, id)
	ok(t, err)
	equals(t, "testUser", userData.Username.String)
}

func TestUserStore_UpdatePassword(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	// create table
	err = Migrate(db)
	ok(t, err)

	store := NewUserStore(db)

	user := &UserModel{
		Username: null.StringFrom("testUser"),
		Password: null.StringFrom("testUser"),
	}

	// create new user
	id, err := store.Create(dbCtx, user)
	ok(t, err)

	// retrieve new user from database
	newUser, err := store.Get(dbCtx, id)
	ok(t, err)

	// set users new password
	newUser.Password = null.StringFrom("updatededUserPassword")

	// update database record
	err = store.UpdatePassword(dbCtx, newUser)
	ok(t, err)

	// retrieve user again
	newUserWithUpdatedPassword, err := store.Get(dbCtx, id)
	ok(t, err)

	// verify new password matches
	err = bcrypt.CompareHashAndPassword([]byte(newUserWithUpdatedPassword.Password.String), []byte("updatededUserPassword"))
	ok(t, err)
}

func TestUserStore_SoftDelete(t *testing.T) {
	db, err := dbSetup()
	ok(t, err)
	defer db.Close()

	err = Migrate(db)
	ok(t, err)

	store := NewUserStore(db)

	user := &UserModel{
		Username: null.StringFrom("testUser"),
		Password: null.StringFrom("testUser"),
	}

	id, err := store.Create(dbCtx, user)
	ok(t, err)

	userData, err := store.Get(dbCtx, id)
	ok(t, err)

	err = store.SoftDelete(dbCtx, userData)
	ok(t, err)

	deletedUserData, err := store.Get(dbCtx, id)
	ok(t, err)

	equals(t, true, deletedUserData.Deleted.Valid)
}
