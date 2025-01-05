package file

import (
	"bytes"
	"encoding/binary"
	"unicode/utf8"
)

// A Page object holds the contents of a disk block.
type Page struct {
	bb *bytes.Buffer
}

func NewPage(blocksize int) *Page {
	return &Page{bb: bytes.NewBuffer(make([]byte, blocksize))}
}

func NewPageFromBytes(b []byte) *Page {
	return &Page{bb: bytes.NewBuffer(b)}
}

func (p *Page) GetInt(offset int) int {
	return int(binary.BigEndian.Uint32(p.bb.Bytes()[offset : offset+4]))
}

func (p *Page) SetInt(offset int, value int) {
	binary.BigEndian.PutUint32(p.bb.Bytes()[offset:offset+4], uint32(value))
}

func (p *Page) GetBytes(offset int) []byte {
	len := p.GetInt(offset)
	s := offset + 4
	e := s + len
	return p.bb.Bytes()[s:e]
}

func (p *Page) SetBytes(offset int, b []byte) {
	p.SetInt(offset, len(b))
	copy(p.bb.Bytes()[offset+4:], b)
}

func (p *Page) GetString(offset int) string {
	return string(p.GetBytes(offset))
}

func (p *Page) SetString(offset int, s string) {
	p.SetBytes(offset, []byte(s))
}

func MaxLength(strlen int) int {
	bytesPerChar := utf8.UTFMax
	return 4 + strlen*bytesPerChar // 4 バイトは長さ用
}

func (p *Page) Contents() *bytes.Buffer {
	return p.bb
}
