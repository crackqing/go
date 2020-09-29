package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

//INIT 引入包处理
type GlobalObj struct {

	//SERVER
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string

	Version        string
	MaxConn        int
	MaxPackageSize uint32

	WorkerPoolSize    uint32 //当前业务的池数量 goroutine
	MaxWorkerPoolSize uint32 //最多开启worker限制条件
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &GlobalObject); err != nil {
		panic(err)
	}

}

func init() {
	GlobalObject = &GlobalObj{
		Host:           "0.0.0.0",
		TcpPort:        8889,
		Name:           "ServerApp",
		Version:        "v0.0.1",
		MaxConn:        10000,
		MaxPackageSize: 100000,
		WorkerPoolSize: 10,
		MaxWorkerPoolSize : 1024,
	}
	GlobalObject.Reload()

}
