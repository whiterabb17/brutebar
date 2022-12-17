package smbdancer

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLIConfig struct {
	host     string
	port     int
	debug    bool
	domain   string
	threads  int
	sleep    string
	userFile string
	pwdFile  string
}

type Config struct {
	host    string
	port    int
	debug   bool
	domain  string
	threads int
	sleep   float64
	users   *WordlistInput
	passwds *WordlistInput
}

const (
	BANNER = ` 
 
`
	SEP = "-----------------------------------------------------"
)

func createConfig(cliconf *CLIConfig) (Config, error) {
	var conf Config
	var err error
	if cliconf.host == "" {
		return conf, fmt.Errorf("Host (-h) not defined")
	}
	conf.host = cliconf.host
	if cliconf.domain == "" {
		return conf, fmt.Errorf("Domain (-d) not defined")
	}
	conf.domain = cliconf.domain
	if cliconf.userFile == "" {
		return conf, fmt.Errorf("Userfile (-u) not defined")
	}
	conf.users, err = NewWordlistInput(cliconf.userFile)
	if err != nil {
		return conf, fmt.Errorf("Could not read user file: %s", err)
	}
	if cliconf.pwdFile == "" {
		return conf, fmt.Errorf("Passwordfile (-w) not defined")
	}
	conf.passwds, err = NewWordlistInput(cliconf.pwdFile)
	if err != nil {
		return conf, fmt.Errorf("Could not read password file: %s", err)
	}
	if cliconf.sleep != "" {
		conf.sleep, err = strconv.ParseFloat(cliconf.sleep, 64)
		if err != nil {
			return conf, fmt.Errorf("Erroneus sleep (-s) value")
		}
	}

	conf.threads = cliconf.threads
	conf.port = cliconf.port
	conf.debug = cliconf.debug
	return conf, nil
}

func SmbBrute(host string, port int, threadCnt int, userlist string, passlist string, domain string, debug bool, sleep string) {
	var cliconf CLIConfig
	cliconf.host = host
	cliconf.port = port
	cliconf.threads = threadCnt
	cliconf.userFile = userlist
	cliconf.pwdFile = passlist
	cliconf.domain = domain
	cliconf.debug = debug
	cliconf.sleep = sleep
	conf, err := createConfig(&cliconf)
	if err != nil {
		fmt.Printf("  [!] Error: %s\n\n", err)
		flag.Usage()
		os.Exit(1)
	}
	runner := NewRunner(&conf)
	runner.Start()
}
