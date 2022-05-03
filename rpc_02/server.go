package rpc

import (
	"fmt"
	"net"
	"reflect"
)

type Server struct {
	addr  string
	funcs map[string]reflect.Value
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, funcs: make(map[string]reflect.Value)}
}

//服务端注册方法
func (s *Server) Register(rpcName string, f interface{}) {
	if _, ok := s.funcs[rpcName]; !ok {
		return
	}
	fVal := reflect.ValueOf(f)
	s.funcs[rpcName] = fVal
}

//服务端等待调用
func (s *Server) Run() {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		fmt.Println("监听出错...err=", err)
		return
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("Accept...err=", err)
		}
		srvSession := NewSession(conn)
		b, err := srvSession.Read()
		if err != nil {
			fmt.Println("Accept...err=", err)
			return
		}
		//对数据解码
		rpcData, err := decode(b)
		//根据读取到的数据的name，得到调用的函数名
		f, ok := s.funcs[rpcData.Name]
		if !ok {
			fmt.Println("函数不存在...err=", err)
			return
		}
		//解析遍历客户端传来的参数，放到一个数组中
		inArgs := make([]reflect.Value, 0, len(rpcData.Args))
		for _, arg := range rpcData.Args {
			inArgs = append(inArgs, reflect.ValueOf(arg))
		}
		//动态调用
		out := f.Call(inArgs)

		outArgs := make([]interface{}, 0, len(out))

		for _, o := range out {
			outArgs = append(outArgs, o.Interface())
		}
		//包装数据，返回个客户端
		respRPCData := RPCData{rpcData.Name, outArgs}

		respBytes, err := encode(respRPCData)
		if err != nil {
			fmt.Println("encode...err=", err)
			return
		}
		err = srvSession.Write(respBytes)
		if err != nil {
			fmt.Println("srvSession_Write...err=", err)
			return
		}
	}
}
