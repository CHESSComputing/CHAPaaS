package main

import "encoding/json"

// Workflow define CHAP workflow record
type Workflow struct {
	MetaData    map[string]interface{} `json:"meta_data"`   // meta-data information about ML model
	Name        string                 `json:"name"`        // workflow name
	Type        string                 `json:"type"`        // workflow type
	Group       string                 `json:"group"`       // workflow group
	Version     string                 `json:"version"`     // workflow version
	Description string                 `json:"description"` // workflow description
	Reference   string                 `json:"reference"`   // reference URL
	UserName    string                 `json:"user_name"`   // user name
	UserID      string                 `json:"user_id"`     // user id
}

// ToJSON provides string representation of Record
func (r Workflow) ToJSON() string {
	// create pretty JSON representation of the record
	data, _ := json.MarshalIndent(r, "", "    ")
	return string(data)
}
