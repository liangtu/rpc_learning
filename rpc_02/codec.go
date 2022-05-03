package rpc

import (
	"bytes"
	"encoding/gob"
)

type RPCData struct {
	//访问时的函数
	Name string
	//访问时的参数
	Args []interface{}
}

//编码
func encode(data RPCData) ([]byte, error) {
	var buf bytes.Buffer
	bufEnc := gob.NewEncoder(&buf)
	if err := bufEnc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//解码
func decode(b []byte) (RPCData, error) {
	buf := bytes.NewBuffer(b)
	bufEnc := gob.NewDecoder(buf)
	var data RPCData
	if err := bufEnc.Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}
