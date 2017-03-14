package mgo

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

const (
	URL        = "localhost:27017"
	DB         = "promtest"
	COLLECTION = "prom"
)

var (
	mgoSession  *mgo.Session
	dataChannel chan User
)

type UserDao struct {
	BulkSize      int
	FlushInterval uint32
}

func NewUserDao(size int, interval uint32) *UserDao {
	u := new(UserDao)
	u.BulkSize = size
	u.FlushInterval = interval
	dataChannel = make(chan User, size*5)
	return u
}

func (u *UserDao) getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(URL)
		if err != nil {
			//panic(err) //直接终止程序运行
			fmt.Errorf("无法连接mongodb : ", err.Error())
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

func (u *UserDao) Add(user User) error {
	dataChannel <- user
	if len(dataChannel) >= u.BulkSize {
		fmt.Println("批量保存")
		users := u.getData()
		return u.insert(users...)
	}
	return nil
}
func (u *UserDao) UserDaoTimer() {
	timer1 := time.NewTicker(time.Duration(u.FlushInterval) * time.Second)
	for {
		select {
		case <-timer1.C:
			u.callback(nil)
		}
	}
}

func (u *UserDao) callback(args interface{}) error {
	fmt.Println("超时自动保存")
	users := u.getData()
	return u.insert(users...)
}

func (u *UserDao) getData() []User {
	size := len(dataChannel)
	if size > u.BulkSize {
		size = u.BulkSize
	}
	fmt.Println("size = ", size)
	users := make([]User, size)
	for i := 0; i < size; i++ {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("channel timeout 1 second")
			break
		case user := <-dataChannel:
			users[i] = user
		}
	}
	index := 0
	for i := 0; i < size; i++ {
		if users[i].Name == "" {
			index = i
			break
		}
	}
	var result []User = users[:index]
	return result
}

func (u *UserDao) insert(users ...User) error {
	if len(users) == 0 {
		return nil
	}
	session := u.getSession()
	defer session.Close()
	c := session.DB(DB).C(COLLECTION)
	return c.Insert(users)
}
