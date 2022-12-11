package main

import (
	"github.com/whiterabb17/brutebar/bruteshed"
)

func main() {
	bruteshed.SshBrute("usrpsw.list", "ip.list", "10", "60")
}
