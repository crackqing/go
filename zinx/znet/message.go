package znet

type Message struct {
	Id      uint32 //消息ID
	DataLen uint32 //消息的长度
	Data    []byte //消息的内容
}

func NewMessage(id uint32,data[]byte) *Message{
	return &Message{
		Id : id ,
		DataLen: uint32(len(data)),
		Data:data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}
func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(Id uint32) {
	m.Id = Id
}
func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}