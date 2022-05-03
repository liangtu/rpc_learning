package rpc

import (
	"fmt"
	"net"
	"reflect"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

//实现通用的RPC客户端
//绑定RPC访问的方法
//传入访问的函数名
//函数具体实现在Server端，Client只有原型函数
//使用makeFunc()完成原型到函数的调用

//fPtr指向函数原型
//xxx.callRPC("queryUser",&query)

func (c *Client) callRPC(rpcName string, fPtr interface{}) {
	//通过反射，获取fPtr未初始化的函数原型
	fn := reflect.ValueOf(fPtr).Elem()
	fmt.Println(fn)
	//另一个函数，作用是对第一个函数参数操作
	//完成与Server的交互
	f := func(args []reflect.Value) []reflect.Value {
		//处理输入的参数
		inArgs := make([]interface{}, 0, len(args))
		for _, arg := range args {
			inArgs = append(inArgs, arg.Interface())
		}

		cliSession := NewSession(c.conn)
		reqRPC := RPCData{rpcName, inArgs}

		//对数据解码
		b, err := encode(reqRPC)
		if err != nil {
			panic(err)
		}
		err = cliSession.Write(b)
		if err != nil {
			panic(err)
		}
		respBytes, err := cliSession.Read()
		if err != nil {
			panic(err)
		}
		//解码数据
		respRPC, err := decode(respBytes)
		if err != nil {
			panic(err)
		}
		outArgs := make([]reflect.Value, 0, len(respRPC.Args))
		for i, arg := range respRPC.Args {
			//必须进行nil转换
			if arg == nil {
				//必须填充一个真正的类型，不能是nil
				outArgs = append(outArgs, reflect.Zero(fn.Type().Out(i)))
				continue
			}
			outArgs = append(outArgs, reflect.ValueOf(arg))
		}
		return outArgs
	}
	//参数1：一个未初始化函数额方法值，类型是reflect.Type
	//参数2：另一个函数，作用是对第一个函数参数操作
	//返回 reflect.Value 类型
	//MakeFunc 使用传入的函数原型，创建一个绑定 参数2 的新函数
	v := reflect.MakeFunc(fn.Type(), f)
	fmt.Println(fn)
	fn.Set(v)
}
