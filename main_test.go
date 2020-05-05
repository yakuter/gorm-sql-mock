package main

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"

	"github.com/jinzhu/gorm"
)

func Setup() (*gorm.DB, sqlmock.Sqlmock) {

	db, mock, _ := sqlmock.New()
	// t.Log(err)

	DB, _ = gorm.Open("postgres", db)
	// t.Log(err)

	DB.LogMode(true)

	return DB, mock
}

func TestGetAll(t *testing.T) {
	DB, mock := Setup()

	const sqlSelectAll = `SELECT * FROM "users"`
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectAll)).
		WillReturnRows(sqlmock.NewRows(nil))
	res := getUsers(DB)

	assert.Nil(t, err)

	expected := []User{}
	assert.Nil(t, deep.Equal(expected, res))
}

func TestGetByID(t *testing.T) {
	DB, mock := Setup()

	// users := []User{}
	user := &User{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
		FirstName: "Erhan",
	}
	// users = append(users, *user)

	rows := sqlmock.
		NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "first_name"}).
		AddRow(user.ID, user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.FirstName)

	const sqlSelectOne = `SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND ((id = $1))`

	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
		WithArgs(user.ID).
		WillReturnRows(rows)

	res, err := getUser(DB, user.ID)

	assert.Nil(t, err)

	assert.Nil(t, deep.Equal(user, res))
}

func TestCreate(t *testing.T) {
	DB, mock := Setup()

	user := &User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
		FirstName: "Erhan",
	}

	const sqlInsert = `
	INSERT INTO "users" ("created_at","updated_at","deleted_at","first_name")
		VALUES ($1,$2,$3,$4) RETURNING "users"."id"`

	mock.ExpectBegin() // start transaction
	mock.ExpectQuery(regexp.QuoteMeta(sqlInsert)).
		WithArgs(user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.FirstName).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))
	mock.ExpectCommit() // commit transaction

	res, err := saveUser(DB, *user)

	assert.Nil(t, err)
	assert.Nil(t, deep.Equal(user.ID, res.ID))
}

func TestUpdate(t *testing.T) {
	DB, mock := Setup()

	user := &User{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
		FirstName: "Erhan",
	}

	const sqlUpdate = `UPDATE "users" SET "created_at" = $1, "updated_at" = $2, "deleted_at" = $3, "first_name" = $4 WHERE "users"."deleted_at" IS NULL AND "users"."id" = $5`
	const sqlSelectOne = `SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL AND "users"."id" = $1 ORDER BY "users"."id" ASC LIMIT 1`

	mock.ExpectBegin() // start transaction
	mock.ExpectExec(regexp.QuoteMeta(sqlUpdate)).
		WithArgs(AnyTime{}, AnyTime{}, user.DeletedAt, user.FirstName, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit() // commit transaction

	// select after update
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectOne)).
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))

	res, err := saveUser(DB, *user)

	assert.Nil(t, err)
	assert.Nil(t, deep.Equal(user.ID, res.ID))
}

func TestAll(t *testing.T) {

	t.Run("GetAll", TestGetAll)
	t.Run("GetByID", TestGetByID)
	t.Run("Create", TestCreate)
	t.Run("Update", TestUpdate)

}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
