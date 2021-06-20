package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type Args struct {
	Action, Valor float64
}

func main() {
	var reply int64
	var info string

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")

	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text + "\n")
		instruction := strings.TrimSpace(string(data))
		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Println("Saliendo Cliente Remoto!")
			return
		}
		if instruction == "info" {
			err = client.Call("API.Information", 0, &info)
			if err != nil {
				log.Fatal("Connection error: ", err)
			}
			fmt.Println("Informacion: ", info)
		} else {
			split := strings.Split(instruction, ".")
			cadena := split[0] //inc.25=inc
			aux := split[1]    //inc.25=25
			val, err := strconv.Atoi(aux)
			action := 0

			if cadena == "inc" {
				action = 2
			}
			if cadena == "dec" {
				action = -1
			}
			if cadena == "res" {
				action = 8
			}

			a := Args{float64(action), float64(val)}

			err = client.Call("API.Operation", a, &reply)
			if err != nil {
				log.Fatal("Connection error: ", err)
			}

			fmt.Println("Estado Del Contador: ", reply)
		}

	}

}
