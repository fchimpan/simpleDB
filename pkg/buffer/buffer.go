package buffer

import (
	"sync"

	"github.com/fchimpan/simpleDB/pkg/file"
	"github.com/fchimpan/simpleDB/pkg/log"
)

type Buffer struct {
	mu      sync.Mutex
	content *file.Page // 割り当てられたブロック（nilなら未割当）
	block   *file.BlockID
	pins    int // 現在のピン数（利用中クライアントの数）
	txnum   int // 修正したトランザクション番号、未修正なら -1
	lsn     int // 最終ログシーケンス番号（未記録なら -1）

	fm *file.FileMgr
	lm *log.LogMgr
}

func NewBuffer(fm *file.FileMgr, lm *log.LogMgr) *Buffer {
	return &Buffer{
		content: file.NewPage(fm.BlockSize()),
		block:   nil,
		pins:    0,
		txnum:   -1,
		lsn:     -1,
		fm:      fm,
		lm:      lm,
	}
}

func (b *Buffer) Pin() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.pins++
}

func (b *Buffer) Unpin() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.pins > 0 {
		b.pins--
	}
}

func (b *Buffer) IsPinned() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.pins > 0
}

func (b *Buffer) ModifyingTx() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.txnum
}

func (b *Buffer) SetModifyingTx(txnum, lsn int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.txnum = txnum
	b.lsn = lsn
}

func (b *Buffer) Flush() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.flushInternal()
}

func (b *Buffer) flushInternal() error {
	if b.txnum >= 0 {
		if err := b.lm.Flush(b.lsn); err != nil {
			return err
		}
		if err := b.fm.Write(b.block, b.content); err != nil {
			return err
		}
		b.txnum = -1
	}
	return nil
}

func (b *Buffer) ReadFromBlock(blk *file.BlockID) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err := b.flushInternal(); err != nil {
		return err
	}
	b.block = blk
	if err := b.fm.Read(blk, b.content); err != nil {
		return err
	}
	b.pins = 0
	return nil
}
