package ziface

type IConnManager interface {
	//ADD
	Add(conn IConnection)
	//DELETE
	Remove(conn IConnection)
	 //CONN_ID
	Get(connID uint32) (IConnection, error)
	 //CURRENT_COUNT
	Len() int
	 //CLEAR CONN ALL
	 ClearConn()
}