package preread

import (
	"errors"
	"github.com/gamexg/memstream"
	"io"
)

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

type preRead struct {
	r     io.Reader           // 来源 reader
	ms    memstream.MemStream // peek 数据保存区
	tee   io.Reader           // 同步写 peek
	multi io.Reader           // 优先读 peek (读到尾的会被移出)
	po    []int64             // peek offset 偏移
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

func (pr *preRead) NewPre() (err error) {
	if len(pr.po) == 0 {
		if pr.ms == nil {
			pr.ms = memstream.NewMemStream()
			pr.tee = io.TeeReader(pr.r, pr.ms)
			pr.multi = io.MultiReader(pr.ms, pr.tee)
		} else {
			if err := pr.ms.DeleteRead(); err != nil {
				return errors.New("DeleteRead")
			}
		}
	}
	offset, err := pr.ms.Seek(0, 1)
	if err != nil {
		return errors.New("[PeekReader] Internal error 2")
	}

	pr.po = append(pr.po, offset)

	return nil
}

func (pr *preRead) ClosePre() error {
	if len(pr.po) == 0 {
		return errors.New("There is no pre reading data")
	}
	pr.po = pr.po[:len(pr.po)-1]

	return nil
}

func (pr *preRead) ResetPreOffset() error {
	if len(pr.po) == 0 {
		return errors.New("peek off")
	}
	if _, err := pr.ms.Seek(pr.po[len(pr.po)-1], 0); err != nil {
		return errors.New("[PeekReader] Internal error 1")
	}
	pr.multi = io.MultiReader(pr.ms, pr.tee)
	return nil
}

func (pr *preRead) Read(p []byte) (n int, err error) {
	if pr.ms != nil && len(pr.po) == 0 && pr.ms.Len() == 0 {
		// 预读已关闭 并且 不存在预读数据
		// 清空缓冲区
		pr.ms.Close()
		pr.ms = nil
		pr.tee = nil
		pr.multi = pr.r
	}
	return pr.multi.Read(p)
}
