package computer

import (
	"fpsmonitor/internal/model"

	"gopkg.in/guregu/null.v3"
)

type NetworkAdapter struct {
	ID         null.Int    `db:"id" json:"-"`
	ComputerID null.Int    `db:"computer_id" json:"-"`
	Name       null.String `db:"name" json:"name"`
	MacAddress null.String `db:"mac_address" json:"mac_address"`
	IPAddress  null.String `db:"ip_address" json:"ip_address"`
	Created    null.String `db:"created" json:"-"`
	Updated    null.String `db:"updated" json:"-"`
	Deleted    null.String `db:"deleted" json:"-"`
}

type Computer struct {
	model.Model

	Name null.String `db:"name" json:"name"`
}

type User struct {
	ID         null.Int    `db:"id" json:"-"`
	ComputerID null.Int    `db:"computer_id" json:"-"`
	Username   null.String `db:"username" json:"username"`
	Created    null.String `db:"created" json:"-"`
	Updated    null.String `db:"updated" json:"-"`
	Deleted    null.String `db:"deleted" json:"-"`
}
