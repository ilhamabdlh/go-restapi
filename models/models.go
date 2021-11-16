package models

type Config struct {
	Id       string `json:"id" bson:"id"`
	Type     string `json:"type,omitempty" bson:"type,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Protocol []struct{
		*Protocol
	} `json:"protocol" bson:"protocol"`
}

type Status struct {
	Id       string `json:"id" bson:"id"`
	Type     string `json:"type,omitempty" bson:"type,omitempty"`
	Name     string `json:"name" bson:"name,omitempty"`
	Protocol []struct{
		*Protocol
	} `json:"protocol" bson:"protocol"`
}

type Protocol struct {
	Id    string `json:"id" bson:"id"`
	Type  string `json:"type,omitempty" bson:"type,omitempty"`
	Name  string `json:"name" bson:"name,omitempty"`
	Items []struct{
		*Device
	} `json:"items" bson:"items"`
}

type Descriptor struct {
	Id      string `json:"id" bson:"id"`
	Type    string `json:"type,omitempty" bson:"type,omitempty"`
	Name    string `json:"name,omitempty" bson:"name,omitempty"`
	Version string `json:"version,omitempty" bson:"version,omitempty"`
	Modules string `json:"modules,omitempty" bson:"modules,omitempty"`
	Configs struct{
		*Config
	} `json:"configs,omitempty" bson:"configs,omitempty"`
	Status struct{
		*Status
	} `json:"status,omitempty" bson:"status,omitempty"`
}

type Device struct {
	Id          string `json:"id" bson:"id"`
	Type        string `json:"type,omitempty" bson:"type,omitempty"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Priority         string    `json:"priority,omitempty" bson:"priority,omitempty"`
	Default struct{
		Max         int    `json:"max,omitempty" bson:"max,omitempty"`
		Min         int    `json:"min,omitempty" bson:"min,omitempty"`
	} `json:"default,omitempty" bson:"default,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Ui          bool   `json:"ui,omitempty" bson:"ui,omitempty"`
	Persist     bool   `json:"persist" bson:"persist"`
}

type Response struct {
	Success string `json:"success,omitempty" bson:"success,omitempty"`
	Msg     string `json:"msg,omitempty" bson:"msg,omitempty"`
	Data    string `json:"data,omitempty" bson:"data,omitempty"`
}
