package dsl

type Resources map[string]Resource

type Resource struct {
	Get    GetMethod    `json:"get,omitzero"`
	Post   PostMethod   `json:"post,omitzero"`
	Put    PutMethod    `json:"put,omitzero"`
	Delete DeleteMethod `json:"delete,omitzero"`
}
