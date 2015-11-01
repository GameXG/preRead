# preRead
golang 无限 预读 Read 。

解析协议时经常出现一个协议不匹配需要切换到下一个协议尝试解析，但是前一个协议已经读取了一部分字节，标准的 io.net 只能回退一个 byte ，很多时候都不够用。

本库可以提供无限长度的回退功能，并且支持多层回退。

## 实现的借口

```go

// 预读接口
// 不是线程安全的！
// 允许多次多层打开关闭预读。
type PreReader interface {
	io.Reader

	// 开启一层预读
	NewPre() error
	// 关闭一层预读
	// 不会移动当前 offset 。
	// 最后一层预读也被关闭并数据读取完毕时将清空缓冲区。
	ClosePre() error
	//复位预读偏移
	ResetPreOffset() error
}

// NewPeekReader 新建预读
// 默认不会开启预读功能
func NewPreReader(r io.Reader) PreReader {
	pr := preRead{}
	pr.r = r
	pr.ms = nil
	pr.tee = nil
	pr.multi = r
	pr.po = make([]int64, 0, 5)
	return &pr
}

````