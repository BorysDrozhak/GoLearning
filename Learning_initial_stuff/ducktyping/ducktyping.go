package main

import "fmt"

type ducker interface {
	cry()
	// jump() # not allowed
}

type eater interface {
	eat()
}

type wildDuck struct {
}

type homeDuck struct {
}

func (d wildDuck) cry() {
	fmt.Println("cry")
}

func (d wildDuck) fly() {
	fmt.Println("fly")
}

func (d homeDuck) cry() {
	fmt.Println("cry like homeDuck", "cry")
}

func (d homeDuck) eat() {
	fmt.Println("eaaat!")
}

func main() {
	var e eater
	var d ducker

	if false {
		d = wildDuck{}
	} else {
		d = homeDuck{}
	}
	d.cry()

	e = homeDuck{}
	e.eat()

	e.(homeDuck).cry() // casting

	a := wildDuck{}
	a.fly()
}
