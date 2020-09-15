package file

import (
	"bytes"
	"errors"
	"io"
	"os"
	"syscall"
)

// OpenFile opens and returns a file instance.
func OpenFile(filename string) (*File, error) {
	f, err := os.OpenFile(filename, syscall.O_CREAT|syscall.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return &File{f}, nil
}

// File instance of file.
type File struct {
	*os.File
}

func (f File) WriteBefore(pattern, dat []byte) (err error) {
	if len(pattern) == 0 {
		return ErrEmptyPattern
	}
	fstat, err := f.Stat()
	if err != nil {
		return err
	}
	size := fstat.Size()

	buf := make([]byte, size)
	pos, err := f.lastIndex(buf, pattern)

	if err != nil {
		return err
	}

	_, err = f.WriteAt(append(dat, buf[:size-pos]...), pos)
	return err
}

// WriteAfter writes data after the pattern if it was found in the file.
//
// NOTE: move starts from the end.
func (f File) WriteAfter(pattern, dat []byte) (err error) {
	if len(pattern) == 0 {
		return ErrEmptyPattern
	}
	fstat, err := f.Stat()
	if err != nil {
		return err
	}
	size := fstat.Size()

	buf := make([]byte, size)
	pos, err := f.lastIndex(buf, pattern)

	if err != nil {
		return err
	}

	_, err = f.WriteAt(append(dat, buf[len(pattern):size-pos]...), pos+int64(len(pattern)))
	return err
}

// lastIndex returns the value of the start position of the last matched pattern.
// If not found matched pattern returns position value -1 and error 'not found pattern'.
func (f File) lastIndex(buf []byte, pattern []byte) (pos int64, err error) {
	var seek int64 = 0
	var found bool

	for err == nil || pos >= 0 {

		if (int64(len(pattern)) + seek) > int64(cap(buf)) {
			break
		}

		pos, err = f.Seek((int64(len(pattern))+seek)*(-1), os.SEEK_END)
		if err != nil {
			return -1, err
		}

		_, err := f.ReadAt(buf, pos)
		if err != nil && err != io.EOF {
			return -1, err
		}

		found = bytes.Equal(pattern, buf[:len(pattern)])
		if found {
			break
		}

		seek++
	}

	if !found {
		return -1, ErrNotFoundPattern
	}

	return pos, nil
}

var (
	ErrNotFoundPattern = errors.New("not found pattern")
	ErrFailedSeeking   = errors.New("failed seeking")
	ErrEmptyPattern    = errors.New("empty pattern")
)

// type file interface {
// 	io.Seeker
// 	io.Closer
// 	io.ReaderAt
// 	io.WriterAt
// 	Sync() error
// 	Close() error
// }
