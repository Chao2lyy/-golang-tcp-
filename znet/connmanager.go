package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

/*
 链接管理模块
*/

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的链接集合
	connLock    sync.RWMutex                  //保护链接集合的读写锁
}

// 创建当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	//加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//conn加入ConnManager
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("conneID=", conn.GetConnID(), "add to ConnManager successfully:conn num =", connMgr.Len())
}

// 删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	//加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除链接
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("conneID=", conn.GetConnID(), "delete from ConnManager successfully:conn num =", connMgr.Len())
}

// 根据id获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND!\n")
	}
}

// 得到当前链接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并终止所有链接
func (connMgr *ConnManager) ClearConn() {
	//加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear all connection! num =", connMgr.Len())
}
