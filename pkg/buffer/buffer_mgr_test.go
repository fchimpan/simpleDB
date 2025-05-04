package buffer

import (
	"testing"
	"time"

	"github.com/fchimpan/simpleDB/pkg/file"
	"github.com/stretchr/testify/require"
)

func TestBufferManager(t *testing.T) {
	maxWaitTime := 500 * time.Millisecond

	t.Run("Basic operations", func(t *testing.T) {
		dir := t.TempDir()
		fm := setupFileMgr(t, dir)
		lm := setupLogMgr(t, fm)

		numBuffers := 3
		bm := NewBufferManager(fm, lm, numBuffers, maxWaitTime)

		require.Equal(t, numBuffers, bm.Available(), "初期状態では全てのバッファが利用可能であること")

		// いくつかのブロックをピン留め
		block1 := file.NewBlockID(testLogFile, 1)
		block2 := file.NewBlockID(testLogFile, 2)
		block3 := file.NewBlockID(testLogFile, 3)

		buf1, err := bm.Pin(block1)
		require.NoError(t, err, "Pin should not return an error")
		require.Equal(t, numBuffers-1, bm.Available(), "バッファをピン留めした後は利用可能数が減少する")

		_, err = bm.Pin(block2)
		require.NoError(t, err)
		require.Equal(t, numBuffers-2, bm.Available())

		_, err = bm.Pin(block3)
		require.NoError(t, err)
		require.Equal(t, numBuffers-3, bm.Available())

		// すべてのバッファが使用中の場合、新たなピン留めはエラーになる（タイムアウト）
		block4 := file.NewBlockID(testLogFile, 4)
		_, err = bm.Pin(block4)
		require.Error(t, err, "利用可能なバッファがない場合はエラーになる")

		// バッファを解放して再度テスト
		bm.Unpin(buf1)
		require.Equal(t, 1, bm.Available(), "バッファをアンピンすると利用可能数が増加する")

		// 解放したので再度ピン留めできるはず
		_, err = bm.Pin(block4)
		require.NoError(t, err, "アンピン後は再度ピン留めできる")
		require.Equal(t, 0, bm.Available())

		// 同じブロックを再度ピン留めした場合
		buf4, err := bm.Pin(block4)
		require.NoError(t, err, "既存のブロックはバッファプールから見つかるべき")
		require.NotNil(t, buf4, "buf4 should not be nil")
		// 既にピン留めされているので利用可能数は変わらない
		require.Equal(t, 0, bm.Available())
	})

	t.Run("Transaction flush", func(t *testing.T) {
		t.Skip()
		dir := t.TempDir()
		fm := setupFileMgr(t, dir)
		lm := setupLogMgr(t, fm)

		numBuffers := 3
		bm := NewBufferManager(fm, lm, numBuffers, maxWaitTime)

		// いくつかのブロックをピン留め
		block1 := file.NewBlockID("testfile", 1)
		block2 := file.NewBlockID("testfile", 2)

		buf1, err := bm.Pin(block1)
		require.NoError(t, err)

		buf2, err := bm.Pin(block2)
		require.NoError(t, err)

		// トランザクション1で修正
		txnum := 1
		lsn := 100
		buf1.SetModifyingTx(txnum, lsn)

		// トランザクション2で修正
		buf2.SetModifyingTx(2, 200)

		// トランザクション1に関連するバッファのみフラッシュ
		bm.FlushAll(txnum)

		// buf1はフラッシュされたのでModifyingTxは-1になっているはず
		require.Equal(t, -1, buf1.ModifyingTx(), "フラッシュ後はModifyingTxが-1になる")

		// buf2はフラッシュされていないのでModifyingTxは2のまま
		require.Equal(t, 2, buf2.ModifyingTx(), "他のトランザクションのバッファはフラッシュされない")
	})

	t.Run("Buffer reuse", func(t *testing.T) {
		t.Skip()
		dir := t.TempDir()
		fm := setupFileMgr(t, dir)
		lm := setupLogMgr(t, fm)

		numBuffers := 2
		bm := NewBufferManager(fm, lm, numBuffers, maxWaitTime)

		// 2つのブロックをピン留め
		block1 := file.NewBlockID("testfile", 1)
		block2 := file.NewBlockID("testfile", 2)

		buf1, err := bm.Pin(block1)
		require.NoError(t, err)

		_, err = bm.Pin(block2)
		require.NoError(t, err)

		// バッファ1をアンピン
		bm.Unpin(buf1)

		// 新しいブロックをピン留め（バッファ1が再利用されるはず）
		block3 := file.NewBlockID("testfile", 3)
		buf3, err := bm.Pin(block3)
		require.NoError(t, err)

		// buf1とbuf3が同じバッファを指しているはず（ただしコンテンツは異なる）
		require.True(t, buf1 == buf3, "アンピンされたバッファが再利用されるべき")
		require.True(t, buf3.block.Equals(block3), "再利用されたバッファには新しいブロックがロードされる")
	})
}
