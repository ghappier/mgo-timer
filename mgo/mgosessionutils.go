package mgo

import (
	"fmt"
	"sync"

	"gopkg.in/mgo.v2"
)

var m *MgoSessionUtils
var once sync.Once

func GetInstance() *MgoSessionUtils {
	once.Do(func() {
		session, err := mgo.Dial(URL)
		if err != nil {
			//panic(err) //直接终止程序运行
			fmt.Errorf("无法连接mongodb : ", err.Error())
		}
		m = &MgoSessionUtils{
			mgoSession: session,
		}
	})
	return m
}

type MgoSessionUtils struct {
	mgoSession *mgo.Session
}

func (p MgoSessionUtils) GetSession() *mgo.Session {
	//最大连接池默认为4096
	return p.mgoSession.Clone()
}
