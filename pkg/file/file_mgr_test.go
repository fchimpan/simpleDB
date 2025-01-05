package file

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileMgr(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	blockSize := 64
	fm, err := NewFileMgr(tmpDir, blockSize)
	if err != nil {
		t.Fatalf("failed to create FileMgr: %v", err)
	}

	t.Run("Write and read", func(t *testing.T) {
		t.Parallel()
		fileName := "WriteAndReadTestfile"
		data := make([]byte, blockSize)
		copy(data, []byte("Testing write and read"))
		page := NewPageFromBytes(data)
		blkID := NewBlockID(fileName, 0)

		if err := fm.Write(blkID, page); err != nil {
			t.Fatalf("Write failed: %v", err)
		}

		readPage := NewPage(blockSize)
		if err := fm.Read(blkID, readPage); err != nil {
			t.Fatalf("Read failed: %v", err)
		}

		assert.True(t, bytes.Equal(page.Contents().Bytes(), readPage.Contents().Bytes()))
	})

	t.Run("Append", func(t *testing.T) {
		t.Parallel()
		fileName := "AppendTestfile"
		blkID, err := fm.Append(fileName)
		require.NoError(t, err, "Append failed")

		blockNum, err := fm.blockNum(fileName)
		require.NoError(t, err, "blockNum failed")
		assert.Equal(t, 1, blockNum, "expected blockNum to be 1")

		assert.Equal(t, 0, blkID.Number(), "expected block number to be 0")
		assert.Equal(t, fileName, blkID.FileName(), "expected file name to match")
	})
}
