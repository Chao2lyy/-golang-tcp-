package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

/*
	存储全局参数，供其他模块使用
*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer //全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口
	Name      string         //当前服务器的名称

	/*
		Zinx
	*/
	Version          string //版本号
	MaxConn          int    //当前服务器主机允许的最大链接数
	MaxPackageSize   uint32 //当前数据包的最大值
	WorkerPoolSize   uint32 //工作池数量
	MaxWorkerTaskLen uint32 //最多worker队列等待长度
}

/*
定义一个全局的对外Globalobj
*/
var GlobalObject *GlobalObj

/*
从zinx.json去加载自定义参数
*/
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//数据解析
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
	提供一个init方法，初始化当前的GlobalObject
*/

func init() { //init方法导包就会调用
	//如果配置文件没有加载，默认
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "1.0",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	//从conf/zinx.json加载
	GlobalObject.Reload()
}
