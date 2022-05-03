package main

import (
	"log"
	"net/http"
	"net/rpc"
)

//矩形对象
type Rect struct {
}

type Params struct {
	Heigh int
	With  int
}

func (r *Rect) Area(p Params, res *int) error {
	*res = p.With * p.Heigh
	return nil
}

func main() {
	//1、注册服务
	rect := new(Rect)
	rpc.Register(rect)
	//2、把服务处理绑定到http协议上
	rpc.HandleHTTP()
	//3、监听服务，等待客户端调用求面积
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
