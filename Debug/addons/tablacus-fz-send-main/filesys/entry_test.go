package filesys_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/AWtnb/go-testtree"
	"github.com/AWtnb/tablacus-fz-send/filesys"
)

func TestMakeTestTree(t *testing.T) {
	p := `C:\Personal\gotemp\piyo`
	err := makeTestTree(p)
	if err != nil {
		t.Error(err)
	}
	showTreeContent(t, p)
}

func showTreeContent(t *testing.T, root string) {
	d := filepath.Dir(root)
	for _, c := range testtree.GetChildItems(root) {
		rel := strings.TrimPrefix(c, d)
		t.Log(rel)
	}
}

func TestNameDir(t *testing.T) {
	p := `C:\Personal\gotemp`
	if err := testtree.MakeTestDir(p); err != nil {
		t.Error(err)
	}
	e := filesys.Entry{Path: p}
	t.Logf("name of '%s': '%s'", p, e.Name())
}

func TestNameFile(t *testing.T) {
	p := `C:\Personal\gotemp\hoge.txt`
	if err := testtree.MakeTestFile(p); err != nil {
		t.Error(err)
	}
	e := filesys.Entry{Path: p}
	t.Logf("name of '%s' is '%s'", p, e.Name())
}

func TestExistsOn(t *testing.T) {
	p1 := `C:\Personal\gotemp\fuga`
	if err := testtree.MakeTestDir(p1); err != nil {
		t.Error(err)
	}
	p2 := filepath.Join(p1, "hoge.txt")
	if err := testtree.MakeTestFile(p2); err != nil {
		t.Error(err)
	}
	e := filesys.Entry{Path: p2}
	if e.ExistsOn(p1) {
		t.Logf("'%s' exists in '%s'", p2, p1)
	} else {
		t.Logf("'%s' not exists in '%s'", p2, p1)
	}
}

func TestDirCopy(t *testing.T) {
	p := `C:\Personal\gotemp\fuga`
	if err := makeTestTree(p); err != nil {
		t.Error(err)
	}
	showTreeContent(t, p)
	e := filesys.Entry{Path: p}
	d := `C:\Personal\gotemp\to`
	if err := testtree.MakeTestDir(d); err != nil {
		t.Error(err)
	}
	if err := e.CopyTo(d); err != nil {
		t.Error(err)
	}
	showTreeContent(t, d)
}

func TestFileCopy(t *testing.T) {
	p := `C:\Personal\gotemp\fuga.txt`
	if err := testtree.MakeTestFile(p); err != nil {
		t.Error(err)
	}
	e := filesys.Entry{Path: p}
	d := `C:\Personal\gotemp\to`
	if err := testtree.MakeTestDir(d); err != nil {
		t.Error(err)
	}
	if err := e.CopyTo(d); err != nil {
		t.Error(err)
	}
}

func TestMember(t *testing.T) {
	p := `C:\Personal\gotemp\aaa`
	if err := testtree.MakeTestDir(p); err != nil {
		t.Error(err)
	}
	e := filesys.Entry{Path: p}
	for _, m := range e.Member() {
		t.Log(m.Name())
	}
}
