package log

import (
	"io"
	"os"
	"testing"
	"github.com/stretchr/testify/require"
	api "Proglog/api/v1"
)

var (
	record_val_1 = []byte("Hello world 1")
	record_val_2 = []byte("Hello world 2")
)

func TestSegment(t *testing.T) {
	c := Config{}

	c.Segment.MaxIndexBytes = entWidth * 2
	c.Segment.MaxStoreBytes = 1024

	segment, err := newSegment(os.TempDir(), 0, c)
	defer segment.Remove()

	require.NoError(t, err)

	testSegmentAppend(t, segment)
	testSegmentRead(t, segment)
	testSegmentClose(t, segment)
	testSegmentRemove(t, segment)
}

func testSegmentAppend(t *testing.T, s *segment) {
	offset, err := s.Append(&api.Record{Value: record_val_1})
	require.NoError(t, err)
	require.EqualValues(t, 0, offset)

	offset, err = s.Append(&api.Record{Value: record_val_2})
	require.NoError(t, err)
	require.EqualValues(t, 1, offset)

	require.Equal(t, s.IsMaxed(), true)

	_, err = s.Append(&api.Record{Value: record_val_2})
	require.Equal(t, io.EOF, err)
}

func testSegmentRead(t *testing.T, s *segment) {
	record, err := s.Read(0)
	require.NoError(t, err)
	require.EqualValues(t, record_val_1, record.Value)

	record, err = s.Read(1)
	require.NoError(t, err)
	require.EqualValues(t, record_val_2, record.Value)

	_, err = s.Read(2)
	require.Equal(t, io.EOF, err)
}

func testSegmentRemove(t *testing.T, s *segment) {
	err := s.Close()
	require.NoError(t, err)
}

func testSegmentClose(t *testing.T, s *segment) {
	err := s.Close()
	require.NoError(t, err)
}