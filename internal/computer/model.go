package computer

import "gopkg.in/guregu/null.v3"

type ComputerModel struct {
	ID      null.Int    `db:"id" json:"id"`
	Created null.String `db:"created" json:"created"`
	Updated null.String `db:"updated" json:"updated"`
	Deleted null.String `db:"deleted" json:"deleted"`

	Name           null.String           `db:"name" json:"name"`
	NetworkAdapter []NetworkAdapterModel `db:"-" json:"network_adapters"`
	Users          []UserModel           `db:"-" json:"users"`
}

type UserModel struct {
	ID      null.Int    `db:"id" json:"id"`
	Created null.String `db:"created" json:"created"`
	Updated null.String `db:"updated" json:"updated"`
	Deleted null.String `db:"deleted" json:"deleted"`

	ComputerID null.Int    `db:"computer_id" json:"-"`
	Username   null.String `db:"username" json:"username"`
}

type NetworkAdapterModel struct {
	ID      null.Int    `db:"id" json:"id"`
	Created null.String `db:"created" json:"created"`
	Updated null.String `db:"updated" json:"updated"`
	Deleted null.String `db:"deleted" json:"deleted"`

	ComputerID null.Int    `db:"computer_id" json:"-"`
	Name       null.String `db:"name" json:"name"`
	MacAddress null.String `db:"mac_address" json:"mac_address"`
	IPAddress  null.String `db:"ip_address" json:"ip_address"`
}
