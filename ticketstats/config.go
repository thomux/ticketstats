package ticketstats

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config groups all configuration values.
type Config struct {
	Template string
	Types    ConfigTypeNames
	States   ConfigStateNames
	Customs  ConfigCustomFields
	Formats  ConfigFormats
}

// ConfigFormats groups format strings.
type ConfigFormats struct {
	Date     string
	JiraDate string
}

// ConfigTypeNames groups the type name strings.
type ConfigTypeNames struct {
	Feature     string
	Bug         string
	Improvement string
}

// ConfigStateNames groups the state name strings.
type ConfigStateNames struct {
	Closed string
}

// ConfigCustomFields groups the custom field names.
type ConfigCustomFields struct {
	ExternalId        string
	SupplierReference string
	Variant           string
	Account           string
	Category          string
}

// DefaultConfig creates a new Config with all settings initialized using
// default values.
func DefaultConfig() Config {
	var config Config

	config.Template = ""

	config.Formats.Date = "2006-01-02"
	config.Formats.JiraDate = "02/Jan/06 3:04 PM"

	config.Types.Bug = "Bug"
	config.Types.Feature = "New Feature"
	config.Types.Improvement = "Improvement"

	config.States.Closed = "Closed"

	config.Customs.ExternalId = "Custom field (External ID)"
	config.Customs.SupplierReference = "Custom field (Supplier reference)"
	config.Customs.Variant = "Custom field (ICAS Variant)"
	config.Customs.Account = "Custom field (Booking Account)"
	config.Customs.Category = "Custom field (Bug-Category)"

	return config
}

// saveConfig saves the default config as "config.json".
func saveConfig() {
	data, err := json.MarshalIndent(DefaultConfig(), "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("config.json", data, 0644)
	if err != nil {
		panic(err)
	}
}

// LoadConfig loads the config or returns a default config if loading fails.
func LoadConfig() Config {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("INFO: no config file found")
		return DefaultConfig()
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("ERROR: Config:", err)
		return DefaultConfig()
	}
	return config
}
