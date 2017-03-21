package mgo

import (
	"fmt"
	"testing"
	"time"
)

var bulkSize int = 1000
var flushInterval uint32 = 10

/*
func TestAdd(t *testing.T) {
	fmt.Println("开始执行")
	now := time.Now()
	metricDao := NewMetricDao(bulkSize, flushInterval)
	for i := 0; i < 200; i++ {
		go func(metricDao *MetricDao, now time.Time, index int, t *testing.T) {
			year, month, day := now.Date()
			today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

			data := Data{}
			data.Date = now
			data.Value = float64(index)
			datas := make([]Data, 0, 1)
			datas = append(datas, data)

			ts := Timeseries{}
			ts.Hour = now.Hour()
			ts.Data = datas
			tss := make([]Timeseries, 0, 1)
			tss = append(tss, ts)

			tag := make(map[string]string)
			tag["instance"] = "localhost"
			tag["aa"] = "lizq"
			tag["bb"] = "hezj"

			metric := Metric{}
			metric.Date = today
			metric.MetricName = "gc"
			metric.HostName = "localhost"
			metric.EnvLocation = "www.baidu.com"
			metric.MetricTag = tag
			metric.LastestValue = float64(index)
			metric.LastestUpdateDate = time.Now()
			metric.Timeseries = tss

			for i := 0; i < 1000; i++ {
				err := metricDao.Add(metric)
				if err != nil {
					t.Error("error : ", err)
				}
			}

		}(metricDao, now, i, t)
	}
	time.Sleep(30 * time.Second)
	fmt.Println("完成")
}
*/

func TestMerge(t *testing.T) {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	data1 := Data{}
	data1.Date = time.Now()
	data1.Value = float64(now.Hour() + 1)
	datas1 := make([]Data, 0, 1)
	datas1 = append(datas1, data1)

	ts1 := Timeseries{}
	ts1.Hour = now.Hour()
	ts1.Data = &datas1
	tss1 := make([]Timeseries, 0, 1)
	tss1 = append(tss1, ts1)

	tag1 := make(map[string]string)
	tag1["instance"] = "localhost"
	tag1["aa"] = "lizq"
	tag1["bb"] = "hezj"

	metric1 := Metric{}
	metric1.Date = today
	metric1.MetricName = "gc"
	metric1.HostName = "localhost"
	metric1.EnvLocation = "www.baidu.com"
	metric1.MetricTag = tag1
	metric1.LastestValue = float64(now.Hour() + 1)
	metric1.LastestUpdateDate = time.Now()
	metric1.Timeseries = tss1

	time.Sleep(1 * time.Second)

	data2 := Data{}
	data2.Date = time.Now()
	data2.Value = float64(now.Hour() + 2)
	datas2 := make([]Data, 0, 1)
	datas2 = append(datas2, data2)

	ts2 := Timeseries{}
	ts2.Hour = now.Hour()
	ts2.Data = &datas2
	tss2 := make([]Timeseries, 0, 1)
	tss2 = append(tss2, ts2)

	tag2 := make(map[string]string)
	tag2["instance"] = "localhost"
	tag2["aa"] = "lizq"
	tag2["bb"] = "hezj"

	metric2 := Metric{}
	metric2.Date = today
	metric2.MetricName = "gc"
	metric2.HostName = "localhost"
	metric2.EnvLocation = "www.baidu.com"
	metric2.MetricTag = tag2
	metric2.LastestValue = float64(now.Hour() + 2)
	metric2.LastestUpdateDate = time.Now()
	metric2.Timeseries = tss2

	time.Sleep(1 * time.Second)

	data3 := Data{}
	data3.Date = time.Now()
	data3.Value = float64(now.Hour() + 3)
	datas3 := make([]Data, 0, 1)
	datas3 = append(datas3, data3)

	ts3 := Timeseries{}
	ts3.Hour = now.Hour()
	ts3.Data = &datas3
	tss3 := make([]Timeseries, 0, 1)
	tss3 = append(tss3, ts3)

	tag3 := make(map[string]string)
	tag3["instance"] = "localhost"
	tag3["aa"] = "lizq"
	tag3["bb"] = "hezj"

	metric3 := Metric{}
	metric3.Date = today
	metric3.MetricName = "gc3"
	metric3.HostName = "localhost"
	metric3.EnvLocation = "www.baidu.com"
	metric3.MetricTag = tag3
	metric3.LastestValue = float64(now.Hour() + 3)
	metric3.LastestUpdateDate = time.Now()
	metric3.Timeseries = tss3

	metricDao := NewMetricDao(bulkSize, flushInterval)

	//metrics := make([]Metric, 0, 3)
	metrics := new([]Metric)

	metrics = metricDao.merge(metrics, metric1)
	if len(*metrics) != 1 {
		t.Error("length of metrics should be 1")
	}
	if (*metrics)[0].LastestValue != metric1.LastestValue || (*metrics)[0].LastestUpdateDate != metric1.LastestUpdateDate {
		t.Error("LastestValue and LastestUpdateDate can not be updated")
	}

	metrics = metricDao.merge(metrics, metric2)
	for _, v := range *metrics {
		for _, ts := range v.Timeseries {
			if ts.Hour == now.Hour() {
				if len(*ts.Data) != 2 {
					t.Error("when ts.Hour = ", now.Hour(), ", length of Data should be 2")
				}
			}
		}
	}

	metrics = metricDao.merge(metrics, metric3)
	if len(*metrics) != 2 {
		t.Error("length of metrics should be 2")
	}

	fmt.Println("metric1 = ", metric1.String())
	fmt.Println("metric2 = ", metric2.String())
	fmt.Println("metric3 = ", metric3.String())
	fmt.Println("metrics = ", MetricSlice(*metrics).String())
}
