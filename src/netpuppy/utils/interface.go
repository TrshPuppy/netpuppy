package utils

import "fmt"

type dummyInterface interface {
	contact() string
	//pickupLine() string
}

type DummyStruct struct {
	name   string
	number int
}

func (d DummyStruct) contact() string {
	contactString := fmt.Sprintf("%v: %v", d.name, d.number)
	return contactString
}

func buildContact(d dummyInterface) string {
	contactString := d.contact()
	return contactString
}

func Tiddies() DummyStruct {
	Bradley := DummyStruct{name: "Bradley", number: 69}

	getContact := buildContact(Bradley)
	fmt.Printf("Bradley's contact: %v\n", getContact)
	return Bradley
}
