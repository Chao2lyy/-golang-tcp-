package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

/*
消息处理模块实现
*/

type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter

	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池工作数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中哦获取
	}
}

// 调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1 从Request中找到msgID
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgId(), "is NOT FOUND!")
	}
	//2 调度
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 判断方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		//id已经注册了
		panic("repeat api,msgID =" + strconv.Itoa((int(msgID))))
	}
	//2 添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID=", msgID, "succ!")
}

// 启动一个worker工作池(开启工作池的动作只能发生一次，一个zinx框架只能有一个worker工作池)
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker启动
		//1 当前的worker对应channerl消息队列 开辟空间 第0个worker就第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2 启动当前的worker，阻塞等待消息从channel传递
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个worker工作流程
func (mh *MsgHandle) startOneWorker(workID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID=", workID, "is started")
	//不断阻塞等待对应消息队列消息
	for {
		select {
		//如果有消息过来，出列的就是一个客户端的request，执行当前所绑定业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}

	}
}

// 将消息交给TaskQueue
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//1 平均分配消息给不同的worker
	//根据客户端建立的connID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), "request MsgID=", request.GetMsgId(), "to WorkerID=", workerID)

	//2 将消息发送给对应的worker的taskqueue即可
	mh.TaskQueue[workerID] <- request
}
