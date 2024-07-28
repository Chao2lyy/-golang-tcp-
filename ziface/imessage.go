package ziface

/*
	请求消息封装，定义接口
*/

type IMessage interface {
	//消息id
	GetMsgId() uint32
	//长度
	GetMsgLen() uint32
	//获取消息
	GetData() []byte

	SetMsgId(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
