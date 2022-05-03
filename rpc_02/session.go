package rpc

import (
	"encoding/binary"
	"io"
	"net"
)

//测试会话中数据读写
//会话连接的结构体

type Session struct {
	Conn net.Conn
}

//创建连接
func NewSession(conn net.Conn) *Session {
	return &Session{conn}
}

//连接中写数据
func (s *Session) Write(data []byte) error {
	//4字节+数据长的切片
	buf := make([]byte, 4+len(data))
	//写入头部数据，记录长度
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data)))
	//写入数据
	copy(buf[4:], data)
	_, err := s.Conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

//连接中读数据
func (s *Session) Read() ([]byte, error) {
	//读取头部长度
	header := make([]byte, 4)
	//按头部长度，读取头部数据
	_, err := io.ReadFull(s.Conn, header)
	if err != nil {
		return nil, err
	}
	//读取数据长度
	dataLen := binary.BigEndian.Uint32(header)
	//按照数据长度去读取数据
	data := make([]byte, dataLen)
	_, err = io.ReadFull(s.Conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
