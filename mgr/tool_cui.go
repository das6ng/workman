package mgr

import (
	"log"

	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
)

func (m *WorkManager) CUIMain() {
	cursor.Hide()
	defer cursor.Show()
	for {
		println("~")
		m.ShowUsed()
		var add, drop []string
		exit := false
		if err := keyboard.Listen(func(key keys.Key) (stop bool, err error) {
			switch key.Code {
			case keys.RuneKey:
				if key.String() == "q" || key.String() == "Q" {
					exit = true
					return true, nil
				}
			case keys.Space, keys.Tab:
				add, drop, err = m.ShowSelectModules()
				return true, err
			case keys.Esc, keys.CtrlC:
				exit = true
				return true, nil
			}
			return false, nil
		}); err != nil {
			log.Println("keyboard listen err:", err)
		}
		if exit {
			m.Printer.Msg("Bye")
			break
		}
		if err := m.Update(add, drop); err != nil {
			m.Printer.Msg("update go.work", err)
		}
	}
}

func (m *WorkManager) ShowUsed() {
	used := m.getUsedModules()
	nearby := m.getNearbyModules()
	items := make([]pterm.BulletListItem, 0, len(used)+len(nearby))
	for _, u := range used {
		items = append(items, pterm.BulletListItem{
			Level:       0,
			Text:        u,
			TextStyle:   pterm.NewStyle(pterm.FgBlue),
			Bullet:      "√",
			BulletStyle: pterm.NewStyle(pterm.FgGreen),
		})
	}
	for _, n := range subStrList(nearby, used) {
		items = append(items, pterm.BulletListItem{
			Level:       1,
			Text:        n,
			TextStyle:   pterm.NewStyle(pterm.FgGray),
			Bullet:      "·",
			BulletStyle: pterm.NewStyle(pterm.FgRed),
		})
	}
	text, _ := pterm.DefaultBulletList.WithItems(items).Srender()
	area, _ := pterm.DefaultArea.Start()
	area.Update(
		pterm.Sprintfln("Workspace: %s go%s", m.workFilePath, m.workFile.Go.Version),
		"Module usage:\n",
		text,
		pterm.ThemeDefault.SecondaryStyle.Sprintfln(
			"change: <Tab> | exit: <Esc>/<q>",
		),
	)
}

func (m *WorkManager) ShowSelectModules() (add, drop []string, err error) {
	used := m.getUsedModules()
	nearby := m.getNearbyModules()
	opts := loadSortedMods(used, nearby).Mods()
	sel := &pterm.InteractiveMultiselectPrinter{
		TextStyle:      &pterm.ThemeDefault.PrimaryStyle,
		DefaultText:    "Modules",
		Options:        opts,
		OptionStyle:    &pterm.ThemeDefault.DefaultText,
		DefaultOptions: used,
		MaxHeight:      9,
		Selector:       ">",
		SelectorStyle:  &pterm.ThemeDefault.SecondaryStyle,
	}
	selected, err := sel.Show()
	if err != nil {
		return
	}
	mm := make(map[string]bool, len(used))
	for _, s := range used {
		mm[s] = false
	}
	for _, s := range selected {
		if _, ok := mm[s]; ok {
			mm[s] = true
			continue
		}
		add = append(add, s)
	}
	for s, ok := range mm {
		if !ok {
			drop = append(drop, s)
		}
	}
	return
}
