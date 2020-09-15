package file

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteBefore(t *testing.T) {
	dir, err := ioutil.TempDir("", "genembed")
	if err != nil {
		t.Error(err)
	}

	file := filepath.Join(dir, "f1")
	defer os.RemoveAll(dir)

	t.Run("replaceLarge", func(t *testing.T) {
		prepareFile(t, file, "2")
		err := writeBefore(t, file, "2", "----")
		assertNoError(t, err)
		equalFileContent(t, file, "----2")
	})
	t.Run("replaceSmall", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeBefore(t, file, "b", "-")
		assertNoError(t, err)
		equalFileContent(t, file, "a-bc")
	})
	t.Run("emptySrc", func(t *testing.T) {
		prepareFile(t, file, "")
		err := writeBefore(t, file, "b", "-")
		requireEqualError(t, err, ErrNotFoundPattern)
		equalFileContent(t, file, "")
	})
	t.Run("emptyPattern", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeBefore(t, file, "", "-")
		requireEqualError(t, err, ErrEmptyPattern)
		equalFileContent(t, file, "abc")
	})
	t.Run("emptyReplace", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeBefore(t, file, "b", "")
		assertNoError(t, err)
		equalFileContent(t, file, "abc")
	})
	t.Run("emptyReplaceBoth", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeBefore(t, file, "", "")
		requireEqualError(t, err, ErrEmptyPattern)
		equalFileContent(t, file, "abc")
	})
	t.Run("allEmpty", func(t *testing.T) {
		prepareFile(t, file, "")
		err := writeBefore(t, file, "", "")
		requireEqualError(t, err, ErrEmptyPattern)
		equalFileContent(t, file, "")
	})
	t.Run("patternMoreSrc", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeBefore(t, file, "bbbbbbb", "c")
		requireEqualError(t, err, ErrNotFoundPattern)
		equalFileContent(t, file, "abc")
	})

	t.Run("wultipleWrite", func(t *testing.T) {
		prepareFile(t, file, "abc")
		f, err := OpenFile(file)
		if err != nil {
			t.Error(err, "failed open file")
		}
		defer f.Close()

		err = f.WriteBefore([]byte("b"), []byte("-"))
		assertNoError(t, err)
		err = f.WriteBefore([]byte("b"), []byte("-"))
		assertNoError(t, err)
		err = f.WriteBefore([]byte("a"), []byte("-"))
		assertNoError(t, err)
		err = f.WriteBefore([]byte("c"), []byte("-"))
		assertNoError(t, err)
		err = f.WriteBefore([]byte("b"), []byte("-"))
		assertNoError(t, err)

		equalFileContent(t, file, "-a---b-c")
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
		assertNoError(t, err)
		equalFileContent(t, file, "2----")
	})
	t.Run("replaceSmall", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeAfter(t, file, "b", "-")
		assertNoError(t, err)
		equalFileContent(t, file, "ab-c")
	})
	t.Run("emptySrc", func(t *testing.T) {
		prepareFile(t, file, "")
		err := writeAfter(t, file, "b", "-")
		requireEqualError(t, err, ErrNotFoundPattern)
		equalFileContent(t, file, "")
	})
	t.Run("emptyPattern", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeAfter(t, file, "", "-")
		requireEqualError(t, err, ErrEmptyPattern)
		equalFileContent(t, file, "abc")
	})
	t.Run("emptyReplace", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeAfter(t, file, "b", "")
		assertNoError(t, err)
		equalFileContent(t, file, "abc")
	})
	t.Run("emptyReplaceBoth", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeAfter(t, file, "", "")
		requireEqualError(t, err, ErrEmptyPattern)
		equalFileContent(t, file, "abc")
	})
	t.Run("allEmpty", func(t *testing.T) {
		prepareFile(t, file, "")
		err := writeAfter(t, file, "", "")
		requireEqualError(t, err, ErrEmptyPattern)
		equalFileContent(t, file, "")
	})
	t.Run("patternMoreSrc", func(t *testing.T) {
		prepareFile(t, file, "abc")
		err := writeAfter(t, file, "bbbbbbb", "c")
		requireEqualError(t, err, ErrNotFoundPattern)
		equalFileContent(t, file, "abc")
	})

	t.Run("wultipleWrite", func(t *testing.T) {
		prepareFile(t, file, "abc")
		f, err := OpenFile(file)
		if err != nil {
			t.Error(err, "failed open file")
		}
		defer f.Close()

		err = f.WriteAfter([]byte("b"), []byte("-"))
		assertNoError(t, err)
		err = f.WriteAfter([]byte("b"), []byte("-"))
		assertNoError(t, err)
		err = f.WriteAfter([]byte("a"), []byte("-"))
		assertNoError(t, err)
		err = f.WriteAfter([]byte("c"), []byte("-"))
		assertNoError(t, err)
		err = f.WriteAfter([]byte("b"), []byte("-"))
		assertNoError(t, err)

		equalFileContent(t, file, "a-b---c-")
	})
}

func prepareFile(t *testing.T, filename string, dat string) {
	t.Helper()
	err := ioutil.WriteFile(filename, []byte(dat), 0666)
	if err != nil {
		t.Error(err)
	}
}

func writeAfter(t *testing.T, filename string, pattern, dat string) error {
	t.Helper()
	f, err := OpenFile(filename)
	if err != nil {
		t.Error(err, "failed open file")
	}

	defer f.Close()
	return f.WriteAfter([]byte(pattern), []byte(dat))
}

func writeBefore(t *testing.T, filename string, pattern, dat string) error {
	t.Helper()
	f, err := OpenFile(filename)
	if err != nil {
		t.Error(err, "failed open file")
	}

	defer f.Close()
	return f.WriteBefore([]byte(pattern), []byte(dat))
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("got error %v, expected without error", err)
	}
}

func requireEqualError(t *testing.T, got error, want error) {
	if got == nil || want == nil {
		t.Error("invalid arguments: must be errors")
	}
	if got.Error() != want.Error() {
		t.Error("not equal errors")
	}
}

func equalFileContent(t *testing.T, filename string, want string) {
	t.Helper()
	got, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(err, "failed to get data from file")
	}
	if !bytes.Equal([]byte(want), got) {
		t.Errorf("not equal contents: got=%q, want=%q", got, want)
	}
}
