package dir

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AWtnb/go-walk"
	"github.com/AWtnb/tablacus-fz-send/filesys"
	"github.com/fatih/color"
	"github.com/ktr0731/go-fuzzyfinder"
)

var ErrNoItem = errors.New("no item to send")

func orderSelectable(ss []string) []string {
	sep := string(os.PathSeparator)
	sort.Slice(ss, func(i, j int) bool {
		return filepath.Base(ss[i]) > filepath.Base(ss[j])
	})
	sort.SliceStable(ss, func(i, j int) bool {
		return len(filepath.Base(ss[i])) > len(filepath.Base(ss[j]))
	})
	sort.SliceStable(ss, func(i, j int) bool {
		return len(strings.Split(ss[i], sep)) > len(strings.Split(ss[j], sep))
	})
	return ss
}

func getChildItem(root string, depth int, all bool, self bool) (paths []string) {
	var d walk.Dir
	d.Init(root, all, depth, filesys.TrashName)
	found, err := d.GetChildItem()
	if err != nil {
		return
	}
	for _, f := range found {
		n := filepath.Base(f)
		if (!self && f == root) || strings.HasSuffix(n, ".ini") || strings.HasPrefix(n, "~$") {
			continue
		}
		paths = append(paths, f)
	}
	return
}

func countDir(paths []string) int {
	n := 0
	for _, p := range paths {
		if fs, err := os.Stat(p); err == nil && fs.IsDir() {
			n += 1
		}
	}
	return n
}

func groupExt(paths []string) map[string]int {
	var total map[string]int = make(map[string]int)
	for _, p := range paths {
		if e := filepath.Ext(p); 0 < len(e) {
			total[e] += 1
		}
	}
	return total
}

func Show(path string) {
	pe := filesys.Entry{Path: path}
	dn := pe.DecoName()
	left := getChildItem(path, 1, true, false)
	if len(left) < 1 {
		fmt.Printf("(now %s is empty)\n", dn)
		return
	}
	if len(left) == 1 {
		fmt.Printf("Left item on %s:\n", dn)
		e := filesys.Entry{Path: left[0]}
		fmt.Printf(" - %s\n", e.DecoName())
		return
	}
	fmt.Printf("Left items on %s:\n", dn)
	fmt.Printf("- %s %s\n", color.YellowString("\U0001F4C1"), color.HiGreenString(fmt.Sprintf("\u00d7%d", countDir(left))))
	exs := groupExt(left)
	for k := range exs {
		fmt.Printf("- %s %s\n", k, color.HiGreenString(fmt.Sprintf("\u00d7%d", exs[k])))
	}
}

func getPerm(path string) fs.FileMode {
	s := string(os.PathSeparator)
	elems := strings.Split(path, s)
	for i := 0; i < len(elems); i++ {
		ln := len(elems) - i
		p := strings.Join(elems[0:ln], s)
		if fs, err := os.Stat(p); err == nil && fs.IsDir() {
			return fs.Mode() & os.ModePerm
		}
	}
	return 0700
}

func Create(dir string, name string) (string, error) {
	p := filepath.Join(dir, name)
	if fs, err := os.Stat(p); err == nil && fs.IsDir() {
		return p, nil
	}
	return p, os.MkdirAll(p, getPerm(p))
}

type Dir struct {
	path   string
	member []string
}

func (d *Dir) Init(path string, depth int, all bool, self bool) {
	d.path = path
	d.member = getChildItem(d.path, depth, all, self)
}

func (d *Dir) Except(path string) {
	paths := []string{}
	for _, p := range d.member {
		if strings.HasPrefix(p, path) {
			continue
		}
		paths = append(paths, p)
	}
	d.member = paths
}

func (d Dir) Member() []string {
	return d.member
}

func (d Dir) rel(path string) string {
	rel, err := filepath.Rel(d.path, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	rel = filepath.ToSlash(rel)
	if fs, err := os.Stat(path); err == nil && fs.IsDir() {
		return fmt.Sprintf("%s \U0001F4C1", rel)
	}
	return rel
}

func (d Dir) SelectItem() (path string, err error) {
	if len(d.member) < 1 {
		err = ErrNoItem
		return
	}
	ss := orderSelectable(d.member)
	idx, err := fuzzyfinder.Find(ss, func(i int) string {
		return d.rel(ss[i])
	}, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop))
	if err != nil {
		return
	}
	path = ss[idx]
	return
}

func (d Dir) SelectItems() (paths []string, err error) {
	if len(d.member) < 1 {
		err = ErrNoItem
		return
	}
	ss := orderSelectable(d.member)
	idxs, err := fuzzyfinder.FindMulti(ss, func(i int) string {
		return d.rel(ss[i])
	}, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop))
	if err != nil {
		return
	}
	for _, i := range idxs {
		paths = append(paths, ss[i])
	}
	return
}
