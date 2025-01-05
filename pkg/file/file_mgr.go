package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileMgr struct {
	dbDirectory string
	blockSize   int
	isNew       bool
	openFiles   map[string]*os.File
	mu          sync.Mutex
}

func NewFileMgr(dbDirectory string, blockSize int) (*FileMgr, error) {
	isNew := false
	info, err := os.Stat(dbDirectory)
	if os.IsNotExist(err) {
		isNew = true
		if err := os.MkdirAll(dbDirectory, 0755); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		if !info.IsDir() {
			return nil, fmt.Errorf("%w is not a directory", err)
		}
	}

	return &FileMgr{
		dbDirectory: dbDirectory,
		blockSize:   blockSize,
		isNew:       isNew,
		openFiles:   make(map[string]*os.File),
	}, nil

}

// The read method seeks to the  appropriate position in the specified file and reads the contents of that block to the  byte buffer of the specified page.
func (fm *FileMgr) Read(blk *BlockID, p *Page) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	f, err := fm.getFile(blk.FileName())
	if err != nil {
		return fmt.Errorf("failed to getFile: %w", err)
	}
	if _, err := f.Seek(int64(blk.Number()*fm.blockSize), 0); err != nil {
		return err
	}

	if _, err := f.Read(p.Contents().Bytes()); err != nil {
		return err
	}
	return nil
}

func (fm *FileMgr) Write(blk *BlockID, p *Page) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	f, err := fm.getFile(blk.FileName())
	if err != nil {
		return fmt.Errorf("failed to getFile: %w", err)
	}

	if _, err := f.Seek(int64(blk.Number()*fm.blockSize), 0); err != nil {
		return err
	}

	if _, err := f.Write(p.Contents().Bytes()); err != nil {
		return err
	}

	return nil
}

// The append  method seeks to the end of the file and writes an empty array of bytes to it, which causes the OS to automatically extend the file.
func (fm *FileMgr) Append(fileName string) (*BlockID, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	size, err := fm.blockNum(fileName)
	if err != nil {
		return &BlockID{}, err
	}

	blkID := NewBlockID(fileName, size)

	f, err := fm.getFile(blkID.fileName)
	if err != nil {
		return &BlockID{}, err
	}

	if _, err := f.Seek(int64(blkID.Number()*fm.blockSize), 0); err != nil {
		return &BlockID{}, err
	}

	if _, err := f.Write(make([]byte, fm.blockSize)); err != nil {
		return &BlockID{}, err
	}

	return blkID, nil
}

func (fm *FileMgr) getFile(fileName string) (*os.File, error) {

	path := filepath.Join(fm.dbDirectory, fileName)
	if f, ok := fm.openFiles[path]; ok {
		return f, nil
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	fm.openFiles[path] = file
	return file, nil

}

func (fm *FileMgr) blockNum(fileName string) (int, error) {
	f, err := fm.getFile(fileName)
	if err != nil {
		return 0, err
	}

	info, err := f.Stat()
	if err != nil {
		return 0, err
	}

	return int(info.Size()) / fm.blockSize, nil
}

func (fm *FileMgr) IsNew() bool {
	return fm.isNew
}

func (fm *FileMgr) BlockSize() int {
	return fm.blockSize
}
