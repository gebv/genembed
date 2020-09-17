package file

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenFile(t *testing.T) {
	t.Run("pathNonExistentDir", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		_, err := OpenFile(file)
		require.NoError(t, err)
	})
	t.Run("pathNonExistentDir", func(t *testing.T) {
		_, err := OpenFile("not/eixsts/path/to/file")
		require.Error(t, err)
		require.EqualError(t, err, ErrNotExistsInputFile.Error())
	})

}

func TestWriteBefore(t *testing.T) {
	t.Run("replaceLarge", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "2")
		close()
		defer remove()

		err := writeBefore(t, file, "2", "----")
		require.NoError(t, err)
		requireEqualFileContent(t, file, "----2")
	})
	t.Run("replaceSmall", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		err := writeBefore(t, file, "b", "-")
		require.NoError(t, err)
		requireEqualFileContent(t, file, "a-bc")
	})
	t.Run("emptySrc", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "")
		close()
		defer remove()

		err := writeBefore(t, file, "b", "-")
		require.Error(t, err)
		require.EqualError(t, err, ErrNotFoundPattern.Error())
		requireEqualFileContent(t, file, "")
	})
	t.Run("emptyPattern", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		err := writeBefore(t, file, "", "-")
		require.Error(t, err)
		require.EqualError(t, err, ErrEmptyPattern.Error())
		requireEqualFileContent(t, file, "abc")
	})
	t.Run("emptyReplace", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		err := writeBefore(t, file, "b", "")
		require.NoError(t, err)
		requireEqualFileContent(t, file, "abc")
	})
	t.Run("emptyReplaceBoth", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		err := writeBefore(t, file, "", "")
		require.Error(t, err)
		require.EqualError(t, err, ErrEmptyPattern.Error())
		requireEqualFileContent(t, file, "abc")
	})
	t.Run("allEmpty", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "")
		close()
		defer remove()

		err := writeBefore(t, file, "", "")
		require.Error(t, err)
		require.EqualError(t, err, ErrEmptyPattern.Error())
		requireEqualFileContent(t, file, "")
	})
	t.Run("patternMoreSrc", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		err := writeBefore(t, file, "bbbbbbb", "c")
		require.Error(t, err)
		require.EqualError(t, err, ErrNotFoundPattern.Error())
		requireEqualFileContent(t, file, "abc")
	})

	t.Run("wultipleWrite", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		f, err := OpenFile(file)
		if err != nil {
			t.Error(err, "failed open file")
		}
		defer f.Close()

		err = f.WriteBefore([]byte("b"), []byte("-"))
		require.NoError(t, err)

		err = f.WriteBefore([]byte("b"), []byte("-"))
		require.NoError(t, err)

		err = f.WriteBefore([]byte("a"), []byte("-"))
		require.NoError(t, err)

		err = f.WriteBefore([]byte("c"), []byte("-"))
		require.NoError(t, err)

		err = f.WriteBefore([]byte("b"), []byte("-"))
		require.NoError(t, err)

		requireEqualFileContent(t, file, "-a---b-c")
	})
}

