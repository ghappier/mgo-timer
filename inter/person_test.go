package inter

import (
	"testing"
)

func TestSave(t *testing.T) {
	vehicle := Vehicle{Name: "bm", Seats: 6}
	myTransformer := MyTransformer{}
	person1 := NewPerson(myTransformer)
	person1.Save(vehicle)
	herTransformer := HerTransformer{}
	person2 := NewPerson(herTransformer)
	person2.Save(vehicle)
}
