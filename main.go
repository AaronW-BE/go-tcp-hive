package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net"
	"strings"
	"time"
)

func pass(err error) {
	if err != nil {
		panic(err.Error())
	}
}

/*
protocol:
	cmd:
		- 0x01 for init feed connect
*/
type protocol struct {
	length int32
	cmd    byte
}

type Node struct {
	listener *net.Listener
	conn     net.Conn
	addr     net.Addr
	uid      string
	birthAt  time.Time
	hive     *NodeHive
}

func (node *Node) Init(port int) {
	listener, err := net.Listen("tcp", ":"+string(rune(port)))
	pass(err)

	node.listener = &listener
	node.uid = uuid.NewV4().String()
}

func (node *Node) Read() {

}

func (node *Node) Write() {

}

func (node *Node) Feed(addr string, port int) {
	builder := strings.Builder{}
	builder.WriteString(addr)
	builder.WriteRune(rune(port))

	conn, err := net.Dial("tcp", builder.String())
	pass(err)

	fmt.Printf("Feed remote node %s", conn.RemoteAddr().String())
}

func (node *Node) Run() {

}

func NewNode() *Node {
	node := &Node{}
	return node
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
			conn: conn,
			uid:  uuid.NewV4().String(),
			addr: conn.RemoteAddr(),
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
