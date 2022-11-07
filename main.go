package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net"
	"time"
)

func pass(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Node struct {
	id          int
	conn        net.Conn
	addr        net.Addr
	uid         string
	connectTime time.Time
}

func (node *Node) Init() {
}

func NewNode() {

}

type NodeHive struct {
	nodes    []Node
	size     int
	lastJoin time.Time
}

var nodeHive = make([]Node, 100)

func main() {
	server, err := net.Listen("tcp", "127.0.0.1:6789")
	pass(err)

	fmt.Printf("server running at %v...\n", 6789)
	//
	defer func(conn net.Listener) {
		err := conn.Close()
		pass(err)
	}(server)

	for {
		conn, err := server.Accept()
		pass(err)

		client := Node{
			id:          len(nodeHive) + 1,
			conn:        conn,
			uid:         uuid.NewV4().String(),
			addr:        conn.RemoteAddr(),
			connectTime: time.Time{},
		}
		nodeHive = append(nodeHive, client)

		for i := range nodeHive {
			go func() {
				for {
					buf := make([]byte, 1024)

					n, readErr := nodeHive[i].conn.Read(buf)

					if readErr != nil {
						for ci := range nodeHive {
							if nodeHive[ci].addr == nodeHive[i].addr {
								println("remove disconnected client", nodeHive[i].uid)
								nodeHive = append(nodeHive[:ci], nodeHive[ci+1:]...)
								return
							}
						}
					}

					if n > 0 {
						fmt.Printf("[SYS] RECV Data: %s \n", buf)
					}
				}
			}()

			i := i
			go func() {
				for {
					input := make([]byte, 1024)
					n, err2 := fmt.Scan(&input)
					pass(err2)
					if n > 0 {
						fmt.Printf("[SYS] Send msg: %s\n", input)
						_, _ = nodeHive[i].conn.Write(input)
					}
				}
			}()
		}

	}
}
