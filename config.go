package main

import (
	"time"

	configo "github.com/jxsl13/simple-configo"
)

const (
	brokerAddressRegex = `^[a-z0-9-\.:]+:([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`
	redisAddressRegex  = `^[a-z0-9-\.:]+:([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`

	folderRegex  = `^[a-zA-Z0-9-]+$`
	errFolderMsg = "The folder name must only contain alphanumeric characters, no special caracters nor whitespaces and must not be empty."
)

// Config is the configuration for this microservice
type Config struct {
	BrokerAddress  string
	BrokerUsername string
	BrokerPassword string

	DataPath                 string
	BroadcastNonAbortActions bool
	DefaultVotebanCommand    string
	DefaultVotebanDuration   time.Duration

	BackupInterval time.Duration
}

// Name is the name of the configuration Cache
func (c *Config) Name() (name string) {
	return "Detect VPN"
}

// Options returns a list of available options that can be configured for this
// config object
func (c *Config) Options() (options configo.Options) {

	// this default value allows for local development
	// while the defalt environment value in the Dockerfile allows for overriding this
	// default value when running in a container.
	optionsList := configo.Options{
		{
			Key:           "BROKER_ADDRESS",
			Description:   "The address of the message broker. In the container environemt it's rabbitmq:5672",
			Mandatory:     true,
			DefaultValue:  "localhost:5672",
			ParseFunction: configo.DefaultParserRegex(&c.BrokerAddress, brokerAddressRegex, "BROKER_ADDRESS must have the format <hostname/ip>:<port>"),
		},
		{
			Key:           "BROKER_USER",
			Description:   "Username of the broker user",
			Mandatory:     true,
			DefaultValue:  "tw-admin",
			ParseFunction: configo.DefaultParserString(&c.BrokerUsername),
		},
		{
			Key:           "BROKER_PASSWORD",
			Mandatory:     true,
			Description:   "Password of the specified username",
			ParseFunction: configo.DefaultParserString(&c.BrokerPassword),
		},
		{
			Key:           "DATA_PATH",
			Description:   "Is the root folder that contains all of the data of this service.",
			DefaultValue:  "data",
			ParseFunction: configo.DefaultParserString(&c.DataPath),
		},
		{
			Key:           "BROADCAST_NON_ABORT_ACTIONS",
			Description:   "If a funvoter is detected, one may use their IP to voteban them on every connected server.",
			DefaultValue:  "false",
			ParseFunction: configo.DefaultParserBool(&c.BroadcastNonAbortActions),
		},
		{
			Key:           "DEFAULT_VOTEBAN_DURATION",
			Description:   "You may use the variables {IP}, {ID}, {DURATION:MINUTES}, {DURATION:SECONDS}, {REASON}",
			DefaultValue:  "30m",
			ParseFunction: configo.DefaultParserDuration(&c.DefaultVotebanDuration),
		},
		{
			Key:           "DEFAULT_VOTEBAN_COMMAND",
			Description:   "You may use the variables {IP}, {ID}, {DURATION:MINUTES}, {DURATION:SECONDS}, {REASON}",
			DefaultValue:  "voteban {IP} {DURATION:SECONDS}",
			ParseFunction: configo.DefaultParserString(&c.DefaultVotebanCommand),
		},
		{
			Key:           "CSV_BACKUP_INTERVAL",
			Description:   "Interval after which a new csv file is created that contains new unclissified reason lines",
			DefaultValue:  "24h",
			ParseFunction: configo.DefaultParserDuration(&c.BackupInterval),
		},
	}

	return optionsList
}
