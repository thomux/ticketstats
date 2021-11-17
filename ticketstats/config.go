package ticketstats

type Config struct {
	Template string
	Types    ConfigTypeNames
	States   ConfigStateNames
	Customs  ConfigCustomFields
}

type ConfigTypeNames struct {
	Feature     string
	Bug         string
	Improvement string
}

type ConfigStateNames struct {
	Closed string
}

type ConfigCustomFields struct {
	Account string
}

func DefaultConfig() Config {
	var config Config

	config.Template = ""

	config.Types.Bug = "Bug"
	config.Types.Feature = "New Feature"
	config.Types.Improvement = "Improvement"

	config.States.Closed = "Closed"

	config.Customs.Account = "Custom field (Booking Account)"

	return config
}
