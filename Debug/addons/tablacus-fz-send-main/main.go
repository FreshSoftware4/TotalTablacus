package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AWtnb/tablacus-fz-send/dir"
	"github.com/AWtnb/tablacus-fz-send/sender"
	"github.com/fatih/color"
	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		src  string
		dest string
	)
	flag.StringVar(&src, "src", "", "location of items to copy or move")
	flag.StringVar(&dest, "dest", "", "destination to copy or move")
	flag.Parse()
	if len(src) < 1 {
		src = os.ExpandEnv(`C:\Users\${USERNAME}\Desktop`)
	}
	os.Exit(run(src, dest))
}

func warn(s string) {
	color.Yellow(fmt.Sprintf("WARNING: %s\n", s))
	fmt.Scanln()
}

func reportError(err error) {
	if err == fuzzyfinder.ErrAbort {
		return
	}
	color.Magenta(fmt.Sprintf("ERROR: %s\n", err.Error()))
	fmt.Scanln()
}

func run(src string, dest string) int {
	if src == dest {
		reportError(errors.New("src and dest path should be different"))
		return 1
	}
	if src == ".." {
		src = filepath.Dir(dest)
	}

	s := sender.Sender{Src: src, Dest: dest}
	err := s.Send()
	if err != nil {
		if err == sender.ErrNoSubDir || err == dir.ErrNoItem {
			warn(err.Error())
			return 0
		}
		reportError(err)
		return 1
	}

	dir.Show(src)
	fmt.Scanln()
	return 0
}
