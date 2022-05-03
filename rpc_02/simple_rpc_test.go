package rpc

import (
	"encoding/gob"
	"fmt"
	"net"
	"testing"
)

type User struct {
	Name string
	Age  int
}

func queryUser(uid int) (User, error) {
	user := make(map[int]User)
	user[0] = User{"zs", 20}
	user[1] = User{"dd", 4}
	user[2] = User{"zdfs", 54}
	user[3] = User{"zsvv", 23}
	if u, ok := user[uid]; ok {
		return u, nil
	}
	return User{}, fmt.Errorf("id %d not int user db", uid)
}

func TestRPC(t *testing.T) {
	gob.Register(User{})
	addr := "127.0.0.1:9001"
	srv := NewServer(addr)

	srv.Register("queryUser", queryUser)
	go srv.Run()

	//客户端获取连接
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
	}

	cli := NewClient(conn)

	var query func(int) (User, error)

	cli.callRPC("queryUser", &query)

	u, err := query(1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(u)

}
