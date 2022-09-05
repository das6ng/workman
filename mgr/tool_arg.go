package mgr

import (
	"encoding/json"
)

type updateMod struct {
	Add  []string `json:"add"`
	Drop []string `json:"drop"`
}

type result struct {
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
	Err string `json:"err"`
}

type workInfo struct {
	result
	Path  string   `json:"path"`
	GoVer string   `json:"go_ver"`
	Used  []string `json:"used"`
	Total []string `json:"total"`
}

func (m *WorkManager) ArgUpdate(info string) result {
	up := updateMod{}
	if err := json.Unmarshal([]byte(info), &up); err != nil {
		return result{
			OK:  false,
			Msg: "parse update info",
			Err: err.Error(),
		}
	}
	if err := m.Update(up.Add, up.Drop); err != nil {
		return result{
			OK:  false,
			Msg: "update go.work",
			Err: err.Error(),
		}
	}
	return result{OK: true}
}

func (m *WorkManager) GetInfo() workInfo {
	if err := m.check(); err != nil {
		return workInfo{result: result{
			OK:  false,
			Msg: "check workspace failed",
			Err: err.Error(),
		}}
	}
	used := m.getUsedModules()
	nearby := m.getNearbyModules()
	total := uniqStrList(append(used, nearby...))
	// sort.Strings(used)
	// sort.Strings(total)
	return workInfo{
		result: result{
			OK: true,
		},
		Path:  m.workFilePath,
		GoVer: m.workFile.Go.Version,
		Used:  used,
		Total: total,
	}
}
