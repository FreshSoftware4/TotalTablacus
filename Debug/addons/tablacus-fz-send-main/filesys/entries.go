package filesys

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

func getDepth(path string) int {
	return len(strings.Split(path, string(os.PathSeparator)))
}

func PadCount(i int, total int) string {
	w := fmt.Sprint(len(fmt.Sprint(total)))
	return fmt.Sprintf("(%"+w+"d/%"+w+"d)", i, total)
}

type Entries struct {
	entries []Entry
}

func (es *Entries) Register(path string) {
	ent := Entry{Path: path}
	es.entries = append(es.entries, ent)
}

func (es *Entries) RegisterMulti(paths []string) {
	for _, p := range paths {
		es.Register(p)
	}
}

func (es Entries) UnMovable(dest string) (paths []string) {
	for _, ent := range es.entries {
		if ent.ExistsOn(dest) {
			paths = append(paths, ent.Path)
		}
	}
	return
}

func (es Entries) Size() int {
	return len(es.entries)
}

func (es *Entries) Exclude(path string) {
	var ents []Entry
	for _, ent := range es.entries {
		if ent.Path != path {
			ents = append(ents, ent)
		}
	}
	es.entries = ents
}

func (es Entries) sorted() []Entry {
	ss := es.entries
	sort.Slice(ss, func(i, j int) bool {
		return filepath.Base(ss[i].Path) < filepath.Base(ss[j].Path)
	})
	sort.SliceStable(ss, func(i, j int) bool {
		return ss[i].Path < ss[j].Path
	})
	sort.SliceStable(ss, func(i, j int) bool {
		return getDepth(ss[i].Path) > getDepth(ss[j].Path)
	})
	return ss
}

func (es Entries) Copy(src string, dest string) error {
	for i, ent := range es.sorted() {
		de := Entry{Path: dest}
		fmt.Printf("- %s Coping to %s <== %s%s\n", PadCount(i+1, len(es.entries)), de.DecoName(), ent.DecoRelPath(src), ent.DecoName())
		if err := ent.CopyTo(dest); err != nil {
			return err
		}
	}
	return nil
}

func (es Entries) Remove(from string) error {
	c := color.New(color.CrossedOut)
	for i, ent := range es.sorted() {
		fmt.Printf("- %s Deleting: %s%s\n", PadCount(i+1, len(es.entries)), ent.DecoRelPath(from), c.Sprint(ent.Name()))
		if err := ent.Remove(); err != nil {
			return err
		}
	}
	return nil
}
