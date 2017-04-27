package inter

type Bicycle struct {
	Name  string
	Speed int
}
type HerTransformer struct {
}

func (t HerTransformer) Transform(vehicle Vehicle) interface{} {
	bicycle := new(Bicycle)
	bicycle.Name = vehicle.Name
	bicycle.Speed = 100
	return bicycle
}
