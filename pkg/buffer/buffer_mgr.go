package buffer

import (
	"errors"
	"sync"
	"time"

	"github.com/fchimpan/simpleDB/pkg/file"
	"github.com/fchimpan/simpleDB/pkg/log"
)

type BufferManager struct {
	mu           sync.Mutex
	bufferPool   []*Buffer
	numAvailable int
	maxWaitTime  time.Duration
}

func NewBufferManager(fm *file.FileMgr, lm *log.LogMgr, numBuffers int, maxWaitTime time.Duration) *BufferManager {
	bufferPool := make([]*Buffer, numBuffers)
	for i := 0; i < numBuffers; i++ {
		bufferPool[i] = NewBuffer(fm, lm)
	}
	return &BufferManager{
		bufferPool:   bufferPool,
		numAvailable: numBuffers,
		maxWaitTime:  maxWaitTime,
	}
}

func (bm *BufferManager) Available() int {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	return bm.numAvailable
}

func (bm *BufferManager) Pin(block *file.BlockID) (*Buffer, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	start := time.Now()
	for {
		buf, err := bm.tryToPin(block)
		if err != nil {
			return nil, err
		}
		if buf != nil {
			return buf, nil
		}

		if time.Since(start) > bm.maxWaitTime {
			return nil, errors.New("BufferAbortException: timeout while waiting for buffer")
		}

		bm.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
		bm.mu.Lock()
	}
}

func (bm *BufferManager) tryToPin(block *file.BlockID) (*Buffer, error) {
	// ブロックがすでにバッファにあれば再利用
	if buf := bm.findExistingBuffer(block); buf != nil {
		if !buf.IsPinned() {
			bm.numAvailable--
		}
		buf.Pin()
		return buf, nil
	}
	// 使えるバッファがなければ nil
	buf := bm.chooseUnpinnedBuffer()
	if buf == nil {
		return nil, nil
	}

	// unpin されているバッファを使う
	if err := buf.ReadFromBlock(block); err != nil {
		return nil, err
	}
	buf.Pin()
	bm.numAvailable--
	return buf, nil

}

// バッファ内に既にあるページを探す
func (bm *BufferManager) findExistingBuffer(block *file.BlockID) *Buffer {
	for _, buf := range bm.bufferPool {
		if buf != nil && buf.block != nil && buf.block.Equals(block) {
			return buf
		}
	}
	return nil
}

// TODO: ナイーブな実装としているので、LRU のようなアルゴリズムにしてみる
func (bm *BufferManager) chooseUnpinnedBuffer() *Buffer {
	for _, buf := range bm.bufferPool {
		if buf != nil && !buf.IsPinned() {
			return buf
		}
	}
	return nil
}

func (bm *BufferManager) Unpin(buf *Buffer) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	buf.Unpin()
	if !buf.IsPinned() {
		bm.numAvailable++
	}
}

// 指定されたトランザクションIDによって変更されたすべてのバッファページをフラッシュ（ログ書き＋ディスク書き）
func (bm *BufferManager) FlushAll(txnum int) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	for _, buf := range bm.bufferPool {
		if buf.ModifyingTx() == txnum {
			if err := buf.Flush(); err != nil {
				return err
			}
		}
	}
	return nil
}
