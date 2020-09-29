package ziface

type IServer interface {
	//start
	Start()
	//stop
	Stop()
	//run
	Serve()

	//router
	AddRouter(msgId uint32,router IRouter)

	GetConnMgr() IConnManager

	SetOnConnStart(func(connection IConnection))
	SetOnConnStop(func(connection IConnection))

	CallOnConnStart(connection IConnection)
	CallOnConnStop(connection IConnection)
}
