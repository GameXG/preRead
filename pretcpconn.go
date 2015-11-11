package preread
import "net"

// 预读 TCP 连接
// 注意：虽然预读连接有 write 函数，但是 write 函数是直接调用的原始链接的 wrtie 函数，意味着回退功能对于write无效！
type PreTcpConn interface {
	TcpConn
	// 开启一层预读
	NewPre() error
	// 关闭一层预读
	// 不会移动当前 offset 。
	// 最后一层预读也被关闭并数据读取完毕时将清空缓冲区。
	ClosePre() error
	//复位预读偏移
	ResetPreOffset() error
}

// TCP 连接 接口
// 注意：虽然预读连接有 write 函数，但是 write 函数是直接调用的原始链接的 wrtie 函数，意味着回退功能对于write无效！
type TcpConn interface {
	net.Conn

	//SetLinger设定当连接中仍有数据等待发送或接受时的Close方法的行为。
	//如果sec < 0（默认），Close方法立即返回，操作系统停止后台数据发送；如果 sec == 0，Close立刻返回，操作系统丢弃任何未发送或未接收的数据；如果sec > 0，Close方法阻塞最多sec秒，等待数据发送或者接收，在一些操作系统中，在超时后，任何未发送的数据会被丢弃。
	SetLinger(sec int) error

	// SetNoDelay设定操作系统是否应该延迟数据包传递，以便发送更少的数据包（Nagle's算法）。默认为真，即数据应该在Write方法后立刻发送。
	SetNoDelay(noDelay bool) error

	//SetReadBuffer设置该连接的系统接收缓冲
	SetReadBuffer(bytes int) error

	//SetWriteBuffer设置该连接的系统发送缓冲
	SetWriteBuffer(bytes int) error
}

type preTcpConn struct {
	TcpConn
	pr PreReader
}


func NewPreTcpConn(conn TcpConn) PreTcpConn{
	return &preTcpConn{conn, NewPreReader(conn)}
}

// 开启一层预读
func (pc *preTcpConn) NewPre() error {
	return pc.pr.NewPre()
}
// 关闭一层预读
// 不会移动当前 offset 。
// 最后一层预读也被关闭并数据读取完毕时将清空缓冲区。
func (pc *preTcpConn) ClosePre() error {
	return pc.pr.ClosePre()
}
//复位预读偏移
func (pc *preTcpConn) ResetPreOffset() error {
	return pc.pr.NewPre()
}

func (pc *preTcpConn) Read(b []byte) (n int, err error) {
	return pc.pr.Read(b)
}