func TestWriteAfter(t *testing.T) {

	t.Run("replaceLarge", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "2")
		close()
		defer remove()

		err := writeAfter(t, file, "2", "----")
		require.NoError(t, err)
		requireEqualFileContent(t, file, "2----")
	})
	t.Run("replaceSmall", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		err := writeAfter(t, file, "b", "-")
		require.NoError(t, err)
		requireEqualFileContent(t, file, "ab-c")
	})
	t.Run("emptySrc", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "")
		close()
		defer remove()

		err := writeAfter(t, file, "b", "-")
		require.Error(t, err)
		require.EqualError(t, err, ErrNotFoundPattern.Error())
		requireEqualFileContent(t, file, "")
	})
	t.Run("emptyPattern", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		err := writeAfter(t, file, "", "-")
		require.Error(t, err)
		require.EqualError(t, err, ErrEmptyPattern.Error())
		requireEqualFileContent(t, file, "abc")
	})
	t.Run("emptyReplace", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		err := writeAfter(t, file, "b", "")
		require.NoError(t, err)
		requireEqualFileContent(t, file, "abc")
	})
	t.Run("emptyReplaceBoth", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		err := writeAfter(t, file, "", "")
		require.Error(t, err)
		require.EqualError(t, err, ErrEmptyPattern.Error())
		requireEqualFileContent(t, file, "abc")
	})
	t.Run("allEmpty", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "")
		close()
		defer remove()

		err := writeAfter(t, file, "", "")
		require.Error(t, err)
		require.EqualError(t, err, ErrEmptyPattern.Error())
		requireEqualFileContent(t, file, "")
	})
	t.Run("patternMoreSrc", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		err := writeAfter(t, file, "bbbbbbb", "c")
		require.Error(t, err)
		require.EqualError(t, err, ErrNotFoundPattern.Error())
		requireEqualFileContent(t, file, "abc")
	})

	t.Run("wultipleWrite", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "abc")
		close()
		defer remove()

		f, err := OpenFile(file)
		if err != nil {
			t.Error(err, "failed open file")
		}
		defer f.Close()

		err = f.WriteAfter([]byte("b"), []byte("-"))
		require.NoError(t, err)
		err = f.WriteAfter([]byte("b"), []byte("-"))
		require.NoError(t, err)
		err = f.WriteAfter([]byte("a"), []byte("-"))
		require.NoError(t, err)
		err = f.WriteAfter([]byte("c"), []byte("-"))
		require.NoError(t, err)
		err = f.WriteAfter([]byte("b"), []byte("-"))
		require.NoError(t, err)

		requireEqualFileContent(t, file, "a-b---c-")
	})
}

func tmpFileWith(t *testing.T, dat string) (filename string, closeFn func(), removeFn func()) {
	t.Helper()

	tmpFile, err := ioutil.TempFile("", "embeded")
	require.NoError(t, err)
	n, err := tmpFile.WriteString(dat)
	require.NoError(t, err)
	require.Len(t, dat, n)

	return tmpFile.Name(),
		func() {
			tmpFile.Close()
		}, func() {
			os.RemoveAll(tmpFile.Name())
		}
}

func writeAfter(t *testing.T, filename string, pattern, dat string) error {
	t.Helper()
	f, err := OpenFile(filename)
	require.NoError(t, err)

	defer f.Close()
	return f.WriteAfter([]byte(pattern), []byte(dat))
}

func writeBefore(t *testing.T, filename string, pattern, dat string) error {
	t.Helper()
	f, err := OpenFile(filename)
	require.NoError(t, err)

	defer f.Close()
	return f.WriteBefore([]byte(pattern), []byte(dat))
}

func requireEqualFileContent(t *testing.T, filename string, want string) {
	t.Helper()
	got, err := ioutil.ReadFile(filename)
	require.NoError(t, err)
	require.EqualValues(t, []byte(want), got)
}

func requireNotExistsFile(t *testing.T, filename string) {
	t.Helper()
	_, err := os.Stat(filename)
	require.True(t, os.IsNotExist(err))
}

