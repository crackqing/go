package znet

import "zinx/ziface"


type BaseRouter struct {
}



//template -> before HOOK
func (b *BaseRouter) PreHandle(request ziface.IRequest) {

}

func (b *BaseRouter) Handle(request ziface.IRequest){

}

//template -> after
func (b *BaseRouter) PostHandle(request ziface.IRequest) {

}