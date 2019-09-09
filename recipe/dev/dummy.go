package dev

import (
	"bufio"
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"errors"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/recpie/app_kitchen"
	"github.com/watermint/toolbox/infra/recpie/app_vo"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type DummyVO struct {
	Path     string
	Dest     string
	MaxEntry int
}

type DummyEntry struct {
	Tag         string `json:".tag"`
	PathDisplay string `json:"path_display"`
}

type Dummy struct {
}

func (z *Dummy) Test(c app_control.Control) error {
	return nil
}

func (z *Dummy) Hidden() {
}

func (z *Dummy) Requirement() app_vo.ValueObject {
	return &DummyVO{}
}

func (z *Dummy) Exec(k app_kitchen.Kitchen) error {
	var vo interface{} = k.Value()
	dvo := vo.(*DummyVO)
	l := k.Log().With(zap.String("path", dvo.Path))

	f, err := os.Open(dvo.Path)
	if err != nil {
		l.Error("Unable to open file", zap.Error(err))
		return err
	}
	defer f.Close()
	br := bufio.NewReader(f)
	entries := 0

	for {
		line, _, err := br.ReadLine()
		switch {
		case err != nil && err == io.EOF:
			l.Info("Done")
			return nil

		case err != nil:
			l.Error("Unable to read", zap.Error(err))
			return err
		}

		de := &DummyEntry{}
		if err = json.Unmarshal(line, de); err != nil {
			l.Error("Unable to unmarshal", zap.Error(err))
			return err
		}

		if err = z.create(k, dvo.Dest, de); err != nil {
			return err
		}

		entries++
		if dvo.MaxEntry != 0 && entries >= dvo.MaxEntry {
			l.Info("Suspend", zap.Int("entries", entries))
			return nil
		}
	}
}

func (z *Dummy) create(k app_kitchen.Kitchen, base string, de *DummyEntry) error {
	l := k.Log()

	switch de.Tag {
	case "file":
		dir := z.anonPath(filepath.ToSlash(filepath.Dir(de.PathDisplay)))
		name := z.anonFileName(filepath.Base(de.PathDisplay))
		path := filepath.Join(dir, name)
		l.Debug("File", zap.String("file", path))
		pp := filepath.Join(base, dir)
		_, err := z.getOrCreate(pp)
		if err != nil {
			l.Debug("Folder create", zap.String("folder", path), zap.Error(err))
			return err
		}
		f, err := os.Create(filepath.Join(pp, name))
		if err != nil {
			l.Debug("Unable to create", zap.Error(err))
			return err
		}
		f.Close()

	case "folder":
		path := z.anonPath(de.PathDisplay)
		l.Debug("Folder", zap.String("folder", path))

		pp := filepath.Join(base, path)
		_, err := z.getOrCreate(pp)
		if err != nil {
			l.Debug("Folder create", zap.String("folder", path), zap.Error(err))
			return err
		}
	}

	return nil
}

func (z *Dummy) anonPath(path string) string {
	pp := strings.Split(path, "/")
	qq := make([]string, 0)
	for _, p := range pp {
		qq = append(qq, z.anonymize(p))
	}
	return filepath.Join(qq...)
}

func (z *Dummy) anonFileName(name string) string {
	ext := filepath.Ext(name)
	if 4 < len(ext) {
		ext = z.anonymize(ext)
	}
	return z.anonymize(name) + ext
}

func (z *Dummy) anonymize(name string) string {
	b := sha1.Sum([]byte(name))
	c := make([]byte, 20)
	copy(c[:], b[:])
	d := base32.StdEncoding.EncodeToString(c)
	l := utf8.RuneCountInString(name)

	if l < len(d) {
		return d[:l]
	} else {
		return d
	}
}

func (z *Dummy) getOrCreate(fqp string) (path string, err error) {
	st, err := os.Stat(fqp)
	switch {
	case err != nil && os.IsNotExist(err):
		err = os.MkdirAll(fqp, 0701)
		if err != nil {
			return "", err
		}
	case err != nil:
		return "", err

	case !st.IsDir():
		return "", errors.New("workspace path is not a directory")

	case st.Mode()&0700 == 0:
		return "", errors.New("no permission")
	}
	return fqp, nil
}