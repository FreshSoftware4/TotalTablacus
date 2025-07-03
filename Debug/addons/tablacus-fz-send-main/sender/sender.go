package sender

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AWtnb/go-asker"
	"github.com/AWtnb/tablacus-fz-send/dir"
	"github.com/AWtnb/tablacus-fz-send/filesys"
)

var (
	ErrNoSubDir    = errors.New("no subdir to move")
	ErrInvalidDest = errors.New("invalid dest path")
)

type Sender struct {
	Src  string
	Dest string
}

func (sdr Sender) isDisposal() bool {
	return sdr.Dest == filesys.TrashName
}

func (sdr Sender) targets() ([]string, error) {
	var d dir.Dir
	if sdr.isDisposal() {
		d.Init(sdr.Src, 1, true, false)
	} else {
		d.Init(sdr.Src, -1, true, false)
	}
	if fs, err := os.Stat(sdr.Dest); err == nil && fs.IsDir() {
		d.Except(sdr.Dest)
	}
	return d.SelectItems()
}

func (sdr Sender) destPath() (string, error) {
	if fs, err := os.Stat(sdr.Dest); err == nil && fs.IsDir() {
		var dd dir.Dir
		dd.Init(sdr.Dest, 1, false, true)
		if 1 < len(dd.Member()) {
			return dd.SelectItem()
		}
		return sdr.Dest, nil
	}
	if strings.Contains(sdr.Dest, string(os.PathSeparator)) {
		return "", ErrInvalidDest
	}
	if sdr.isDisposal() {
		return dir.Create(sdr.Src, sdr.Dest)
	}
	if len(sdr.Dest) < 1 {
		var dd dir.Dir
		dd.Init(sdr.Src, 1, false, false)
		sds := dd.Member()
		if len(sds) < 1 {
			return "", ErrNoSubDir
		}
		if len(sds) == 1 {
			return sds[0], nil
		}
		return dd.SelectItem()
	}
	return dir.Create(sdr.Src, sdr.Dest)
}

func (sdr Sender) sendItems(paths []string, dest string) error {
	var fes filesys.Entries
	fes.RegisterMulti(paths)
	dupls := fes.UnMovable(dest)
	for _, dp := range dupls {
		a := asker.Asker{Accept: "y", Reject: "n"}
		e := filesys.Entry{Path: dp}
		d := filesys.Entry{Path: dest}
		a.Ask(fmt.Sprintf("Name duplicated: %s in %s\nOverwrite?", e.DecoName(), d.DecoName()))
		if !a.Accepted() {
			fmt.Println("==> Skipped")
			fes.Exclude(dp)
		}
	}
	if fes.Size() < 1 {
		return nil
	}
	if err := fes.Copy(sdr.Src, dest); err != nil {
		return err
	}

	if sdr.isDisposal() {
		return fes.Remove(sdr.Src)
	}

	a := asker.Asker{Accept: "y", Reject: "n"}
	a.Ask("Delete original?")
	if a.Accepted() {
		if err := fes.Remove(sdr.Src); err != nil {
			return err
		}
	}
	return nil
}

func (sdr Sender) Send() error {
	ts, err := sdr.targets()
	if err != nil {
		return err
	}

	d, err := sdr.destPath()
	if err != nil {
		return err
	}

	if err := sdr.sendItems(ts, d); err != nil {
		return err
	}

	return nil
}
