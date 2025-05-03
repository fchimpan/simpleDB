package file

import (
	"bytes"
	"encoding/binary"
	"unicode/utf8"
)

const (
	HeaderSize = 4
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
	return int(binary.BigEndian.Uint32(p.bb.Bytes()[offset : offset+HeaderSize]))
}

func (p *Page) SetInt(offset int, value int) {
	binary.BigEndian.PutUint32(p.bb.Bytes()[offset:offset+HeaderSize], uint32(value))
}

// 先頭4バイトはデータ長
func (p *Page) GetBytes(offset int) []byte {
	len := p.GetInt(offset)
	s := offset + HeaderSize
	e := s + len
	return p.bb.Bytes()[s:e]
}

// GetBytesAtPosition は特定の位置から指定された長さのバイト列を取得します
func (p *Page) GetBytesAtPosition(pos int, length int) []byte {
	return p.bb.Bytes()[pos : pos+length]
}

func (p *Page) SetBytes(offset int, b []byte) {
	p.SetInt(offset, len(b))
	copy(p.bb.Bytes()[offset+HeaderSize:], b)
}

func (p *Page) GetString(offset int) string {
	return string(p.GetBytes(offset))
}

func (p *Page) SetString(offset int, s string) {
	p.SetBytes(offset, []byte(s))
}

func MaxLength(strlen int) int {
	bytesPerChar := utf8.UTFMax
	return HeaderSize + strlen*bytesPerChar // 4 バイトは長さ用
}

func (p *Page) Contents() *bytes.Buffer {
	return p.bb
}
