package log

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogIterator_ReadsRecordsInReverseOrder(t *testing.T) {
	dir := t.TempDir()

	fm := setupFileMgr(t, dir)
	logMgr := setupLogMgr(t, fm)

	// ログに追記
	records := [][]byte{
		[]byte("first"),
		[]byte("second"),
		[]byte("third"),
	}

	var lsns []int
	for _, rec := range records {
		lsn, err := logMgr.Append(rec)
		require.NoError(t, err)
		lsns = append(lsns, lsn)
	}

	require.NoError(t, logMgr.Flush(lsns[len(lsns)-1]))

	iter, err := NewLogIterator(fm, logMgr.CurrentBlk())
	require.NoError(t, err)

	var readBack [][]byte
	for iter.HasNext() {
		rec, err := iter.Next()
		require.NoError(t, err)
		readBack = append(readBack, rec)
	}

	// ログは新->旧の順に返される
	require.Equal(t, []string{"third", "second", "first"}, toStrings(readBack))
}

func toStrings(recs [][]byte) []string {
	var result []string
	for _, r := range recs {
		result = append(result, string(r))
	}
	return result
}
