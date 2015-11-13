package preread
import (
	"testing"
	"net"
	"bytes"
)


func TestrePConn(t *testing.T) {
	c, err := net.Dial("tcp", "www.baidu.com:80")
	if err != nil {
		panic(err)
	}
	if _, err := c.Write([]byte("GET / HTTP/1.0\r\nHOST:www.google.com\r\n\r\n")); err != nil {
		panic(err)
	}
	preconn := NewPreConn(c)
	preconn.NewPre()

	// 第一次读取
	b1 := make([]byte, 16384)
	if n, err := preconn.Read(b1); err != nil {
		panic(err)
	}else {
		b1 = b1[:n]
	}

	preconn.ResetPreOffset()

	//　预读的实现是只要缓冲区有数据，即使不够也会直接返回而不是阻塞在Ｒｅａｄ里面
	//　所以这里应该是相同的长度。
	b2 := make([]byte, 16384)
	if n, err := preconn.Read(b1); err != nil {
		panic(err)
	}else {
		b2 = b2[:n]
	}

	if bytes.Equal(b1, b2) != true {
		t.Errorf("预读会退结果不相等")
	}

	preconn.ResetPreOffset()
	conn := net.Conn(preconn)
	b2 = b2[:16384]
	if n, err := conn.Read(b1); err != nil {
		panic(err)
	}else {
		b2 = b2[:n]
	}

	if bytes.Equal(b1, b2) != true {
		t.Errorf("预读会退结果不相等")
	}
}
