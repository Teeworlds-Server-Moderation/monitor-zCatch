package main

import (
	configo "github.com/jxsl13/simple-configo"
)

const (
	// https://stackoverflow.com/questions/53497/regular-expression-that-matches-valid-ipv6-addresses
	// https://stackoverflow.com/questions/12968093/regex-to-validate-port-number
	ipPortRegex = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))|((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]):([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])`

	brokerAddressRegex = `[a-z0-9-\.:]+:([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])`
)

// Config configures the monitoring application
type Config struct {
	EconAddress    string
	EconPassword   string
	BrokerAddress  string
	BrokerUsername string
	BrokerPassword string
}

// Name is the name of the configuration Cache
func (c *Config) Name() (name string) {
	return "zCatch Monitoring"
}

// Options returns a list of available options that can be configured for this
// config object
func (c *Config) Options() (options configo.Options) {

	optionsList := configo.Options{
		{
			Key:           "MONITOR_ECON_ADDRESS",
			Description:   "Please provide the address of your configured server Econ: <IP>:<Port>",
			Mandatory:     true,
			ParseFunction: configo.DefaultParserRegex(&c.EconAddress, ipPortRegex, "Please provide a valid IP:Port address"),
		},
		{
			Key:           "MONITOR_ECON_PASSWORD",
			Description:   "The password to log into your Econ.",
			Mandatory:     true,
			ParseFunction: configo.DefaultParserString(&c.EconPassword),
		},
		{
			Key:           "MONITOR_BROKER_ADDRESS",
			Description:   "The address of your broker in the container is rabbitmq:5672",
			DefaultValue:  "localhost:5672",
			ParseFunction: configo.DefaultParserRegex(&c.BrokerAddress, brokerAddressRegex, "Please provide a valid <IP>:<Port> address'"),
		},
		{
			Key:           "MONITOR_BROKER_USER",
			Description:   "The user that can access the broker, default: tw-admin",
			DefaultValue:  "tw-admin",
			ParseFunction: configo.DefaultParserString(&c.BrokerUsername),
		},
		{
			Key:           "MONITOR_BROKER_PASSWORD",
			Mandatory:     true,
			Description:   "The password to access the broker with the corresonding username.",
			ParseFunction: configo.DefaultParserString(&c.BrokerPassword),
		},
	}
	return optionsList
}
