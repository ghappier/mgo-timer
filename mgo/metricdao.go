package mgo

import (
	"fmt"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type MetricDao struct {
	BulkSize      int
	FlushInterval uint32
	dataChannel   chan Metric
}

func NewMetricDao(size int, interval uint32) *MetricDao {
	u := new(MetricDao)
	u.BulkSize = size
	u.FlushInterval = interval
	u.dataChannel = make(chan Metric, size*5)
	go u.MetricDaoTimer()
	return u
}

func (u *MetricDao) Add(metric Metric) error {
	u.dataChannel <- metric
	if len(u.dataChannel) >= u.BulkSize {
		fmt.Println("批量保存")
		metrics := u.getData()
		return u.insert(metrics)
	}
	return nil
}
func (u *MetricDao) MetricDaoTimer() {
	timer1 := time.NewTicker(time.Duration(u.FlushInterval) * time.Second)
	for {
		select {
		case <-timer1.C:
			u.callback(nil)
		}
	}
}

func (u *MetricDao) callback(args interface{}) error {
	fmt.Println("超时自动保存")
	metrics := u.getData()
	return u.insert(metrics)
}

func (u *MetricDao) getData() *[]Metric {
	size := len(u.dataChannel)
	if size > u.BulkSize {
		size = u.BulkSize
	}
	metrics := make([]Metric, 0, size)
	for i := 0; i < size; i++ {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("channel timeout 1 second")
			break
		case metric := <-u.dataChannel:
			metrics = append(metrics, metric)
		}
	}
	fmt.Println("metric = ", metrics)
	return &metrics
}

func (u *MetricDao) insert(metrics *[]Metric) error {
	size := len(*metrics)
	if size == 0 {
		fmt.Println("没有需要保存的数据")
		return nil
	}
	//session := u.getSession()
	session := GetInstance().GetSession()
	defer session.Close()
	c := session.DB(DB).C(COLLECTION)
	bulk := c.Bulk()
	for _, v := range *metrics {
		selector := bson.M{"metric_name": v.MetricName, "hostname": v.Hostname}
		updateSet := bson.M{"lastest_value": v.LastestValue}

		timeSeries := make([]Timeseries, 0, 24)
		for i := 0; i < 24; i++ {
			ts := Timeseries{Hour: i, Data: make([]Data, 0, 0)}
			timeSeries = append(timeSeries, ts)
		}
		updateSetOnInsert := bson.M{"timeseries": timeSeries}
		bulk.Upsert(selector, bson.M{"$set": updateSet, "$setOnInsert": updateSetOnInsert})
		for _, ts := range v.Timeseries {
			bulk.Update(selector, bson.M{"$set": bson.M{"lastest_value": v.LastestValue}, "$push": bson.M{"timeseries." + strconv.Itoa(ts.Hour) + ".data": bson.M{"$each": ts.Data}}})
		}
	}
	_, err := bulk.Run()
	return err
}

/*
db.prom.update(
    {metric_name: "gc", hostname: "localhost"},
    {
        $set: {lastest_value: 300},
        $setOnInsert: {
            timeseries: [
                {hour: 0,data: []},
                {hour: 1,data: []},
                {hour: 2,data: []},
                {hour: 3,data: []},
                {hour: 4,data: []},
                {hour: 5,data: []},
                {hour: 6,data: []},
                {hour: 7,data: []},
                {hour: 8,data: []},
                {hour: 9,data: []},
                {hour: 10,data: []},
                {hour: 11,data: []},
                {hour: 12,data: []},
                {hour: 13,data: []},
                {hour: 14,data: []},
                {hour: 15,data: []},
                {hour: 16,data: []},
                {hour: 17,data: []},
                {hour: 18,data: []},
                {hour: 19,data: []},
                {hour: 20,data: []},
                {hour: 21,data: []},
                {hour: 22,data: []},
                {hour: 23,data: []}
            ]
        }
    },
    {upsert:true}
)


db.prom.update(
    {metric_name: "gc", hostname: "localhost"},
    {
        $push: {"timeseries.0.data":{name:"lizq",value:100}},
        $set: {lastest_value: 200}
    }
)

*/
