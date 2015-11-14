package preread

import (
	"bytes"
	"io"
	"testing"
)

func TesPreReadt(t *testing.T) {
	r := bytes.NewReader([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	p := NewPreReader(r)

	//比较具体实现结构
	pp, ok := p.(*preRead)
	if ok != true {
		t.Error("1-tpye")
	}

	//读测试
	b := make([]byte, 3)
	if n, err := p.Read(b); n != 3 || err != nil {
		t.Error("3-read")
	}
	if bytes.Equal(b, []byte{0, 1, 2}) != true {
		t.Error("5-read")
	}

	//错误关闭预读测试
	if p.ClosePre() == nil {
		t.Error("7-close")
	}

	//开启预读测试
	if p.NewPre() != nil {
		t.Error("9-NewPre")
	}
	if len(pp.po) != 1 || pp.po[0] != 0 {
		t.Error("11-NewPre", pp.po)
	}

	//开启预读后，读测试
	if n, err := p.Read(b); n != 3 || err != nil {
		t.Error("13-read")
	}
	if bytes.Equal(b, []byte{3, 4, 5}) != true {
		t.Error("15-read")
	}
	if len(pp.po) != 1 || pp.po[0] != 0 {
		t.Error("17-NewPre", pp.po)
	}

	//复位预读
	if p.ResetPreOffset() != nil {
		t.Error("19-reset")
	}
	if len(pp.po) != 1 || pp.po[0] != 0 {
		t.Error("21-NewPre")
	}

	//复位预读后再次读测试
	if n, err := p.Read(b); n != 3 || err != nil {
		t.Error("23-read")
	}
	if bytes.Equal(b, []byte{3, 4, 5}) != true {
		t.Error("25-read")
	}

	// 二次开启预读
	if p.NewPre() != nil {
		t.Error("27-NewPre")
	}
	if len(pp.po) != 2 || pp.po[1] != 3 {
		t.Error("29-NewPre")
	}

	//开启预读后，读测试
	b = make([]byte, 10)
	if n, err := p.Read(b); n != 4 || err != nil {
		t.Error("31-read")
	}
	if bytes.Equal(b[:4], []byte{6, 7, 8, 9}) != true {
		t.Error("33-read")
	}
	if len(pp.po) != 2 || pp.po[1] != 3 {
		t.Error("37-NewPre", pp.po)
	}

	// 读结尾
	if n, err := p.Read(b); n != 0 || err != io.EOF {
		t.Error("39-read eof")
	}

	//复位预读
	if p.ResetPreOffset() != nil {
		t.Error("41-reset")
	}
	if len(pp.po) != 2 || pp.po[1] != 3 {
		t.Error("43-NewPre", pp.po)
	}

	//复位预读后再次读测试
	if n, err := p.Read(b); n != 4 || err != nil {
		t.Error("45-read")
	}
	if bytes.Equal(b[:4], []byte{6, 7, 8, 9}) != true {
		t.Error("47-read")
	}
	if len(pp.po) != 2 || pp.po[1] != 3 {
		t.Error("49-NewPre", pp.po)
	}

	// 关闭一层预读
	if p.ClosePre() != nil {
		t.Error("51-closepre")
	}

	// 读结尾
	if n, err := p.Read(b); n != 0 || err != io.EOF {
		t.Error("53-read eof")
	}

	//复位预读
	if p.ResetPreOffset() != nil {
		t.Error("57-reset")
	}
	if len(pp.po) != 1 || pp.po[0] != 0 {
		t.Error("59-NewPre")
	}

	// 读测试
	b = make([]byte, 6)
	if n, err := p.Read(b); n != 6 || err != nil {
		t.Error("61-read")
	}
	if bytes.Equal(b[:6], []byte{3, 4, 5, 6, 7, 8}) != true {
		t.Error("63-read", b)
	}
	if len(pp.po) != 1 || pp.po[0] != 0 {
		t.Error("67-NewPre", pp.po)
	}

	//复位预读
	if p.ResetPreOffset() != nil {
		t.Error("69-reset")
	}
	if len(pp.po) != 1 || pp.po[0] != 0 {
		t.Error("71-NewPre")
	}

	//复位预读后再次读测试
	b = make([]byte, 3)
	if n, err := p.Read(b); n != 3 || err != nil {
		t.Error("83-read")
	}
	if bytes.Equal(b, []byte{3, 4, 5}) != true {
		t.Error("75-read")
	}

	// 关闭最后一层预读
	if p.ClosePre() != nil {
		t.Error("51-closepre")
	}
	if len(pp.po) != 0 {
		t.Error("59-NewPre")
	}

	// 读测试
	b = make([]byte, 1)
	if n, err := p.Read(b); n != 1 || err != nil {
		t.Error("31-read", n)
	}
	if bytes.Equal(b[:1], []byte{6}) != true {
		t.Error("33-read", b)
	}
	if len(pp.po) != 0 {
		t.Error("37-NewPre")
	}

	// 再次开启预读
	if err := p.NewPre(); err != nil {
		t.Error("newpre")
	}

	// 读测试
	b = make([]byte, 3)
	if n, err := p.Read(b); n != 3 || err != nil {
		t.Error("31-read", n)
	}
	if bytes.Equal(b[:3], []byte{7, 8, 9}) != true {
		t.Error("33-read", b)
	}
	if len(pp.po) != 1 || pp.po[0] == 7 {
		t.Error("37-NewPre")
	}

	//复位预读
	if p.ResetPreOffset() != nil {
		t.Error("69-reset")
	}
	if len(pp.po) != 1 || pp.po[0] != 0 {
		t.Error("71-NewPre", pp.po)
	}

	//复位预读后再次读测试
	b = make([]byte, 10)
	if n, err := p.Read(b); n != 3 || err != nil {
		t.Error("83-read")
	}
	if bytes.Equal(b[:3], []byte{7, 8, 9}) != true {
		t.Error("75-read", b)
	}

	// 读结尾
	if n, err := p.Read(b); n != 0 || err != io.EOF {
		t.Error("53-read eof")
	}
}
