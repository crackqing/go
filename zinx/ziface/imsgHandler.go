package ziface

type IMsgHandle interface {
	DoMsgHandler(request IRequest)

	AddRouter(msgID uint32,router IRouter)

	StartWorkerPool()
	StartOneWorker(workerId int, taskQueue chan IRequest)
	SendMsgToTaskQueue(request IRequest)
}