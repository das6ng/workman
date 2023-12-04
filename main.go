package main

import (
	"flag"
	"path/filepath"

	"github.com/das6ng/workman/mgr"
)

var (
	fNoCUI  bool
	fList   bool
	fUpdate string
)

var out mgr.Printer

func main() {
	flag.BoolVar(&fNoCUI, "n", false, "no command line ui")
	flag.BoolVar(&fList, "l", false, "list go.work info")
	flag.StringVar(&fUpdate, "u", "{}", "update go.work")
	flag.Parse()

	out.SetOutput(func() mgr.Output {
		if fNoCUI {
			return mgr.OutputJSON
		}
		return mgr.OutputCUI
	}())
	cur, _ := filepath.Abs(".")
	mgr := &mgr.WorkManager{Printer: out}
	if err := mgr.Load(cur); err != nil {
		out.Msg("cannot load workspace", err)
		return
	}

	// command line ui mode
	if !fNoCUI {
		mgr.CUIMain()
		return
	}

	// arg mode
	if fList {
		out.Print(mgr.GetInfo())
		return
	}
	if fUpdate != "{}" && fUpdate != "" {
		out.Print(mgr.ArgUpdate(fUpdate))
		return
	}
	out.Msg("nothing to do")
}
