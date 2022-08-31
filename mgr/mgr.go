package mgr

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

type WorkManager struct {
	Printer
	workFilePath string
	workFileDir  string
	workFile     *modfile.WorkFile
}

func (m *WorkManager) Load(dir string) error {
	m.workFilePath = findWorkspaceFile(dir)
	if m.workFilePath == "" {
		return fmt.Errorf("no go.work file found: %s", dir)
	}
	workFile, err := readWorkFile(m.workFilePath)
	if err != nil {
		return fmt.Errorf("read work file: %s", err.Error())
	}
	m.workFile = workFile
	m.workFileDir = filepath.Dir(m.workFilePath)
	return nil
}

func (m *WorkManager) AddUse(dir string) error {
	modFilePath := filepath.Join(m.workFileDir, dir, "go.mod")
	if fi, err := os.Stat(modFilePath); err != nil || fi.IsDir() {
		return fmt.Errorf("can find go.mod: %s", modFilePath)
	}
	if err := m.check(); err != nil {
		return err
	}
	modFile, err := readModFile(modFilePath)
	if err != nil {
		return err
	}
	return m.workFile.AddUse(dir, modFile.Module.Mod.Path)
}

func (m *WorkManager) FindUse(dir string) (*modfile.Use, bool) {
	for _, u := range m.workFile.Use {
		if u.Path == dir {
			return u, true
		}
	}
	return nil, false
}

func (m *WorkManager) DropUse(dir string) error {
	return m.workFile.DropUse(dir)
}

func (m *WorkManager) Write() error {
	if err := m.check(); err != nil {
		return err
	}
	return writeWorkFile(m.workFilePath, m.workFile)
}

func (m *WorkManager) Update(add, drop []string) error {
	w := false
	if len(add) > 0 {
		for _, u := range add {
			if err := m.AddUse(u); err != nil {
				return err
			}
		}
		w = true
	}
	if len(drop) > 0 {
		for _, u := range drop {
			if err := m.DropUse(u); err != nil {
				return err
			}
		}
		w = true
	}
	if w {
		return m.Write()
	}
	return nil
}

func (m *WorkManager) check() error {
	if m.workFile == nil {
		return fmt.Errorf("no work file loaded")
	}
	return nil
}

func (m *WorkManager) getUsedModules() []string {
	mods := make([]string, 0, len(m.workFile.Use))

	for _, u := range m.workFile.Use {
		mods = append(mods, u.Path)
	}

	return mods
}

func (m *WorkManager) getNearbyModules() []string {
	mods := make([]string, 0, 3)
	di, _ := os.ReadDir(m.workFileDir)
	for _, d := range di {
		if !d.IsDir() {
			continue
		}
		if fi, err := os.Stat(path.Join(m.workFileDir, d.Name(), "go.mod")); err != nil || fi.IsDir() {
			continue
		}
		mods = append(mods, d.Name())
	}
	return mods
}
