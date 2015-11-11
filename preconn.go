package preread
import "net"

// 预读连接
// 注意：虽然预读连接有 write 函数，但是 write 函数是直接调用的原始链接的 wrtie 函数，意味着回退功能对于write无效！
type PreConn interface {
	net.Conn
	// 开启一层预读
	NewPre() error
	// 关闭一层预读
	// 不会移动当前 offset 。
	// 最后一层预读也被关闭并数据读取完毕时将清空缓冲区。
	ClosePre() error
	//复位预读偏移
	ResetPreOffset() error
}

type preConn struct {
	net.Conn
	pr PreReader
}


func NewPreConn(conn net.Conn) PreConn{
	return &preConn{conn, NewPreReader(conn)}
}

// 开启一层预读
func (pc *preConn) NewPre() error {
	return pc.pr.NewPre()
}
// 关闭一层预读
// 不会移动当前 offset 。
// 最后一层预读也被关闭并数据读取完毕时将清空缓冲区。
func (pc *preConn) ClosePre() error {
	return pc.pr.ClosePre()
}
//复位预读偏移
func (pc *preConn) ResetPreOffset() error {
	return pc.pr.NewPre()
}

func (pc *preConn) Read(b []byte) (n int, err error) {
	return pc.pr.Read(b)
}
