package main

import (
	"encoding/json"
	"log"
	"time"
)

// Workflow define CHAP workflow record
type Workflow struct {
	MetaData    map[string]interface{} `json:"meta_data",yaml:"meta_data"`     // meta-data information about ML model
	Name        string                 `json:"name",yaml:"name"`               // workflow name
	Directory   string                 `json:"directory",yaml:"directory"`     // workflow directory
	Config      string                 `json:"config",yaml:"config"`           // workflow configuration
	Type        string                 `json:"type",yaml:"type"`               // workflow type
	Group       string                 `json:"group",yaml:"group"`             // workflow group
	Version     string                 `json:"version",yaml:"version"`         // workflow version
	Description string                 `json:"description",yaml:"description"` // workflow description
	Reference   string                 `json:"reference",yaml:"reference"`     // reference URL
	UserName    string                 `json:"user_name",yaml:"user_name"`     // user name
	UserID      string                 `json:"user_id",yaml:"user_id"`         // user id
}

// ToJSON provides string representation of Record
func (r Workflow) ToJSON() string {
	// create pretty JSON representation of the record
	data, _ := json.MarshalIndent(r, "", "    ")
	return string(data)
}

// ChapWorkflows represent current set of known CHAP workflows
type ChapWorkflows struct {
	Workflows []Workflow
	Expire    time.Time
}

// helper function to get current set of chap workflows
func (c ChapWorkflows) getWorkflows() []Workflow {
	// update local cache of workflows
	if !c.Expire.After(time.Now()) {
		if Config.Verbose > 0 {
			log.Println("UPDATE CHAP workflows")
		}
		c.Workflows = getChapWorkflows()
		c.Expire = time.Now().Add(1 * time.Hour)
	}
	if Config.Verbose > 0 {
		log.Printf("Total number of existing CHAP workflows %d, expire at %v", len(c.Workflows), c.Expire.String())
	}
	return c.Workflows
}

var chapWorkflows ChapWorkflows
