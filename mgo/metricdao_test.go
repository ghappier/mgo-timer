package mgo

import (
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	for i := 0; i < 10000; i++ {
		go save(i, t)
	}
	time.Sleep(30 * time.Second)
}

func save(index int, t *testing.T) {
	var bulkSize int = 10
	var flushInterval uint32 = 5
	metricDao := NewMetricDao(bulkSize, flushInterval)

	data := Data{}
	data.Name = "lizq"
	data.Value = float64(index)
	datas := make([]Data, 0, 1)
	datas = append(datas, data)

	ts := Timeseries{}
	ts.Hour = 0
	ts.Data = datas
	tss := make([]Timeseries, 0, 1)
	tss = append(tss, ts)

	metric := Metric{}
	metric.MetricName = "gc"
	metric.Hostname = "localhost"
	metric.LastestValue = float64(index)
	metric.Timeseries = tss

	err := metricDao.Add(metric)
	if err != nil {
		t.Error("error")
	}
}
