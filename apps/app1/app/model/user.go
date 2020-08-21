package model

type SchemaTableUser struct {
	ID   int    `xorm:"id <-" json:"id"`
	Name string `xorm:"name" json:"name"`
	Pass string `xorm:"pass" json:"pass"`
}
