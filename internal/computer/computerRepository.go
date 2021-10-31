package computer

import (
	"fpsmonitor/internal/model"

	"github.com/jmoiron/sqlx"
)

type ComputerRepository struct {
	model.Repository
}

func NewComputerRepository(db *sqlx.DB) *ComputerRepository {
	repo := &ComputerRepository{}
	repo.SetDB(db)

	repo.InstallCmd = `CREATE TABLE "computers" (
		"id" INTEGER,
		"name" TEXT,
		"created" TEXT,
		"updated" TEXT,
		"deleted" TEXT,
		PRIMARY KEY("id" AUTOINCREMENT)
	)`

	repo.CreateCmd = `INSERT INTO computers(
		name,
		created
	) VALUES (?,?)`

	repo.SelectCmd = `SELECT 
		id,
		name,
		created,
		updated,
		deleted
	FROM computers 
	WHERE name=?`

	repo.UpdateCmd = `UPDATE computers SET
		name=?,
		updated=?
	WHERE id=?`

	repo.DeleteCmd = `UPDATE computers SET
		deleted=?
	WHERE id=?`

	return repo
}
