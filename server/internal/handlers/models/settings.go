package models

type SettingsDataResponse struct {
	Groups  []GroupResponse  `json:"groups"`
	Members []MemberResponse `json:"members"`
	Albums  []AlbumResponse  `json:"albums"`
}
