package main

import (
	"bufio"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	helpMsg = `This is a golang implementation of the RBC Client, using the TCP connection

Command Usages: /command [arguments]

The commands are:
	exit        - exit the program
`
)

func handleConnection(c net.Conn) {
	defer c.Close()

	reader := bufio.NewReader(os.Stdin)
	buf := make([]byte, 1024)

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		//客戶端指令
		if input[0] == '/' {
			switch input[1:] {
			case "exit":
				return
			case "?":
			case "help":
				log.Println(helpMsg)
			default:
				log.Println("Unknown command: ", input, " ,Please enter \"/help\" to see usage")
			}
		}

		_, _ = c.Write([]byte(input))

		cnt, err := c.Read(buf)
		if err != nil {
			log.Fatalln("Reading data failed\n", err)
		}
		log.Println(string(buf[0:cnt]))
	}
}

func main() {
	cfg := make(map[string]interface{}) //Configuration map
	b, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal([]byte(b), &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	//New a connect
	conn, err := net.Dial("tcp", cfg["server-addr"].(string)+":"+strconv.Itoa(cfg["server-port"].(int)))
	log.Println("Try to connect to the server: ", cfg["server-addr"].(string)+":"+strconv.Itoa(cfg["server-port"].(int)))
	if err != nil {
		log.Fatalln(err)
		return
	}
	handleConnection(conn)
}
