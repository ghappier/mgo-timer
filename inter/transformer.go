package inter

type Vehicle struct {
	Name  string
	Seats int
}
type Transformer interface {
	Transform(vehicle Vehicle) interface{}
}
