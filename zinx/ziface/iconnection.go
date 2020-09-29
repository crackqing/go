package ziface

import "net"

type IConnection interface {
	//启动
	Start()
	//停止
	Stop()
	//获取当前链接绑定的socket conn
	GetTcpConnection() *net.TCPConn
	//获取当前链接模块的链接ID
	GetConnID() uint32
	//获取远程客户端的TCP 状态 IP PORT
	RemoteAddr() net.Addr
	//发送数据 ,将数据发送远程客户端
	SendMsg(msgId uint32,data []byte) error

	SetProperty(key string,value interface{})
	GetProperty(key string) (interface{},error)

	RemoveProperty(key string)

}

type HandleFunc func(*net.TCPConn,[]byte,int) error
