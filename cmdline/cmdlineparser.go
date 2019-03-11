package cmdline

import "flag"

const (
	hostFlag     string = "h"
	hostFlagHint string = "host logger will try to connect to"

	portFlag     string = "p"
	portFlagHint string = "hosts port"

	fileFlag     string = "f"
	fileFlagHint string = "file to save log to"
)

type HostConfig struct {
	Host *string
	Port *string
	File *string
}

func GetSettings() *HostConfig {
	hc := HostConfig{}
	hc.Host = flag.String(hostFlag, "", hostFlagHint)
	hc.Port = flag.String(portFlag, "", portFlagHint)
	hc.File = flag.String(fileFlag, "", fileFlagHint)
	flag.Parse()
	return &hc
}

func PrintUsage() {
	flag.Usage()
}

func (hc *HostConfig) GetAddress() string {
	return *hc.Host + ":" + *hc.Port
}
