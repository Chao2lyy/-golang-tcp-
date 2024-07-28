package znet

type Message struct {
	Id      uint32 //消息的id
	Datalen uint32 //消息长度
	Data    []byte //消息内容
}

// 创建
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		Datalen: uint32(len(data)),
		Data:    data,
	}
}

// 消息id
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// 长度
func (m *Message) GetMsgLen() uint32 {
	return m.Datalen
}

// 获取消息
func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}
func (m *Message) SetDataLen(len uint32) {
	m.Datalen = len
}
