package inter

import (
	"fmt"
)

type Person struct {
	transformer Transformer
}

func NewPerson(transformer Transformer) *Person {
	person := new(Person)
	person.transformer = transformer
	return person
}

func (p *Person) Save(vehicle Vehicle) {
	val := p.transformer.Transform(vehicle)
	//以下两种方法均可
	/*
		switch val.(type) {
		case *Car:
			if car, ok := val.(*Car); ok {
				fmt.Printf("Car : {Name:%s, Color:%s}\n", car.Name, car.Color)
			}
		case *Bicycle:
			if bicycle, ok := val.(*Bicycle); ok {
				fmt.Printf("Bicycle : {Name:%s, Speed:%d}\n", bicycle.Name, bicycle.Speed)
			}
		default:
			fmt.Println("Unknow")
		}
	*/
	if car, ok := val.(*Car); ok {
		fmt.Printf("Car : {Name:%s, Color:%s}\n", car.Name, car.Color)
	} else if bicycle, ok := val.(*Bicycle); ok {
		fmt.Printf("Bicycle : {Name:%s, Speed:%d}\n", bicycle.Name, bicycle.Speed)
	} else {
		fmt.Println("Unknow")
	}
}
