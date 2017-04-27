package inter

type Car struct {
	Name  string
	Color string
}
type MyTransformer struct {
}

func (t MyTransformer) Transform(vehicle Vehicle) interface{} {
	car := new(Car)
	car.Name = vehicle.Name
	car.Color = "red"
	return car
}
