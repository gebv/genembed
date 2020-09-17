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
// Empowered the file by adding the methods to insert data before or after the pattern.
type File struct {
	*os.File
}

// WriteBefore writes data before the pattern if it was found in the file.
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
	pos, err := lastIndex(f.File, buf, pattern)

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
	pos, err := lastIndex(f.File, buf, pattern)

	if err != nil {
		return err
	}

	_, err = f.WriteAt(append(dat, buf[len(pattern):size-pos]...), pos+int64(len(pattern)))
	return err
}

var (
	// ErrNotFoundPattern is returned when not found pattern (after or before which should be an insert) in file.
	ErrNotFoundPattern = errors.New("not found pattern")

	// ErrEmptyPattern is returned when invalid pattern.
	ErrEmptyPattern = errors.New("empty pattern")
)

// lastIndex returns the value of the start position of the last matched pattern.
// If not found matched pattern returns position value -1 and error 'not found pattern'.
//
// The buffer contains the tail after the found pattern. If not found, the buffer contains a copy of the reader.
// NOTE: seeker is dropped to os.SEEK_END. Searching from the end every time.
func lastIndex(f readAtSeeker, buf []byte, pattern []byte) (pos int64, err error) {
	if len(pattern) == 0 {
		return -1, ErrEmptyPattern
	}

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

type readAtSeeker interface {
	io.Seeker
	io.ReaderAt
}
