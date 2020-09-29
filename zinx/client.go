package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	fmt.Println("CLIENT START....")
	time.Sleep(1*time.Second)
	 conn,err := net.Dial("tcp","127.0.0.1:8889")
	if err != nil {
		fmt.Println("client start err exit!!!")
		return
	}
	for n:=300; n > 0; n--  {
		dp := znet.NewDataPack()
		msg,_ := dp.Pack(znet.NewMessage(0,[]byte("ZInx client Test Message")))
		if  _,err :=  conn.Write(msg); err != nil {
			fmt.Println("write error err",err)
			return
		}
		//读取head部分 [MsgLength]|[MsgID]|[Data]的封包格式
		headData := make([]byte,dp.GetHeadLen())
		if _,err := io.ReadFull(conn,headData); err != nil {
			fmt.Println("read head error")
			break
		}
		msgHead,err :=  dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:",err )
		}
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte,msg.GetMsgLen())
			if  _,err := io.ReadFull(conn,msg.Data); err != nil {
				fmt.Println("server unpack data err:",err)
				return
			}
			fmt.Println("==> Recv Msg: ID=",msg.Id,"len ==>",msg.DataLen,"data=",string(msg.Data))
		}
		time.Sleep(1*time.Second)

	}




}

