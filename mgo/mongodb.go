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
	mgoSession *mgo.Session
)

type UserDao struct {
	BulkSize      int
	FlushInterval uint32
	dataChannel   chan User
}

func NewUserDao(size int, interval uint32) *UserDao {
	u := new(UserDao)
	u.BulkSize = size
	u.FlushInterval = interval
	u.dataChannel = make(chan User, size*5)
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
	u.dataChannel <- user
	if len(u.dataChannel) >= u.BulkSize {
		fmt.Println("批量保存")
		users := u.getData()
		return u.insert(users)
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
	return u.insert(users)
}

func (u *UserDao) getData() *[]User {
	size := len(u.dataChannel)
	if size > u.BulkSize {
		size = u.BulkSize
	}
	users := make([]User, 0, size)
	for i := 0; i < size; i++ {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("channel timeout 1 second")
			break
		case user := <-u.dataChannel:
			users = append(users, user)
		}
	}
	fmt.Println("users = ", users)
	return &users
}

func (u *UserDao) insert(users *[]User) error {
	size := len(*users)
	if size == 0 {
		fmt.Println("没有需要保存的数据")
		return nil
	}
	docs := make([]interface{}, 0, size)
	for _, v := range *users {
		docs = append(docs, v)
	}
	session := u.getSession()
	defer session.Close()
	c := session.DB(DB).C(COLLECTION)
	err := c.Insert(docs...)
	if err != nil {
		fmt.Println("error : ", err)
	}
	return err
}