func Test_lastIndex(t *testing.T) {
	tests := []struct {
		name    string
		f       readAtSeeker
		buf     []byte
		pattern []byte
		wantPos int64
		wantErr error
	}{
		{
			name:    "ok",
			f:       strings.NewReader("123"),
			buf:     make([]byte, 3),
			pattern: []byte("2"),
			wantPos: 1,
			wantErr: nil,
		},
		{
			name:    "ok",
			f:       strings.NewReader("123"),
			buf:     make([]byte, 3),
			pattern: []byte("1"),
			wantPos: 0,
			wantErr: nil,
		},
		{
			name:    "ok",
			f:       strings.NewReader("123"),
			buf:     make([]byte, 3),
			pattern: []byte("3"),
			wantPos: 2,
			wantErr: nil,
		},

		// small buffer
		{
			name:    "smallBuffer",
			f:       strings.NewReader("123"),
			buf:     make([]byte, 1),
			pattern: []byte("2"),
			wantPos: -1,
			wantErr: ErrNotFoundPattern,
		},
		{
			name:    "ok",
			f:       strings.NewReader("123"),
			buf:     make([]byte, 1),
			pattern: []byte("1"),
			wantPos: -1,
			wantErr: ErrNotFoundPattern,
		},
		{
			name:    "ok",
			f:       strings.NewReader("123"),
			buf:     make([]byte, 1),
			pattern: []byte("3"),
			wantPos: 2,
			wantErr: nil,
		},

		{
			name:    "notFound",
			f:       strings.NewReader("123"),
			buf:     make([]byte, 3),
			pattern: []byte("4"),
			wantPos: -1,
			wantErr: ErrNotFoundPattern,
		},

		{
			name:    "nilBuffer",
			f:       strings.NewReader("123"),
			buf:     nil,
			pattern: []byte("2"),
			wantPos: -1,
			wantErr: ErrNotFoundPattern,
		},

		{
			name:    "largerPattern", // but buffer size is correct)
			f:       strings.NewReader("123"),
			buf:     make([]byte, 3),
			pattern: []byte("2222222"),
			wantPos: -1,
			wantErr: ErrNotFoundPattern,
		},
		{
			name:    "nilPattern",
			f:       strings.NewReader("123"),
			buf:     make([]byte, 3),
			pattern: nil,
			wantPos: -1,
			wantErr: ErrEmptyPattern,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPos, err := lastIndex(tt.f, tt.buf, tt.pattern)
			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr.Error())
			}
			require.Equal(t, tt.wantPos, gotPos)
		})
	}
}

func Test_lastIndex_SpecialCases(t *testing.T) {
	t.Run("multipleCall", func(t *testing.T) {
		r := strings.NewReader("121212")
		gotPos, err := lastIndex(r, make([]byte, r.Size()), []byte("2"))
		require.NoError(t, err)
		require.EqualValues(t, 5, gotPos)

		gotPos, err = lastIndex(r, make([]byte, r.Size()), []byte("2"))
		require.NoError(t, err)
		require.EqualValues(t, 5, gotPos)
	})

	t.Run("bufferСontents", func(t *testing.T) {
		cases := []struct {
			in      string
			pattern string
			// buffer size is equal to the input
			// NOTE: buffer preparing before call method
			wantBuf []byte
			wantErr error
		}{
			{
				"1212",
				"2",
				[]byte{0x32, 0x0, 0x0, 0x0},
				nil,
			},
			{
				"1212",
				"1",
				[]byte{0x31, 0x32, 0x0, 0x0},
				nil,
			},

			{
				"1212",
				"4",
				[]byte{0x31, 0x32, 0x31, 0x32}, // because the buffer for tail
				ErrNotFoundPattern,
			},
			{
				"1212",
				"",
				[]byte{0x00, 0x00, 0x00, 0x00}, // it is ok, buffer preparing before call method
				ErrEmptyPattern,
			},
			{
				"",
				"4",
				[]byte{}, // because input is empty
				ErrNotFoundPattern,
			},
		}

		for _, tt := range cases {
			t.Run(tt.in, func(t *testing.T) {
				r := strings.NewReader(tt.in)
				buf := make([]byte, r.Size())
				_, err := lastIndex(r, buf, []byte(tt.pattern))
				if tt.wantErr == nil {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
					require.EqualError(t, err, tt.wantErr.Error())
				}
				require.EqualValues(t, tt.wantBuf, buf)
			})
		}
	})
}

func TestFile_Size(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "")
		close()
		defer remove()

		f, err := OpenFile(file)
		require.NoError(t, err)
		size, err := f.size()
		require.NoError(t, err)
		require.EqualValues(t, 0, size)
	})

	t.Run("ok", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "123")
		close()
		defer remove()

		f, err := OpenFile(file)
		require.NoError(t, err)
		size, err := f.size()
		require.NoError(t, err)
		require.EqualValues(t, 3, size)
	})

	t.Run("closed", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "123")
		close()
		defer remove()

		f, err := OpenFile(file)
		require.NoError(t, err)
		err = f.Close()
		require.NoError(t, err)

		size, err := f.size()
		require.Error(t, err)
		require.Contains(t, err.Error(), "use of closed file")
		require.EqualValues(t, -1, size)
	})
}

