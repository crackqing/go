package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
	//负责WORKER取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作WORKER池
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgId(), "is NOT FOUND! Need register")
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api,msgID=" + strconv.Itoa(int(msgID)))
	}
	//2
	mh.Apis[msgID] = router
	fmt.Println("ADD API MsgID = ", msgID, " success!")
}

//START WORKER
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//1
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerPoolSize)
		//2
		go mh.StartOneWorker(i, mh.TaskQueue[i])

	}
}

//START WORKER ING
func (mh *MsgHandle) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerID = ", workerId, "is started...")
	for {
		select {
		//如果有消息过来,出列一个客户端
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//平均分配 %取模处理
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("ADD ConnID  =", request.GetConnection().GetConnID(), "request  MsgID", request.GetMsgId(),
		"to workerId", workerId)
	mh.TaskQueue[workerId] <- request
}
