package main

import (
	"fmt"
	"time"
)

type DB interface {
	Find() interface{}
	Update(interface{})
}

type V interface {
	Render(interface{}) string
}

// type COMMON interface {
// 	V
// 	DB
// }

type Controller struct {
	dataProvider DB
	viewRender   V
}

func NewController() *Controller {
	var controler Controller
	// controler = new(Controller)
	// get config
	var db DB
	var v V
	if false {
		db = newRDB("mysql")
	} else {
		db = newODB("S3")
	}
	controler.dataProvider = db

	if true {
		v = new(JSONv)
	} else {
		v = new(HTMLv)
	}

	controler.viewRender = v
	return &controler
}

func (c *Controller) work1() string {
	v := c.dataProvider.Find()
	return c.viewRender.Render(v)
}

type JSONv struct{}
type HTMLv struct{}

func (j *HTMLv) Render(i interface{}) string {
	return fmt.Sprintln("kind if HTML render for", i)
}

func (j *JSONv) Render(i interface{}) string {
	return fmt.Sprintln("damn render of", i)
}

type RDB struct {
	connection string
}

func (i *RDB) Find() interface{} {
	fmt.Println("finding in RDB:", i.connection)
	return "user_in_db"
}

func newRDB(name string) *RDB {
	var db RDB
	db.connection = name
	return &db
}

func (i *RDB) Update(inter interface{}) {
	fmt.Println("update in", i.connection, "doing:", inter)
}

//

type ODB struct {
	connection string
}

func (i *ODB) Find() interface{} {
	fmt.Println("finding in ODB:", i.connection)
	return "REsult_of_find"
}

func newODB(name string) *ODB {
	var db ODB
	db.connection = name
	return &db
}

func (i *ODB) Update(inter interface{}) {
	fmt.Println("save in", i.connection, "doing:", inter)
}

// type COMMON1 struct {
// 	RDB
// 	JSONv
// }

func main() {
	c := NewController()
	fmt.Println(c.work1())

	c.viewRender.Render()

	// v := NewController()

	// var test COMMON

}
