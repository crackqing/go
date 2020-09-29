package ziface

type IRouter interface {
	//template -> before
	PreHandle(request IRequest)

	Handle(request IRequest)
	//template -> after
	PostHandle(request IRequest)
}