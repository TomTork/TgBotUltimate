package SyncStrapi

import "TgBotUltimate/types/Sync"

type Strapi struct {
	Projects  []Sync.Project  `json:"projects"`
	Buildings []Sync.Building `json:"buildings"`
	Sections  []Sync.Section  `json:"sections"`
	Flats     []Sync.Flat     `json:"flats"`
}
