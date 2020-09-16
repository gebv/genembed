package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

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
