package component

import "github.com/yohamta/donburi"

type PlayerData struct {
	rawData map[string]any
}

func NewPlayerData() *PlayerData {
	return &PlayerData{rawData: map[string]any{}}
}
func (p *PlayerData) Has(key string) any {
	ok, _ := p.rawData[key]
	return ok
}
func (p *PlayerData) Get(key string) any {
	return p.rawData[key]
}
func (p *PlayerData) Set(key string, val any) {
	p.rawData[key] = val
}

var PlayerDataComponent = donburi.NewComponentType[PlayerData]()
