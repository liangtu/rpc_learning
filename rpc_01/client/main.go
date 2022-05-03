package main

import (
	"fmt"
	"log"
	"net/rpc"
)

//矩形对象
type Ract struct {
}

type Params struct {
	Heigh int
	With  int
}

func main() {

	//1.连接远程RPC
	rp, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")

	if err != nil {
		log.Fatal(err)
	}
	res := 0 //接受结果
	//2.调用远程方法
	err = rp.Call("Rect.Area", Params{50, 30}, &res)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("面积：", res)

}
