package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/projectdiscovery/gologger"
	"github.com/scottkiss/gosshtool"
)

const banner = `

              __         __              __            
   __________/ /_  _____/ /_  ___  _____/ /_____  _____
  / ___/ ___/ __ \/ ___/ __ \/ _ \/ ___/ //_/ _ \/ ___/
 (__  |__  ) / / / /__/ / / /  __/ /__/ ,< /  __/ /    
/____/____/_/ /_/\___/_/ /_/\___/\___/_/|_|\___/_/v1.0     
                                                       
`
const Version = `1.0`

func showBanner() {
	gologger.Printf("%s\n", banner)
	gologger.Printf("\t\tCreated by Shazeb.\n\n")

	gologger.Labelf("Use with caution. You are responsible for your actions\n")
	gologger.Labelf("Developers assume no liability and are not responsible for any misuse or damage.\n\n")
}

var (
	userList     string
	passwordList string
	//Version      bool
	concurrency int
)

func ParseOptions() {

	flag.StringVar(&userList, "U", "", "List of the default usernames of ssh")
	flag.StringVar(&passwordList, "P", "", "List of the default passwords of ssh")
	//flag.BoolVar(&Version, "version", false, "Show the version of sshchecker.")
	flag.IntVar(&concurrency, "c", 10, "set the concurrency level")
	flag.Parse()

}

func main() {
	ParseOptions()
	//	if Version {
	//		gologger.Infof("Current Version: %s\n", Version)
	//		os.Exit(0)
	//	}

	showBanner()
	f, err := os.Open(userList)
	f.Close()
	reader := bufio.NewReader(f)
	var user string
	for {
		user, err = reader.ReadString('\n')
		if err != nil {
			break
		}
	}
	var wg sync.WaitGroup

	//adding scanner for reading the input from terminal
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		text := sc.Text()
		wg.Add(1)

		go func(ip string, user string) {
			fmt.Printf("Trying sshing on: %v \n %s \n", ip, user)
			sshconfig := &gosshtool.SSHClientConfig{
				User:     user,
				Password: passwordList,
				Host:     ip,
			}
			sshclient := gosshtool.NewSSHClient(sshconfig)
			_, err := sshclient.Connect()
			if err == nil {
				fmt.Println("ssh successful")
			} else {
				fmt.Println("ssh failed")
			}
			wg.Done()
		}(text)
	}
	wg.Wait()
}
