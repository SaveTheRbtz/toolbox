package es_resource

import (
	rice "github.com/GeertJohan/go.rice"
	"testing"
)

func TestNewResource(t *testing.T) {
	r := NewResource(rice.MustFindBox("."))
	if _, err := r.Bytes("no existent"); err == nil {
		t.Error(err)
	}
	hfs := r.HttpFileSystem()
	if _, err := hfs.Open("no existent"); err == nil {
		t.Error(err)
	}
}

func TestNewSecureResource(t *testing.T) {
	r := NewSecureResource(rice.MustFindBox("."))
	if _, err := r.Bytes("no existent"); err == nil {
		t.Error(err)
	}
	hfs := r.HttpFileSystem()
	if _, err := hfs.Open("no existent"); err == nil {
		t.Error(err)
	}
}

func TestNewEmptyResource(t *testing.T) {
	r := EmptyResource()
	if _, err := r.Bytes("no existent"); err == nil {
		t.Error(err)
	}
	hfs := r.HttpFileSystem()
	if _, err := hfs.Open("no existent"); err == nil {
		t.Error(err)
	}
}

func TestEmptyBundle(t *testing.T) {
	b := EmptyBundle()

	if x, err := b.Templates().Bytes("no existent"); err == nil {
		t.Error(x, err)
	}
	if x, err := b.Messages().Bytes("no existent"); err == nil {
		t.Error(x, err)
	}
	if x, err := b.Web().Bytes("no existent"); err == nil {
		t.Error(x, err)
	}
	if x, err := b.Keys().Bytes("no existent"); err == nil {
		t.Error(x, err)
	}
	if x, err := b.Images().Bytes("no existent"); err == nil {
		t.Error(x, err)
	}
	if x, err := b.Data().Bytes("no existent"); err == nil {
		t.Error(x, err)
	}
}
