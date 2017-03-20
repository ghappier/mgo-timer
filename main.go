package main

import (
	"runtime"
	"time"

	"github.com/ghappier/mgo-timer/mgo"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	user := mgo.User{
		Name:      "lizq",
		Age:       30,
		JonedAt:   time.Now(),
		Interests: []string{"aaa", "bbb", "ccc"},
	}
	var bulkSize int = 10
	var flushInterval uint32 = 5
	userDao := mgo.NewUserDao(bulkSize, flushInterval)
	userDao.Add(user)
	go func() {
		time.Sleep(time.Duration(userDao.FlushInterval+10) * time.Second)
		for i := 0; i < userDao.BulkSize+3; i++ {
			userDao.Add(user)
		}
		time.Sleep(time.Duration(userDao.FlushInterval+10) * time.Second)
		for i := 0; i < userDao.BulkSize+3; i++ {
			userDao.Add(user)
		}
	}()
	//userDao.UserDaoTimer()
	time.Sleep(time.Duration(userDao.FlushInterval*10) * time.Second)
}
