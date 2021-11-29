package models

type Response struct {
	Success bool         `json:"success,omitempty" bson:"success,omitempty"`
	Status  string       `json:"status,omitempty" bson:"status,omitempty"`
	Msg     string       `json:"msg,omitempty" bson:"msg,omitempty"`
	Data    []Descriptor `json:"data,omitempty" bson:"data,omitempty"`
}
type Descriptor struct {
	Id      string     `json:"id" bson:"id"`
	Type    string     `json:"type,omitempty" bson:"type,omitempty"`
	Name    string     `json:"name,omitempty" bson:"name,omitempty"`
	Version string     `json:"version,omitempty" bson:"version,omitempty"`
	Modules []string   `json:"modules,omitempty" bson:"modules,omitempty"`
	Configs []Config   `json:"configs" bson:"configs"`
	Status  []Statuses `json:"status" bson:"status"`
}
type Config struct {
	Id       string      `json:"id" bson:"id"`
	Type     string      `json:"type,omitempty" bson:"type,omitempty"`
	Name     string      `json:"name,omitempty" bson:"name,omitempty"`
	Protocol []Protocols `json:"protocol" bson:"protocol"`
}

type Statuses struct {
	Id       string      `json:"id" bson:"id"`
	Type     string      `json:"type,omitempty" bson:"type,omitempty"`
	Name     string      `json:"name" bson:"name,omitempty"`
	Protocol []Protocols `json:"protocol" bson:"protocol"`
}

type Protocols struct {
	Id    string  `json:"id" bson:"id"`
	Type  string  `json:"type,omitempty" bson:"type,omitempty"`
	Name  string  `json:"name" bson:"name,omitempty"`
	Items []Items `json:"items" bson:"items"`
}

type Items struct {
	Id       string `json:"id" bson:"id"`
	Type     string `json:"type,omitempty" bson:"type,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Priority string `json:"priority,omitempty" bson:"priority,omitempty"`
	Default  struct {
		Max int `json:"max,omitempty" bson:"max,omitempty"`
		Min int `json:"min,omitempty" bson:"min,omitempty"`
	} `json:"default,omitempty" bson:"default,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Ui          bool   `json:"ui,omitempty" bson:"ui,omitempty"`
	Persist     bool   `json:"persist,omitempty" bson:"persist,omitempty"`
}
