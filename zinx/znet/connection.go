package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	//当前conn属于那个server
	TcpServer ziface.IServer
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool
	handleAPI ziface.HandleFunc
	ExitChan  chan bool
	//线程池 协程池
	MsgHandle ziface.IMsgHandle
	//读写 Goroutine
	msgChan chan []byte

	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func (c *Connection) StartReader() {
	fmt.Println("StartReader Goroutine is running...", c.ConnID)
	defer fmt.Println("connID=", c.ConnID, "Reader is exit,remote addr is", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//创建一个拆包与解包对像
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTcpConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}
		//拆包 行到MsgID 与 MsgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//开启池子使用连接池而不是对应一个链接开启一个,优化处理
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			//执行注册路由注册
			go c.MsgHandle.DoMsgHandler(&req)
		}

	}

}

//write Goroutine
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine  is running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[CONN WRITER EXIT!]")

	//select
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error", err)
				return
			}
		case <-c.ExitChan:
			return
		}

	}
}

//启动
func (c *Connection) Start() {
	fmt.Println("Conn Start().. ConnID=", c.ConnID)

	go c.StartReader()
	go c.StartWriter()
	//处理路由之前调用 HOOK函数
	c.TcpServer.CallOnConnStart(c)

}

//停止
func (c *Connection) Stop() {
	fmt.Println("Conn stop().. ConnID=", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//处理关闭前 关闭HOOK处理
	c.TcpServer.CallOnConnStop(c)
	_ = c.Conn.Close()
	//writer close goroutine
	c.ExitChan <- true

	close(c.ExitChan)
	close(c.msgChan)

	c.TcpServer.GetConnMgr().Remove(c)

}

//获取当前链接绑定的socket conn
func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP 状态 IP PORT
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据 ,将数据发送远程客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {

	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id", msgId)
		return errors.New("pack error msg")
	}
	//writer goroutine
	c.msgChan <- binaryMsg

	return nil

}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		MsgHandle: handler,
		msgChan:   make(chan []byte),
		ExitChan:  make(chan bool, 1),
		property:make(map[string]interface{}),
	}
	//1对1的服务绑定
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value,ok := c.property[key] ; ok {
		return value,nil
	}
	return nil,errors.New("not property found")

}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property,key)
}
