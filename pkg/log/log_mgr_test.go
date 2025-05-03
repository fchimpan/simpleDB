package log

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fchimpan/simpleDB/pkg/file"
	"github.com/stretchr/testify/assert"
)

func TestLogMgr_AppendAndFlush(t *testing.T) {
	dir := t.TempDir()
	fm := setupFileMgr(t, dir)
	lm := setupLogMgr(t, fm)

	lsn1, err := lm.Append([]byte("record1"))
	assert.NoError(t, err)
	assert.Equal(t, 1, lsn1)

	lsn2, err := lm.Append([]byte("record2"))
	assert.NoError(t, err)
	assert.Equal(t, 2, lsn2)

	lsn3, err := lm.Append([]byte("record3"))
	assert.NoError(t, err)
	assert.Equal(t, 3, lsn3)

	err = lm.Flush(lsn3)
	assert.NoError(t, err)

	readPage := file.NewPage(blockSize)
	blkID := file.NewBlockID(testLogFile, 0)
	err = fm.Read(blkID, readPage)
	assert.NoError(t, err)

	readFile, err := os.ReadFile(filepath.Join(dir, testLogFile))
	assert.NoError(t, err)
	assert.Equal(t, readFile, readPage.Contents().Bytes())

	// ブロックサイズを超えるレコードを追加
	lns4, err := lm.Append([]byte("overflow length records!!!!!!!!!!"))
	assert.NoError(t, err)
	assert.Equal(t, 4, lns4)

	err = lm.Flush(lns4)
	assert.NoError(t, err)
	blkNum, err := lm.fm.BlockNum(testLogFile)
	assert.NoError(t, err)
	assert.Equal(t, 2, blkNum)

}
