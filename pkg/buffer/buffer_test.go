package buffer

import (
	"testing"

	"github.com/fchimpan/simpleDB/pkg/file"
	"github.com/stretchr/testify/require"
)

func TestBuffer_BasicOperations(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		dir := t.TempDir()
		fm := setupFileMgr(t, dir)
		lm := setupLogMgr(t, fm)
		buffer := NewBuffer(fm, lm)

		block := file.NewBlockID(testLogFile, 0)
		err := buffer.ReadFromBlock(block)
		require.NoError(t, err, "ReadFromBlock should not return an error")

		buffer.Pin()
		if !buffer.IsPinned() {
			t.Errorf("expected buffer to be pinned")
		}
		buffer.Unpin()
		if buffer.IsPinned() {
			t.Errorf("expected buffer to be unpinned")
		}

		buffer.SetModifyingTx(1, 10)
		if buffer.ModifyingTx() != 1 {
			t.Errorf("expected modifying transaction to be 1, got %d", buffer.ModifyingTx())
		}

		err = buffer.Flush()
		require.NoError(t, err, "Flush should not return an error")

		if buffer.ModifyingTx() != -1 {
			t.Errorf("expected modifying transaction to be -1 after flush, got %d", buffer.ModifyingTx())
		}

	})

}
