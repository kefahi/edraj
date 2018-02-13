package main

import (
	"fmt"
	"reflect"
)

// Content ...
type Content struct {
	ID string
}

// Addon ...
type Addon struct {
	Name string
}

// Entry ...
type Entry struct {
	Cool  *Content
	Addon *Addon
}

func main() {

	/*
		e := Entry{
			cool:  &Content{ID: "one"},
			Addon: &Addon{Name: "Ali"},
		}*/

	e := reflect.New(reflect.TypeOf(Entry{}))

	e.Elem().FieldByName("Cool").Set(reflect.ValueOf(&Content{ID: "three"}))
	e.Elem().FieldByName("Addon").Set(reflect.ValueOf(&Addon{Name: "Mo"}))

	//t := reflect.TypeOf(e)

	/*
		t := reflect.TypeOf(e)
		s := "cool"
		val, ok := t.FieldByName(s)

		if !ok {
			panic("Lets bail out")
		}
	*/
	//v := reflect.ValueOf(val)
	//val. = &Content{ID: "two"}
	//val.(&Content{ID: "two"})

	fmt.Printf("Hello %T %v\n", e.Elem().FieldByName("Cool").Interface(), e.Elem().FieldByName("Addon"))

	e.Elem().FieldByName("Cool").Set(reflect.ValueOf(&Content{ID: "one"}))
	fmt.Printf("Hello %v %v\n", e.Elem().FieldByName("Cool"), e.Elem().FieldByName("Addon"))

}
