package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

//GCD for 2 numbers
func GCD(a int, b int) int {
	var GCD int
	for i := 1; i <= a && i <= b; i++ {
		if a%i == 0 && b%i == 0 {
			GCD = i
		}
	}
	return GCD
}

//previous function applied to an array
func GCDarr(arr []int) int {
	var GCDVar int = arr[0]
	for i := 1; i < len(arr); i++ {
		var elem int = arr[i]
		GCDVar = GCD(elem, GCDVar)
	}
	return GCDVar
}

func check(err error, message string) {
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", message)
}

func main() {
	//reading the condtion for the minimum integer elements from a file
	data, err := ioutil.ReadFile("config.txt")
	if err != nil {
		fmt.Println("ERROR! File reading.", err)
		return
	}
	length := string(data)
	lengthInt, _ := strconv.Atoi(length)

	clientCount := 0
	allClients := make(map[net.Conn]int)

	ln, err := net.Listen("tcp", ":8080")
	check(err, "Server is ready.")

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		allClients[conn] = clientCount
		fmt.Printf("Client %d Connected.\n", allClients[conn])

		clientCount += 1

		go func() {
			reader := bufio.NewReader(conn)

			for {
			tag:
				incoming, err := reader.ReadString('\n')

				if err != nil {
					fmt.Printf("Client disconnected.\n")
					break
				}

				fmt.Printf("Client %d requested: %s", allClients[conn], incoming)
				conn.Write([]byte("Server received request.\n"))

				incoming = incoming[0 : len(incoming)-2]

				//sepparate the elements by space
				v := strings.Split(incoming, " ")
				var arr []int
				for i := 0; i < len(v); i++ {
					if _, err := strconv.Atoi(v[i]); err == nil {
						fmt.Printf("%q looks like a number.\n", v[i])
						elem, _ := strconv.Atoi(v[i])
						arr = append(arr, elem)
					}
				}

				if len(arr) < lengthInt {
					l := strconv.Itoa(lengthInt)
					conn.Write([]byte("ERROR! Please enter an array with a minimum of " + l + " integers, separated by space.\n"))
					goto tag
				}
				conn.Write([]byte("Server is processing data\n"))
				result := GCDarr(arr[:])
				GCD := strconv.Itoa(result)

				nume := strconv.Itoa(allClients[conn])
				fmt.Printf("Server sends " + incoming + " => " + GCD + "\n")
				conn.Write([]byte("Client " + nume + " recevies: " + incoming + " =>" + GCD + "\n"))

			}
		}()
	}
}
