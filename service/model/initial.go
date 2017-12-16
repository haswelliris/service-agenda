package model

import (
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var Engine *xorm.Engine

func init() {
	// initialize xorm engine
	engine, err := xorm.NewEngine("sqlite3", "./agenda.db")
	checkErr(err)
	engine.SetMapper(core.SameMapper{})
	err = engine.Sync2(new(User), new(Meeting))
	checkErr(err)
	Engine = engine
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
