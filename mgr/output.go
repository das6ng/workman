package mgr

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pterm/pterm"
)

type Output int

const (
	OutputCUI Output = iota
	OutputJSON
)

type Printer struct {
	output Output
}

func (p *Printer) SetOutput(o Output) {
	p.output = o
}

func (p Printer) Msg(msg string, err ...error) {
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	switch p.output {
	case OutputCUI:
		if e != nil {
			msg = pterm.NewRGB(241, 76, 76).Sprintfln("ERROR %s %s", msg, e.Error())
		} else {
			msg = pterm.NewStyle(pterm.FgCyan).Sprintln(msg)
		}
		area, _ := pterm.DefaultArea.Start()
		area.Update(msg)
	case OutputJSON:
		r := result{
			OK:  e == nil,
			Msg: msg,
			Err: func() string {
				if e != nil {
					return e.Error()
				}
				return ""
			}(),
		}
		bs, _ := json.Marshal(r)
		print(string(bs))
	}
}

func (p Printer) Print(v any) {
	switch p.output {
	case OutputCUI:
		log.Println(v)
	case OutputJSON:
		bs, e := json.Marshal(v)
		if e != nil {
			bs, _ = json.Marshal(result{
				OK:  true,
				Msg: fmt.Sprint(v),
			})
		}
		print(string(bs))
	}
}
