package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenFile(t *testing.T) {
	_, err := OpenFile("not/eixsts/path/to/file")
	require.Error(t, err)
	require.EqualError(t, err, "open not/eixsts/path/to/file: no such file or directory")
}

func TestWriteBefore(t *testing.T) {
	dir, err := ioutil.TempDir("", "genembed")
	require.NoError(t, err)

	file := filepath.Join(dir, "f1")
	defer os.RemoveAll(dir)

	t.Run("replaceLarge", func(t *testing.T) {
		prepareFile(t, file, "2")
		err := writeBefore(t, file, "2", "----")
		require.NoError(t, err)
		requireEqualFileContent(t, file, "----2")
	})
	t.Run("replaceSmall", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeBefore(t, file, "b", "-")
		require.NoError(t, err)
		requireEqualFileContent(t, file, "a-bc")
	})
	t.Run("emptySrc", func(t *testing.T) {
		prepareFile(t, file, "")
		err := writeBefore(t, file, "b", "-")
		require.Error(t, err)
		require.EqualError(t, err, ErrNotFoundPattern.Error())
		requireEqualFileContent(t, file, "")
	})
	t.Run("emptyPattern", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeBefore(t, file, "", "-")
		require.Error(t, err)
		require.EqualError(t, err, ErrEmptyPattern.Error())
		requireEqualFileContent(t, file, "abc")
	})
	t.Run("emptyReplace", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeBefore(t, file, "b", "")
		require.NoError(t, err)
		requireEqualFileContent(t, file, "abc")
	})
	t.Run("emptyReplaceBoth", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeBefore(t, file, "", "")
		require.Error(t, err)
		require.EqualError(t, err, ErrEmptyPattern.Error())
		requireEqualFileContent(t, file, "abc")
	})
	t.Run("allEmpty", func(t *testing.T) {
		prepareFile(t, file, "")
		err := writeBefore(t, file, "", "")
		require.Error(t, err)
		require.EqualError(t, err, ErrEmptyPattern.Error())
		requireEqualFileContent(t, file, "")
	})
	t.Run("patternMoreSrc", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeBefore(t, file, "bbbbbbb", "c")
		require.Error(t, err)
		require.EqualError(t, err, ErrNotFoundPattern.Error())
		requireEqualFileContent(t, file, "abc")
	})

	t.Run("wultipleWrite", func(t *testing.T) {
		prepareFile(t, file, "abc")
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
	dir, err := ioutil.TempDir("", "genembed")
	if err != nil {
		t.Error(err)
	}

	file := filepath.Join(dir, "f1")
	defer os.RemoveAll(dir)

	t.Run("replaceLarge", func(t *testing.T) {
		prepareFile(t, file, "2")
		err := writeAfter(t, file, "2", "----")
		require.NoError(t, err)
		requireEqualFileContent(t, file, "2----")
	})
	t.Run("replaceSmall", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeAfter(t, file, "b", "-")
		require.NoError(t, err)
		requireEqualFileContent(t, file, "ab-c")
	})
	t.Run("emptySrc", func(t *testing.T) {
		prepareFile(t, file, "")
		err := writeAfter(t, file, "b", "-")
		require.Error(t, err)
		require.EqualError(t, err, ErrNotFoundPattern.Error())
		requireEqualFileContent(t, file, "")
	})
	t.Run("emptyPattern", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeAfter(t, file, "", "-")
		require.Error(t, err)
		require.EqualError(t, err, ErrEmptyPattern.Error())
		requireEqualFileContent(t, file, "abc")
	})
	t.Run("emptyReplace", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeAfter(t, file, "b", "")
		require.NoError(t, err)
		requireEqualFileContent(t, file, "abc")
	})
	t.Run("emptyReplaceBoth", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeAfter(t, file, "", "")
		require.Error(t, err)
		require.EqualError(t, err, ErrEmptyPattern.Error())
		requireEqualFileContent(t, file, "abc")
	})
	t.Run("allEmpty", func(t *testing.T) {
		prepareFile(t, file, "")
		err := writeAfter(t, file, "", "")
		require.Error(t, err)
		require.EqualError(t, err, ErrEmptyPattern.Error())
		requireEqualFileContent(t, file, "")
	})
	t.Run("patternMoreSrc", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeAfter(t, file, "bbbbbbb", "c")
		require.Error(t, err)
		require.EqualError(t, err, ErrNotFoundPattern.Error())
		requireEqualFileContent(t, file, "abc")
	})

	t.Run("wultipleWrite", func(t *testing.T) {
		prepareFile(t, file, "abc")
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

func prepareFile(t *testing.T, filename string, dat string) {
	t.Helper()
	err := ioutil.WriteFile(filename, []byte(dat), 0666)
	require.NoError(t, err)
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
