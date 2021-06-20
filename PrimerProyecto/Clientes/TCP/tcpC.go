package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	CONNECT := "127.0.0.1:" + "2020"
	c, err := net.Dial("tcp4", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")
		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("Estado Del Contador->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP cliente saliendo...")
			return
		}
	}
}
