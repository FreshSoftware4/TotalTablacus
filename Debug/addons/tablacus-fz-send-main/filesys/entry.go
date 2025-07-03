package filesys

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/go-dircopy"
	"github.com/fatih/color"
)

var TrashName = "_obsolete"

func copyFile(src string, newPath string) error {
	if src == newPath {
		return errors.New("src and dest are the same path")
	}
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	nf, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer nf.Close()
	if _, err = io.Copy(nf, sf); err != nil {
		return err
	}
	return nil
}

type Entry struct {
	Path string
}

func (e Entry) Name() string {
	return filepath.Base(e.Path)
}

func (e Entry) DecoRelPath(base string) string {
	b := filepath.Dir(base)
	rel, _ := filepath.Rel(b, e.Path)
	c := strings.TrimSuffix(rel, e.Name())
	return color.HiBlackString(filepath.ToSlash(c))
}

func (e Entry) DecoName() string {
	b := filepath.Base(e.Path)
	fs, err := os.Stat(e.Path)
	if err != nil {
		return fmt.Sprintf("%s (non-exists)", b)
	}
	if fs.IsDir() {
		c := color.New(color.BgYellow, color.FgBlack)
		return c.Sprint(b)
	}
	x := filepath.Ext(b)
	n := strings.TrimSuffix(b, x)
	return color.CyanString(n) + x
}

func (e Entry) isDir() bool {
	fi, err := os.Stat(e.Path)
	return err == nil && fi.IsDir()
}

func (e Entry) reborn(dest string) string {
	return filepath.Join(dest, filepath.Base(e.Path))
}

func (e Entry) ExistsOn(dest string) bool {
	p := e.reborn(dest)
	_, err := os.Stat(p)
	return err == nil
}

func (e Entry) CopyTo(dest string) error {
	fs, err := os.Stat(e.Path)
	if err != nil {
		return err
	}

	newPath := e.reborn(dest)
	if fs.IsDir() {
		return dircopy.Copy(e.Path, newPath)
	}

	return copyFile(e.Path, newPath)
}

func (e Entry) Remove() error {
	if e.isDir() {
		return os.RemoveAll(e.Path)
	}
	return os.Remove(e.Path)
}

func (e Entry) Member() (entries []Entry) {
	if !e.isDir() {
		return
	}
	fs, err := os.ReadDir(e.Path)
	if err != nil {
		return
	}
	for _, f := range fs {
		p := filepath.Join(e.Path, f.Name())
		ent := Entry{Path: p}
		entries = append(entries, ent)
	}
	return
}
