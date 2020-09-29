package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理链接信息集合

	connLock sync.RWMutex //读写保护
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

//ADD
func (c *ConnManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnID()] = conn
	fmt.Println("connection add",conn.GetConnID()," to ConnManager  successfully; conn num=", c.Len())
}

//DELETE
func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.connections,conn.GetConnID())
	fmt.Println("connection Remove",conn.GetConnID()," to ConnManager  successfully; conn num=", c.Len())

}

//CONN_ID
func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn,ok := c.connections[connID]; ok {
		return conn,nil
	} else {
		return nil,errors.New("connection not FOUND")
	}


}

//CURRENT_COUNT
func (c *ConnManager) Len() int {
	return len(c.connections)
}

//CLEAR CONN ALL
func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connId,conn := range c.connections {
		//STOP
		conn.Stop()
		//DELETE
		delete(c.connections,connId)
	}
	fmt.Println("Clear All connections success! conn num=",c.Len())
}
