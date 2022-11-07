package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:6789")
	pass(err)

	go func() {
		for true {
			input := make([]byte, 1024)
			fmt.Scan(&input)
			_, err = conn.Write(input)
			pass(err)
			fmt.Printf("[SYS] Send msg: %s\n", input)
		}
	}()

	for {
		output := make([]byte, 1024)
		n, _ := conn.Read(output)
		if n > 0 {
			fmt.Printf("[SYS] RECV msg: %s\n", output)
		}
	}
}

func pass(err error) {
	if err != nil {
		panic(err.Error())
	}
}
