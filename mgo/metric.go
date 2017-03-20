package mgo

type Data struct {
	Name  string  `bson:"name"`
	Value float64 `bson:"value"`
}

type Timeseries struct {
	Hour int    `bson:"hour"`
	Data []Data `bson:"data"`
}
type Metric struct {
	MetricName   string       `bson:"metric_name"`
	Hostname     string       `bson:"hostname"`
	LastestValue float64      `bson:"lastest_value"`
	Timeseries   []Timeseries `bson:"timeseries"`
}
