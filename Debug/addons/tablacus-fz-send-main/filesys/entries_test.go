package filesys_test

import (
	"path/filepath"
	"testing"

	"github.com/AWtnb/go-testtree"
	"github.com/AWtnb/tablacus-fz-send/filesys"
)

func getChildItems(root string) (paths []string) {
	e := filesys.Entry{Path: root}
	for _, m := range e.Member() {
		paths = append(paths, m.Path)
	}
	return
}

func makeTestTree(root string) error {
	ds := []string{"aa/bb", "aa/cc", "bb/ee"}
	fs := []string{"aa/bb/cc.txt", "aa/ff.txt", "dd.txt"}
	return testtree.MakeTestTree(root, ds, fs)
}

func TestRegister(t *testing.T) {
	p := `C:\Personal\gotemp\aaa`
	if err := makeTestTree(p); err != nil {
		t.Error(err)
	}
	ps := getChildItems(p)
	var es filesys.Entries
	es.RegisterMulti(ps)
	t.Log(es.Size())
}

func TestUnMovable(t *testing.T) {
	p1 := `C:\Personal\gotemp\aaa`
	if err := makeTestTree(p1); err != nil {
		t.Error(err)
	}
	p2 := `C:\Personal\gotemp\bbb`
	if err := makeTestTree(p2); err != nil {
		t.Error(err)
	}
	ps1 := getChildItems(p1)
	var es filesys.Entries
	es.RegisterMulti(ps1)
	uns := es.UnMovable(p2)
	for _, u := range uns {
		t.Logf("unmovable '%s' to '%s' because '%s' already exists", u, p2, filepath.Join(p2, filepath.Base(u)))
	}
}
