package mgr

import (
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func findWorkspaceFile(dir string) (root string) {
	if dir == "" {
		panic("dir not set")
	}
	dir = filepath.Clean(dir)

	// Look for enclosing go.mod.
	for {
		f := filepath.Join(dir, "go.work")
		if fi, err := os.Stat(f); err == nil && !fi.IsDir() {
			return f
		}
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		if gr, ok := os.LookupEnv("GOROOT"); ok && gr == d {
			// As a special case, don't cross GOROOT to find a go.work file.
			// The standard library and commands built in go always use the vendored
			// dependencies, so avoid using a most likely irrelevant go.work file.
			return ""
		}
		dir = d
	}
	return ""
}

func readWorkFile(path string) (*modfile.WorkFile, error) {
	workData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return modfile.ParseWork(path, workData, nil)
}

func readModFile(path string) (*modfile.File, error) {
	modData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return modfile.Parse(path, modData, nil)
}

func writeWorkFile(path string, wf *modfile.WorkFile) error {
	wf.SortBlocks()
	wf.Cleanup()
	out := modfile.Format(wf.Syntax)

	return os.WriteFile(path, out, 0666)
}

func uniqStrList(in []string) []string {
	m := make(map[string]struct{}, len(in))
	for i := 0; i < len(in); i++ {
		if _, ok := m[in[i]]; ok {
			in = append(in[:i], in[i+1:]...)
			i--
			continue
		}
		m[in[i]] = struct{}{}
	}
	return in
}

func subStrList(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, bb := range b {
		mb[bb] = struct{}{}
	}
	for i := 0; i < len(a); i++ {
		if _, ok := mb[a[i]]; ok {
			a = append(a[:i], a[i+1:]...)
			i--
			continue
		}
	}
	return a
}
