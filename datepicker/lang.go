package datepicker

import (
	_ "embed"
	"encoding/json"
)

// LangsData contains all languages data
// key is language code, value is language data with key:value pairs
type LangsData map[string]map[string]string

//go:embed langs.json
var langsData string

func loadLangs() LangsData {
	data := LangsData{}
	_ = json.Unmarshal([]byte(langsData), &data)
	return data
}

func (d *DatePicker) lang(key string) string {
	l, ok := d.langs[d.language]
	if !ok {
		l = d.langs["en"]
	}
	if l == nil {
		return "<" + key + ">"
	}
	s, ok2 := l[key]
	if !ok2 {
		s = "<" + key + ">"
	}
	return s
}
