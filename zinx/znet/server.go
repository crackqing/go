package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	//SERVER ROUTER BIND MSG_HANDLER
	MsgHandle ziface.IMsgHandle
	//CONN_MANAGER_SET
	ConnManager ziface.IConnManager

	//广播使用
	OnConnStart func(conn ziface.IConnection)
	OnConnStop func(conn ziface.IConnection)
}

//Start
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP: %s,Port %d, is Starting\n", s.IP, s.Port)

	go func() {
		//0 开启工作池 队列worker -> TaskQueue
		s.MsgHandle.StartWorkerPool()

		//1
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error", err)
			return
		}
		//2
		Listen, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err ", err)
			return
		}
		fmt.Println("start ZInx server success", s.Name, " success,Listen...")
		var cid uint32
		cid = 0
		//3
		for {
			conn, err := Listen.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("Too Many Connections MaxConn=",utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}



			dealConn := NewConnection(s,conn, cid, s.MsgHandle)
			cid++
			go dealConn.Start()
		}
	}()

}

//stop server
func (s *Server) Stop() {
	fmt.Println("[STOP ZInx server name]",s.Name)
	s.ConnManager.ClearConn()
}



func (s *Server) Serve() {
	//启动server服务功能
	s.Start()
	//阻塞状态
	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgId, router)
	fmt.Println("Add Router Success!!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnManager
}


//INIT
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return s
}

func (s *Server)  SetOnConnStart(hookFunc  func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}
func (s *Server)  SetOnConnStop(hookFunc func(connection ziface.IConnection)){
	s.OnConnStop  = hookFunc
}

func (s *Server)  CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil  {
		fmt.Println("---> Call OnConnStart()")
		s.OnConnStart(conn)
	}
}
func (s *Server)  CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil  {
		fmt.Println("---> Call OnConnStop()")
		s.OnConnStop(conn)
	}
}
