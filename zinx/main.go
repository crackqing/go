package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PintRouter struct {
	znet.BaseRouter
}

func (p *PintRouter) Handle(request ziface.IRequest) {
	fmt.Println("recv from client : msgId =", request.GetMsgId(), "data=", string(request.GetData()))
	if err := request.GetConnection().SendMsg(0, []byte("ping...ping...ping...")); err != nil {
		fmt.Println(err)
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("---->DoConnectionBegin is Call")
	if err := conn.SendMsg(202, []byte("DoConnectionBegin BEGIN")); err != nil {
		fmt.Println(err)
	}
	conn.SetProperty("Name","JK_YANG")
}

func DoConnectionStop(conn ziface.IConnection) {
	fmt.Println("--->DoConnectionStop is Call")
	fmt.Println("conn ID=", conn.GetConnID(), "is Lost....")

	if name ,err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name--->",name)
	}
}

func main() {
	//1 SERVER TCP 句柄
	s := znet.NewServer("[JK_YANG TCP SERVER V0.10.0]")

	//2广播 创建连接 与关闭 方法 处理  注册 与调用处理
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionStop)

	//3.添加路由
	s.AddRouter(0, &PintRouter{})

	//4 服务启动 支持连接池与单个链接的处理  根据配置文件
	s.Serve()
}
