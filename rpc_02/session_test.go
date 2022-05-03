package rpc

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

//测试读写

func TestSession_ReadWrite(t *testing.T) {

	addr := "127.0.0.1:9000"
	my_data := "hello"
	//等待组
	wg := sync.WaitGroup{}
	wg.Add(2)

	//一个协程写一个协程读
	go func() {
		defer wg.Done()
		//创建tcp连接
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		conn, _ := lis.Accept()
		s := Session{Conn: conn}
		err = s.Write([]byte(my_data))
		if err != nil {
			t.Fatal(err)
		}
	}()
	go func() {
		defer wg.Done()
		//创建tcp连接
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		s := Session{Conn: conn}
		data, err := s.Read()
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != my_data {
			t.Fatal(err)
		}
		fmt.Println(string(data))
	}()
	wg.Wait()
}
