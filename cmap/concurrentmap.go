package cmap

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ghappier/mgo-timer/mgo"
	"gopkg.in/mgo.v2/bson"
)

type ConcurrentMap struct {
	mutex          sync.RWMutex
	dataChannelMap map[string]chan mgo.Metric
	bulk           int
	interval       int
	saveCount      int
	insertCount    int
}

func NewConcurrentMap(bulkSize int, flushInterval int) *ConcurrentMap {
	c := ConcurrentMap{}
	c.dataChannelMap = make(map[string]chan mgo.Metric)
	c.bulk = bulkSize
	c.interval = flushInterval
	go c.StorageTimer()
	return &c
}

func (c *ConcurrentMap) Save(collectionName string, metric mgo.Metric) {
	c.saveCount += 1
	//fmt.Println("总共save", c.saveCount)
	var mChannel chan mgo.Metric
	var ok bool
	c.mutex.RLock()
	if mChannel, ok = c.dataChannelMap[collectionName]; !ok {
		c.mutex.RUnlock()
		c.mutex.Lock()
		if mChannel, ok = c.dataChannelMap[collectionName]; !ok {
			fmt.Println("未存在集合", collectionName)
			mChannel = make(chan mgo.Metric, c.bulk)
			c.dataChannelMap[collectionName] = mChannel
		}
		c.mutex.Unlock()
	} else {
		c.mutex.RUnlock()
	}
	c.inSave(mChannel, metric, collectionName)
}
func (c *ConcurrentMap) inSave(myChannel chan mgo.Metric, metric mgo.Metric, collectionName string) {
	select {
	case myChannel <- metric: //未满
		//fmt.Println("指标入channel：", metric)
	default: //已满,阻塞
		fmt.Println("批量保存")
		c.bulkSave(collectionName, myChannel)
		myChannel <- metric
	}
}
func (c *ConcurrentMap) bulkSave(collectionName string, myChannel chan mgo.Metric) {
	var metrics []mgo.Metric = make([]mgo.Metric, 0, c.bulk)
LABEL_BREAK:
	for i := 0; i < c.bulk; i++ {
		select {
		case metric := <-myChannel:
			metrics = append(metrics, metric)
		default: //myChannel为空
			break LABEL_BREAK
		}
	}
	go func(colName string, mets *[]mgo.Metric) {
		size := len(*mets)
		if size == 0 {
			fmt.Println("没有需要保存的数据")
		}
		c.insertCount += size
		//fmt.Println("总共insert ", c.insertCount)
		//fmt.Println("指标入mongodb：", *mets)
		session := mgo.GetInstance().GetSession()
		defer session.Close()
		c := session.DB(mgo.DB).C(colName)
		bulk := c.Bulk()
		for _, v := range *mets {
			selector := bson.M{"date": v.Date, "metric_name": v.MetricName, "hostname": v.HostName}
			for k, v := range v.MetricTag {
				selector["metric_tag."+k] = v
			}
			updateSet := bson.M{"lastest_value": v.LastestValue, "lastest_update_date": v.LastestUpdateDate}
			timeSeries := make([]mgo.Timeseries, 0, 24)
			for i := 0; i < 24; i++ {
				ts := mgo.Timeseries{Hour: i, Data: make([]mgo.Data, 0, 0)}
				timeSeries = append(timeSeries, ts)
			}
			updateSetOnInsert := bson.M{"timeseries": timeSeries}
			bulk.Upsert(selector, bson.M{"$set": updateSet, "$setOnInsert": updateSetOnInsert})
			for _, ts := range v.Timeseries {
				bulk.Update(selector, bson.M{"$set": bson.M{"lastest_value": v.LastestValue}, "$push": bson.M{"timeseries." + strconv.Itoa(ts.Hour) + ".data": bson.M{"$each": ts.Data}}})
			}
		}
		bulkResult, err := bulk.Run()
		if err != nil {
			fmt.Printf("bulk save error : %s\n", err.Error())
		} else {
			fmt.Printf("bulk save success : {matched: %d, modified: %d}\n", bulkResult.Matched, bulkResult.Modified)
		}
	}(collectionName, &metrics)
}

//定期保存
func (c *ConcurrentMap) StorageTimer() {
	ticker := time.NewTicker(time.Duration(c.interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("超时自动保存")
			for k, v := range c.dataChannelMap {
				c.bulkSave(k, v)
			}
		}
	}
}