func Test_NegativeCases(t *testing.T) {
	t.Run("openNotExistsFile", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "123")
		close()
		remove()

		f, err := OpenFile(file)
		require.Nil(t, f)
		require.Error(t, err)
		require.EqualError(t, err, ErrNotExistsInputFile.Error())

		requireNotExistsFile(t, file)
	})

	t.Run("onFlyRemoveFile", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "123")
		close()
		defer remove()

		f, err := OpenFile(file)
		require.NoError(t, err)
		require.NotEmpty(t, f)

		remove()

		requireNotExistsFile(t, file)

		err = f.WriteAfter([]byte("2"), []byte("-"))
		require.NoError(t, err)

		requireNotExistsFile(t, file)
	})

	t.Run("nil", func(t *testing.T) {
		f := File{nil}
		err := f.WriteAfter([]byte("2"), []byte("-"))
		require.EqualError(t, err, ErrInvalid.Error())

		err = f.WriteBefore([]byte("2"), []byte("-"))
		require.EqualError(t, err, ErrInvalid.Error())

		_, err = f.size()
		require.EqualError(t, err, ErrInvalid.Error())
	})

	t.Run("perm", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "123")
		close()
		defer remove()

		err := os.Chmod(file, 0400)
		require.NoError(t, err)

		f, err := OpenFile(file)
		require.Nil(t, f)
		require.Error(t, err)
		require.Contains(t, err.Error(), "permission denied")

		requireEqualFileContent(t, file, "123")
	})

	t.Run("closed", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "123")
		close()
		defer remove()

		f, err := OpenFile(file)
		require.NotNil(t, f)
		require.NoError(t, err)

		err = f.Close()
		require.NoError(t, err)

		err = f.WriteAfter([]byte("2"), []byte("-"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "use of closed file")

		requireEqualFileContent(t, file, "123")
	})

	failCases := []struct {
		in        readAtSeeker
		buf       []byte
		pattern   string
		wantErr   bool
		conainErr string
	}{
		{
			strings.NewReader("123"),
			nil,
			"",
			true,
			ErrEmptyPattern.Error(),
		},

		{
			strings.NewReader("123"),
			nil,
			"2",
			true,
			ErrNotFoundPattern.Error(), // because buffer is small
		},
		{
			strings.NewReader("123"),
			make([]byte, 1),
			"2",
			true,
			ErrNotFoundPattern.Error(), // because buffer is small
		},

		{
			strings.NewReader("123"),
			make([]byte, 1), // invalid size
			"1234",          // pattern more than content and buffer
			true,
			ErrNotFoundPattern.Error(),
		},
		{
			strings.NewReader("123"),
			make([]byte, 3), // ok size
			"1234",          // pattern more than content and buffer
			true,
			ErrNotFoundPattern.Error(),
		},
		{
			strings.NewReader("123"),
			make([]byte, 4), // more than content but less than buf
			"12345",
			true,
			ErrNotFoundPattern.Error(),
		},
		{
			strings.NewReader("123"),
			make([]byte, 5), // same size as pattern
			"12345",
			true,
			"negative position", // seeker returns error
		},
		{
			strings.NewReader("123"),
			make([]byte, 8), // size is equal pattern + content
			"12345",
			true,
			"negative position", // seeker returns error
		},
		{
			strings.NewReader("123"),
			make([]byte, 10), // much larger size
			"12345",
			true,
			"negative position", // seeker returns error
		},
	}

	for _, tt := range failCases {
		t.Run(tt.pattern, func(t *testing.T) {
			_, err := lastIndex(tt.in, tt.buf, []byte(tt.pattern))
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.conainErr)
		})
	}

}

func Test_isExistsFile(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "123")
		close()
		defer remove()
		require.True(t, isExistsFile(file))
	})
	t.Run("notExists", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "123")
		close()
		remove()
		require.False(t, isExistsFile(file))
	})
	t.Run("perm", func(t *testing.T) {
		file, close, remove := tmpFileWith(t, "123")
		close()
		defer remove()

		err := os.Chmod(file, 0222)
		require.NoError(t, err)

		require.True(t, isExistsFile(file))
	})
}
