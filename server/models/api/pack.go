package api

type PackRequest struct {
	Items int `json:"items"`
}

type PackResponse struct {
	Packs map[int]int `json:"packs"`
}
