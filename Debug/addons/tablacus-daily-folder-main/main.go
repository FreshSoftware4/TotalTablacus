package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"

	"gopkg.in/yaml.v2"
)

func main() {
	var (
		cur string
	)
	flag.StringVar(&cur, "cur", "", "current dir path")
	flag.Parse()
	os.Exit(run(cur))
}

type Rule struct {
	Prefix      string
	Description string
}

type Menu struct {
	options []Rule
}

func (m *Menu) load(path string) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return
	}
	rules := []Rule{}
	if err := yaml.Unmarshal(buf, &rules); err != nil {
		return
	}
	m.options = rules
}

func (m Menu) pick() (string, error) {
	if len(m.options) < 1 {
		return "", nil
	}
	idx, err := fuzzyfinder.Find(m.options, func(i int) string {
		o := m.options[i]
		return fmt.Sprintf("%s - %s", o.Prefix, o.Description)
	})
	if err != nil {
		return "", err
	}
	return m.options[idx].Prefix, nil
}

type DirName struct {
	timestamp string
	prefix    string
	name      string
}

func (dn DirName) getFullPrefix() (prefix string) {
	prefix = dn.timestamp + "_"
	if 0 < len(dn.prefix) {
		prefix = prefix + dn.prefix + "_"
	}
	return
}

func (dn *DirName) setName() error {
	p := dn.getFullPrefix()
	fmt.Printf("Enter after '%s': ", p)
	var answer string
	scn := bufio.NewScanner(os.Stdin)
	if !scn.Scan() {
		return fmt.Errorf("abort")
	}
	answer = scn.Text()
	answer = strings.TrimSpace(answer)
	dn.name = answer
	return nil
}

func (dn DirName) getName() string {
	p := dn.getFullPrefix()
	if len(dn.name) < 1 {
		return strings.TrimSuffix(p, "_")
	}
	return p + dn.name
}

type WorkDir struct {
	path string
}

func (wd WorkDir) getExample(prefix string) (examples []string) {
	if len(prefix) < 1 {
		return
	}
	items, err := os.ReadDir(wd.path)
	if err != nil {
		return
	}
	for _, item := range items {
		n := item.Name()
		if item.IsDir() && 9 < len(n) {
			suf := n[9:]
			if strings.HasPrefix(suf, prefix+"_") {
				examples = append(examples, n)
			}
		}
	}
	return
}

func (wd WorkDir) showExamples(prefix string) {
	es := wd.getExample(prefix)
	if 0 < len(es) {
		if 1 < len(es) {
			fmt.Println("Examples:")
		} else {
			fmt.Println("Example:")
		}
		for _, e := range es {
			fmt.Printf("- '%s'\n", e)
		}
	}
}

func (wd WorkDir) newDir(name string) error {
	np := filepath.Join(wd.path, name)
	if _, err := os.Stat(np); err == nil {
		return fmt.Errorf("'%s' already exists in '%s'", name, wd.path)
	}
	return os.Mkdir(np, os.ModePerm)
}

func run(path string) int {
	y := filepath.Join(path, "rule.yml")
	var menu Menu
	menu.load(y)

	p, err := menu.pick()
	if err != nil && err != fuzzyfinder.ErrAbort {
		return 1
	}

	cur := WorkDir{path: path}
	cur.showExamples(p)

	ts := time.Now().Format("20060102")
	dn := DirName{timestamp: ts, prefix: p}
	if err := dn.setName(); err != nil {
		return 1
	}

	n := dn.getName()
	if err := cur.newDir(n); err != nil {
		return 1
	}
	return 0
}
