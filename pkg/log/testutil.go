package log

import (
	"testing"

	"github.com/fchimpan/simpleDB/pkg/file"
	"github.com/stretchr/testify/assert"
)

const (
	testLogFile = "test_log.db"
	blockSize   = 64
)

func setupFileMgr(t *testing.T, dir string) *file.FileMgr {
	t.Helper()

	fm, err := file.NewFileMgr(dir, blockSize)
	assert.NoError(t, err)
	assert.True(t, !fm.IsNew())
	return fm
}

func setupLogMgr(t *testing.T, fm *file.FileMgr) *LogMgr {
	t.Helper()
	lm, err := NewLogMgr(fm, testLogFile)
	assert.NoError(t, err)
	lm.logPage.SetInt(0, blockSize) // ヘッダーの境界を設定
	return lm
}
