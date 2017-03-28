package cmap

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/ghappier/mgo-timer/mgo"
)

var bulkSize int = 5000
var flushInterval int = 5
var collectionNames = [...]string{
	"lizq000",
	"lizq001",
	"lizq002",
	"lizq003",
	"lizq004",
	"lizq005",
	"lizq006",
	"lizq007",
	"lizq008",
	"lizq009",
	"lizq010",
	"lizq011",
	"lizq012",
	"lizq013",
	"lizq014",
	"lizq015",
	"lizq016",
	"lizq017",
	"lizq018",
	"lizq019",
	"lizq020",
	"lizq021",
	"lizq022",
	"lizq023",
	"lizq024",
	"lizq025",
	"lizq026",
	"lizq027",
	"lizq028",
	"lizq029",
	"lizq030",
	"lizq031",
	"lizq032",
	"lizq033",
	"lizq034",
	"lizq035",
	"lizq036",
	"lizq037",
	"lizq038",
	"lizq039",
	"lizq040",
	"lizq041",
	"lizq042",
	"lizq043",
	"lizq044",
	"lizq045",
	"lizq046",
	"lizq047",
	"lizq048",
	"lizq049",
	"lizq050",
	"lizq051",
	"lizq052",
	"lizq053",
	"lizq054",
	"lizq055",
	"lizq056",
	"lizq057",
	"lizq058",
	"lizq059",
	"lizq060",
	"lizq061",
	"lizq062",
	"lizq063",
	"lizq064",
	"lizq065",
	"lizq066",
	"lizq067",
	"lizq068",
	"lizq069",
	"lizq070",
	"lizq071",
	"lizq072",
	"lizq073",
	"lizq074",
	"lizq075",
	"lizq076",
	"lizq077",
	"lizq078",
	"lizq079",
	"lizq080",
	"lizq081",
	"lizq082",
	"lizq083",
	"lizq084",
	"lizq085",
	"lizq086",
	"lizq087",
	"lizq088",
	"lizq089",
	"lizq090",
	"lizq091",
	"lizq092",
	"lizq093",
	"lizq094",
	"lizq095",
	"lizq096",
	"lizq097",
	"lizq098",
	"lizq099",
	"lizq100",
	"lizq101",
	"lizq102",
	"lizq103",
	"lizq104",
	"lizq105",
	"lizq106",
	"lizq107",
	"lizq108",
	"lizq109",
	"lizq110",
	"lizq111",
	"lizq112",
	"lizq113",
	"lizq114",
	"lizq115",
	"lizq116",
	"lizq117",
	"lizq118",
	"lizq119",
	"lizq120",
	"lizq121",
	"lizq122",
	"lizq123",
	"lizq124",
	"lizq125",
	"lizq126",
	"lizq127",
	"lizq128",
	"lizq129",
	"lizq130",
	"lizq131",
	"lizq132",
	"lizq133",
	"lizq134",
	"lizq135",
	"lizq136",
	"lizq137",
	"lizq138",
	"lizq139",
	"lizq140",
	"lizq141",
	"lizq142",
	"lizq143",
	"lizq144",
	"lizq145",
	"lizq146",
	"lizq147",
	"lizq148",
	"lizq149",
	"lizq150",
	"lizq151",
	"lizq152",
	"lizq153",
	"lizq154",
	"lizq155",
	"lizq156",
	"lizq157",
	"lizq158",
	"lizq159",
	"lizq160",
	"lizq161",
	"lizq162",
	"lizq163",
	"lizq164",
	"lizq165",
	"lizq166",
	"lizq167",
	"lizq168",
	"lizq169",
	"lizq170",
	"lizq171",
	"lizq172",
	"lizq173",
	"lizq174",
	"lizq175",
	"lizq176",
	"lizq177",
	"lizq178",
	"lizq179",
	"lizq180",
	"lizq181",
	"lizq182",
	"lizq183",
	"lizq184",
	"lizq185",
	"lizq186",
	"lizq187",
	"lizq188",
	"lizq189",
	"lizq190",
	"lizq191",
	"lizq192",
	"lizq193",
	"lizq194",
	"lizq195",
	"lizq196",
	"lizq197",
	"lizq198",
	"lizq199",
}

func TestMerge(t *testing.T) {
	fmt.Println("开始执行")
	session := mgo.GetInstance().GetSession()
	defer session.Close()
	for _, v := range collectionNames {
		c := session.DB(mgo.DB).C(v)
		c.RemoveAll(nil)
	}
	now := time.Now()
	metricDao := NewConcurrentMap(bulkSize, flushInterval)
	times := 100
	num := 400
	var wg sync.WaitGroup
	for i := 0; i < times; i++ {
		wg.Add(1)
		go func(metricDao *ConcurrentMap, now time.Time, index int, t *testing.T) {
			for j := 0; j < num; j++ {
				year, month, day := now.Date()
				today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

				data := mgo.Data{}
				data.Date = now
				data.Value = float64(index)
				datas := make([]mgo.Data, 0, 1)
				datas = append(datas, data)

				ts := mgo.Timeseries{}
				ts.Hour = now.Hour()
				ts.Data = datas
				tss := make([]mgo.Timeseries, 0, 1)
				tss = append(tss, ts)

				tag := make(map[string]string)
				tag["instance"] = "localhost"
				tag["aa"] = "lizq"
				tag["bb"] = "hezj"

				metric := mgo.Metric{}
				metric.Date = today
				metric.MetricName = collectionNames[rand.Intn(len(collectionNames))]
				metric.HostName = "localhost"
				metric.EnvLocation = "www.baidu.com"
				metric.MetricTag = tag
				metric.LastestValue = float64(index)
				metric.LastestUpdateDate = time.Now()
				metric.Timeseries = tss
				collectionName := metric.MetricName
				metricDao.Save(collectionName, metric)
			}
			wg.Done()
		}(metricDao, now, i, t)
	}
	wg.Wait()
	time.Sleep(time.Duration(flushInterval*2) * time.Second)
	end := time.Now()
	fmt.Println("time costs :", end.Sub(now))
	mcount := 0
	for _, v := range collectionNames {
		c := session.DB(mgo.DB).C(v)
		var mm mgo.Metric
		err := c.Find(nil).One(&mm)
		if err != nil {
			t.Error("error : ", err.Error())
		} else {
			mcount += len(mm.Timeseries[now.Hour()].Data)
		}
	}
	if mcount != times*num {
		t.Error("mcount ", mcount, ", expected ", times*num, ", error")
	}
	fmt.Println("完成")
}
