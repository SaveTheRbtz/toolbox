package qt_file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func MakeDummyFile(name string) (path string, err error) {
	d, err := ioutil.TempFile("", name)
	if err != nil {
		return "", err
	}
	_, err = d.Write([]byte(time.Now().Format(time.RFC3339)))
	if err != nil {
		os.Remove(d.Name())
		return "", err
	}
	d.Close()
	return d.Name(), nil
}

func TestWithTestFile(t *testing.T, name, content string, f func(path string)) {
	tf, err := MakeTestFile(name, content)
	if err != nil {
		t.Error(err)
		return
	}
	f(tf)
	_ = os.Remove(tf)
}

func MakeTestFile(name string, content string) (path string, err error) {
	d, err := ioutil.TempFile("", name)
	if err != nil {
		return "", err
	}
	_, err = d.Write([]byte(content))
	if err != nil {
		_ = os.Remove(d.Name())
		return "", err
	}
	_ = d.Close()
	return d.Name(), nil
}

func TestWithTestFolder(t *testing.T, name string, withContent bool, f func(path string)) {
	path, err := MakeTestFolder(name, withContent)
	if err != nil {
		t.Error(err)
		return
	}
	f(path)
	if err = os.RemoveAll(path); err != nil {
		t.Error(err)
	}
}

func MakeTestFolder(name string, withContent bool) (path string, err error) {
	path, err = ioutil.TempDir("", name)
	if err != nil {
		return "", err
	}
	if withContent {
		err := ioutil.WriteFile(filepath.Join(path, "test.dat"), []byte(time.Now().String()), 0644)
		if err != nil {
			os.RemoveAll(path)
			return "", err
		}
	}
	return path, nil
}
