package auth

import "gopkg.in/guregu/null.v3"

type UserModel struct {
	ID      null.Int    `db:"id" json:"id"`
	Created null.String `db:"created" json:"created"`
	Updated null.String `db:"updated" json:"updated"`
	Deleted null.String `db:"deleted" json:"deleted"`

	Username null.String `db:"username" json:"username"`
	Password null.String `db:"password" json:"password"`
}
