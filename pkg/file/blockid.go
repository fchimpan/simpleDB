package file

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

// A BlockId object identifies a specific block by its file name and logical block number.
type BlockID struct {
	fileName string
	blkNum   int
}

func NewBlockID(fileName string, blkNum int) *BlockID {
	return &BlockID{
		fileName: fileName,
		blkNum:   blkNum,
	}
}

func (b *BlockID) FileName() string {
	return b.fileName
}

func (b *BlockID) Number() int {
	return b.blkNum
}

func (b *BlockID) Equals(other *BlockID) bool {
	return b.fileName == other.fileName && b.blkNum == other.blkNum
}

func (b *BlockID) String() string {
	return b.fileName + ":" + strconv.Itoa(b.blkNum)
}

func (b *BlockID) HashCode() string {
	h := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(h[:])
}
