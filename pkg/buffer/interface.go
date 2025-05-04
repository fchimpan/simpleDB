package buffer

import "github.com/fchimpan/simpleDB/pkg/file"

type BufferManagerInterface interface {
	Available() int
	Pin(block *file.BlockID) (BufferInterface, error)
	Unpin(buf BufferInterface) error
	FlushAll() error // 指定されたトランザクションIDによって変更されたすべてのバッファページをフラッシュ（ログ書き＋ディスク書き）
}

type BufferInterface interface {
	Pin()
	Unpin()
	IsPinned() bool
	SetModifyingTx(txnum, lsn int) // バッファを編集したトランザクションとその修正に対応するログシーケンス番号をセット
	ModifyingTx() int              // バッファを編集したトランザクション番号を取得
	Flush() error
	ReadFromBlock(blk *file.BlockID) // ブロックをpage読み込む
}
