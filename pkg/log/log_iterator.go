package log

import (
	"io"

	"github.com/fchimpan/simpleDB/pkg/file"
)

type LogIterator struct {
	fm         *file.FileMgr
	currentBlk *file.BlockID
	page       *file.Page
	pos        int
	boundary   int
}

func NewLogIterator(fm *file.FileMgr, blk *file.BlockID) (*LogIterator, error) {
	p := file.NewPage(fm.BlockSize())
	if err := fm.Read(blk, p); err != nil {
		return nil, err
	}
	boundary := p.GetInt(0)
	return &LogIterator{
		fm:         fm,
		currentBlk: blk,
		page:       p,
		pos:        boundary,
		boundary:   boundary,
	}, nil
}

func (li *LogIterator) HasNext() bool {
	return li.pos < li.fm.BlockSize() || li.currentBlk.Number() > 0
}

func (li *LogIterator) Next() ([]byte, error) {
	if li.pos >= li.fm.BlockSize() {
		// 次のブロックがないなら終了
		if li.currentBlk.Number() == 0 {
			return nil, io.EOF
		}

		// 1つ前のブロックに切り替え
		newBlk := file.NewBlockID(li.currentBlk.FileName(), li.currentBlk.Number()-1)
		p := file.NewPage(li.fm.BlockSize())
		if err := li.fm.Read(newBlk, p); err != nil {
			return nil, err
		}
		li.currentBlk = newBlk
		li.page = p
		li.boundary = p.GetInt(0)
		li.pos = li.boundary
	}

	// レコード長の読み取り
	recLen := li.page.GetInt(li.pos)
	recData := li.page.GetBytesAtPosition(li.pos+file.HeaderSize, recLen)
	li.pos += file.HeaderSize + recLen
	return recData, nil
}
