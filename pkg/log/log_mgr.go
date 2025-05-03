package log

import (
	"sync"

	"github.com/fchimpan/simpleDB/pkg/file"
)

type LogMgr struct {
	fm           *file.FileMgr
	logFileName  string
	logPage      *file.Page
	currentBlk   *file.BlockID
	latestLNS    int // page に書き込んだ最新のLSN
	lastSavedLNS int // 最後にディスクに書き込んだLSN
	mu           sync.Mutex
}

func (lm *LogMgr) CurrentBlk() *file.BlockID {
	return lm.currentBlk
}

func NewLogMgr(fm *file.FileMgr, logFileName string) (*LogMgr, error) {
	p := file.NewPage(fm.BlockSize())
	logFileSize, err := fm.BlockNum(logFileName)
	if err != nil {
		return nil, err
	}

	var currentBlk *file.BlockID
	if logFileSize == 0 {
		currentBlk, err = fm.Append(logFileName)
		if err != nil {
			return nil, err
		}
	} else {
		currentBlk = file.NewBlockID(logFileName, logFileSize-1)
		if err := fm.Read(currentBlk, p); err != nil {
			return nil, err
		}
	}
	return &LogMgr{
		fm:          fm,
		logFileName: logFileName,
		logPage:     p,
		currentBlk:  currentBlk,
	}, nil
}

func (lm *LogMgr) Append(rec []byte) (int, error) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	recSize := len(rec)
	boundary := lm.logPage.GetInt(0)
	bytesNeeded := recSize + file.HeaderSize // 先頭4バイトはヘッダーのサイズ

	if boundary-bytesNeeded < file.HeaderSize {
		if err := lm.Flush(lm.latestLNS); err != nil {
			return 0, err
		}
		blk, err := lm.fm.Append(lm.logFileName)
		if err != nil {
			return 0, err
		}
		lm.logPage.SetInt(0, lm.fm.BlockSize())
		if err := lm.fm.Write(blk, lm.logPage); err != nil {
			return 0, err
		}
		lm.currentBlk = blk
		boundary = lm.logPage.GetInt(0)
	}
	newBoundary := boundary - bytesNeeded

	lm.logPage.SetBytes(newBoundary, rec)
	lm.logPage.SetInt(0, newBoundary)
	lm.latestLNS++
	return lm.latestLNS, nil
}

// Flush は、指定された LSN までのログをディスクに書き込み
// logPage を空にするわけではない
func (lm *LogMgr) Flush(lsn int) error {
	if lsn >= lm.lastSavedLNS {
		if err := lm.fm.Write(lm.currentBlk, lm.logPage); err != nil {
			return err
		}
		lm.lastSavedLNS = lm.latestLNS
	}
	return nil
}
