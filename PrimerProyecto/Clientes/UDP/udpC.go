package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	CONNECT := "127.0.0.1:" + "2002"

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Conectado a servidor UDP  %s\n", c.RemoteAddr().String())
	defer c.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text + "\n")
		_, err = c.Write(data)
		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Println("Saliendo UDP cliente!")
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		buffer := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Estado Del Contador: %s\n", string(buffer[0:n]))
	}
}
