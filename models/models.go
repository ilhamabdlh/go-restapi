package models

type Config struct {
	
	Id		string			`json:"id,omitempty" bson:"id,omitempty"`
	Type   string             `json:"type,omitempty" bson:"type,omitempty"`
	Name  	string             `json:"name" bson:"name,omitempty"`
	Protocol string            `json:"protocol" bson:"protocol,omitempty"`
}

type Status struct {
	
	Id		string			`json:"id,omitempty" bson:"id,omitempty"`
	Type   string             `json:"type,omitempty" bson:"type,omitempty"`
	Name  	string             `json:"name" bson:"name,omitempty"`
	Protocol string            `json:"protocol" bson:"protocol,omitempty"`
}

type Protocol struct {
	Id	   string `json:"id,omitempty" bson:"id,omitempty"`
	Type  string `json:"type,omitempty" bson:"type,omitempty"`
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Items string `json:"items,omitempty" bson:"items,omitempty"`
}
type Descriptor struct {
	Id	   string `json:"id,omitempty" bson:"id,omitempty"`
	Type  string `json:"type,omitempty" bson:"type,omitempty"`
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Version string `json:"version,omitempty" bson:"version,omitempty"`
	Modules string `json:"modules,omitempty" bson:"modules,omitempty"`
	Configs string `json:"configs,omitempty" bson:"configs,omitempty"`
	Status string `json:"status,omitempty" bson:"status,omitempty"`
}