package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/antonPalmFolkmann/DISYS_MandatoryAssignment2/DistributedMutualExclusion"
)

type Node struct {
	DistributedMutualExclusion.UnimplementedCriticalSectionServiceServer
	ports               []string
	port                string
	isInCriticalSection bool //for node
}

var queue []string //only for leader
var leaderPort string
var nodeInCriticalSection string //only for leader

/*func main() {

	fmt.Println("Main")

	ports := []string{"1111", "2223", "3333", "4444"}
	fmt.Println("ports created")

	nodes := make([]*Node, 0)
	fmt.Println("nodes created")

	for _, p := range ports {
		nodes = append(nodes, setNodeValues(p, ports))
	}

	fmt.Println("Nodes created")

	nodes[0].startElection()

	for i := 0; i < len(nodes); i++ {
		fmt.Printf("i: %v, node: %v, nodes leaderport: %v <---\n", i, nodes[i].port, nodes[i].leaderPort)
		go nodes[i].sendQueueUpRequestAfter(56)
	}

	writeStatus(nodes)

}*/

func main() {
	// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please write your port")
	input, _ := reader.ReadString('\n')
	// convert CRLF to LF
	input = strings.Replace(input, "\n", "", -1)

	var n Node
	n.port = input
	n.ports = make([]string, 0)
	n.isInCriticalSection = false

	//portString := ":" + input
	// Create listener tcp on port *input*
	go n.startListen(n.port, n.ports)

	fmt.Println("Please write the leaders port")
	input2, _ := reader.ReadString('\n')
	// convert CRLF to LF
	input2 = strings.Replace(input2, "\n", "", -1)
	leaderPort = input2

	n.sendJoinRequest()

	if n.port == leaderPort {
		nodeInCriticalSection = "empty"
		go n.checkCriticalSection()
	}

	for {

		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)

		if strings.Compare("/add", input) == 0 {
			fmt.Println("Write the port:")
			portInput, _ := reader.ReadString('\n')
			// convert CRLF to LF
			portInput = strings.Replace(portInput, "\n", "", -1)
			portInput = strings.Replace(portInput, "\r", "", -1)

			n.addPort(portInput)

		} else if strings.Compare("/q", input) == 0 {
			n.sendQueueUpRequest()
		} else {
			fmt.Println("Invalid message")
		}
	}

}

func (n *Node) addPort(port string) {
	n.ports = append(n.ports, port)
}

/*func (n *Node) fillPorts()  {
	n.addPort("1111")
	n.addPort("2222")
	n.addPort("3333")
	n.addPort("4444")
}*/

/*func writeStatus(nodes []*Node) {
	for {
		time.Sleep(3 * time.Second)
		for _, n := range nodes {
			fmt.Printf("%v is in critial section: %v\n", n, n.isInCriticalSection)
		}
	}
}*/

/*func setNodeValues(port string, listOfPorts []string) *Node {

	var n Node
	fmt.Printf("Set node values before: %v\n", n)
	n.port = port
	n.ports = listOfPorts
	n.isInCriticalSection = false

	fmt.Printf("Set node values after: %v\n", n)

	fmt.Printf("start goroutine\n")
	go n.startListen(port, listOfPorts)
	fmt.Printf("goroutine started\n")
	return &n
}*/

func (n *Node) startListen(port string, listOfPorts []string) {
	//SERVER
	portString := ":" + port
	// Create listener tcp on portString
	list, err := net.Listen("tcp", portString)
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", port, err)
	}
	grpcServer := grpc.NewServer()
	DistributedMutualExclusion.RegisterCriticalSectionServiceServer(grpcServer, &Node{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}

func (n *Node) startElection() {
	sent := false
	recievedResponse := false
	fmt.Printf("election node: %v\n", n)
	for _, port := range n.ports {
		fmt.Println("Inside startElection loop")
		if port > n.port {
			sent = true
			// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
			var conn *grpc.ClientConn
			portString := ":" + port
			conn, err := grpc.Dial(portString, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}

			// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
			defer conn.Close()

			//  Create new Client from generated gRPC code from proto
			c := DistributedMutualExclusion.NewCriticalSectionServiceClient(conn)

			// Send election request
			if sendElectionRequest(c) {
				recievedResponse = true
			}
		}
	}
	if !sent || !recievedResponse {
		nodeInCriticalSection = "empty"
		n.sendLeaderRequest()
		fmt.Println("Leader Node chosen")
	}
}

func sendElectionRequest(c DistributedMutualExclusion.CriticalSectionServiceClient) bool {
	message := DistributedMutualExclusion.ElectionRequest{
		Message: "Election",
	}
	response, err := c.Election(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Elction: %s", err)
		return false
	}
	if response == nil {
		log.Printf("Response was nil")
		return false
	}

	log.Printf("Election reply: ", response.Reply)
	return true
}

func (n *Node) Election(ctx context.Context, in *DistributedMutualExclusion.ElectionRequest) (*DistributedMutualExclusion.ElectionReply, error) {
	n.startElection()
	return &DistributedMutualExclusion.ElectionReply{
		Reply: "OK",
	}, nil
}

func (n *Node) sendLeaderRequest() {
	fmt.Printf("SendLeaderRequest before loop, ports: %v, port: %v\n", n.ports, n.port)
	for _, port := range n.ports {
		fmt.Printf("SendLeaderRequest inside loop, port to: %v\n", port)

		// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
		var conn *grpc.ClientConn
		portString := ":" + port
		conn, err := grpc.Dial(portString, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}

		// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
		defer conn.Close()

		//  Create new Client from generated gRPC code from proto
		c := DistributedMutualExclusion.NewCriticalSectionServiceClient(conn)

		// Send leader request
		message := DistributedMutualExclusion.LeaderRequest{
			Port: n.port,
		}

		response, err := c.LeaderDeclaration(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling LeaderDeclaration: %s", err)
		}

		fmt.Printf("LeaderDeclaration response: %s\n", response.Reply)
	}
	fmt.Printf("checkCriticalSection go routine start\n")
	go n.checkCriticalSection()
	fmt.Printf("checkCriticalSection go routine started\n")

}

func (n *Node) LeaderDeclaration(ctx context.Context, in *DistributedMutualExclusion.LeaderRequest) (*DistributedMutualExclusion.LeaderReply, error) {
	fmt.Printf("before leader declaration, port: %v, leaderport: %v\n", n.port, leaderPort)
	leaderPort = in.Port
	fmt.Printf("after leader declaration, port: %v, leaderport: %v\n", n.port, leaderPort)

	return &DistributedMutualExclusion.LeaderReply{
		Reply: "OK",
	}, nil
}

func (n *Node) sendQueueUpRequest() {
	fmt.Printf("%v send queue up request start\n", n.port)
	if leaderPort != "" {
		fmt.Printf("send queue up request inside if\n")
		// Creat a virtual RPC Client Connection on leaderPort
		var conn *grpc.ClientConn
		portString := ":" + leaderPort
		conn, err := grpc.Dial(portString, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}

		// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
		defer conn.Close()

		//  Create new Client from generated gRPC code from proto
		c := DistributedMutualExclusion.NewCriticalSectionServiceClient(conn)

		// Send critical section request
		message := DistributedMutualExclusion.CriticalSectionRequest{
			Port: n.port,
		}

		fmt.Printf("sned queue: Leaderport: %v", leaderPort)
		response, err := c.QueueUp(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling critical section request: %s", err)
		}

		fmt.Printf("QueueUp response: %s\n", response.Reply)
	}

}

func (n *Node) QueueUp(ctx context.Context, in *DistributedMutualExclusion.CriticalSectionRequest) (*DistributedMutualExclusion.CriticalSectionReply, error) {
	queue = append(queue, in.Port)
	fmt.Printf("n.queue: %v\n", queue)
	return &DistributedMutualExclusion.CriticalSectionReply{
		Reply: "OK",
	}, nil

}

func (n *Node) checkCriticalSection() {
	fmt.Printf("port: %v, queue: %v\n", n.port, queue)

	//i := 0

	for {
		/*fmt.Printf("i: %v, node in CS: %v\n", i, n.nodeInCriticalSection)
		i++*/

		if strings.Compare(nodeInCriticalSection, "empty") == 0 {
			fmt.Printf("send Grant access request: port: %v, queue: %v\n", n.port, queue)
			n.sendGrantAccessRequest()
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (n *Node) sendQueueUpRequestAfter(milliseconds int) {
	for {
		time.Sleep(10 * time.Second)
		n.sendQueueUpRequest()
		time.Sleep(10 * time.Second)
		//time.Sleep(time.Duration(milliseconds) * time.Millisecond)
	}
}

func (n *Node) sendGrantAccessRequest() {
	if len(queue) != 0 {
		// Creat a virtual RPC Client Connection on the node that can access
		fmt.Println("sendGrantAccessRequest")
		var conn *grpc.ClientConn
		portString := ":" + queue[0] //next in queue
		nodeInCriticalSection = queue[0]
		fmt.Printf("queue: %v\n", queue)
		queue = dequeue(queue)

		conn, err := grpc.Dial(portString, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}

		// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
		defer conn.Close()

		//  Create new Client from generated gRPC code from proto
		c := DistributedMutualExclusion.NewCriticalSectionServiceClient(conn)

		// Send grant access request
		message := DistributedMutualExclusion.GrantAccessRequest{
			Message: "You may now enter the critical section",
		}

		response, err := c.GrantAccess(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling grant access request: %s", err)
		}

		nodeInCriticalSection = response.Reply

		fmt.Printf("Grant Access response: %s\n", nodeInCriticalSection)
	}
}

func dequeue(queue []string) []string {
	newQueue := make([]string, 0)

	for i, n := range queue {
		if i != 0 {
			newQueue = append(newQueue, n)
		}
	}

	return newQueue
}

func (n *Node) GrantAccess(ctx context.Context, in *DistributedMutualExclusion.GrantAccessRequest) (*DistributedMutualExclusion.GrantAccessReply, error) {
	n.isInCriticalSection = true
	fmt.Printf("We did it!")
	go n.doingCriticalStuff()
	return &DistributedMutualExclusion.GrantAccessReply{
		Reply: n.port,
	}, nil
}

func (n *Node) doingCriticalStuff() {
	time.Sleep(10000 * time.Millisecond)
	n.sendLeaveCriticalSectionRequest()
}

func (n *Node) sendLeaveCriticalSectionRequest() {
	// Creat a virtual RPC Client Connection on leaderPort
	var conn *grpc.ClientConn
	portString := ":" + leaderPort
	conn, err := grpc.Dial(portString, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	c := DistributedMutualExclusion.NewCriticalSectionServiceClient(conn)

	// Send leave request
	message := DistributedMutualExclusion.LeaveCriticalSectionRequest{
		Message: "I am leaving the critical section",
	}

	response, err := c.LeaveCriticalSection(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling leave critical section: %s", err)
	}

	n.isInCriticalSection = false

	fmt.Printf("Leave critical section response: %s\n", response.Reply)
}

func (n *Node) LeaveCriticalSection(ctx context.Context, in *DistributedMutualExclusion.LeaveCriticalSectionRequest) (*DistributedMutualExclusion.LeaveCriticalSectionReply, error) {
	nodeInCriticalSection = "empty"
	return &DistributedMutualExclusion.LeaveCriticalSectionReply{
		Reply: "OK",
	}, nil
}

func (n *Node) sendJoinRequest() {
	// Creat a virtual RPC Client Connection on leaderPort
	var conn *grpc.ClientConn
	portString := ":" + leaderPort
	conn, err := grpc.Dial(portString, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	c := DistributedMutualExclusion.NewCriticalSectionServiceClient(conn)

	// Send leave request
	message := DistributedMutualExclusion.JoinRequest{
		Port: n.port,
	}

	response, err := c.Join(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling join: %s", err)
	}

	n.ports = strings.Split(response.Ports, " ")

	fmt.Printf("Join response: %s\n", response.Ports)
}

func (n *Node) Join(ctx context.Context, in *DistributedMutualExclusion.JoinRequest) (*DistributedMutualExclusion.JoinReply, error) {
	var portsstring string

	for _, p := range n.ports {
		portsstring = portsstring + p + " "
	}

	portsstring = portsstring + in.Port

	n.ports = append(n.ports, in.Port)

	n.sendUpdatePortsRequest(portsstring)

	return &DistributedMutualExclusion.JoinReply{
		Ports: portsstring,
	}, nil
}

func (n *Node) sendUpdatePortsRequest(ports string) { //called on leader
	for _, port := range n.ports {

		// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
		var conn *grpc.ClientConn
		portString := ":" + port
		conn, err := grpc.Dial(portString, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}

		// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
		defer conn.Close()

		//  Create new Client from generated gRPC code from proto
		c := DistributedMutualExclusion.NewCriticalSectionServiceClient(conn)

		// Send leader request
		message := DistributedMutualExclusion.UpdatePortsRequest{
			Ports: ports,
		}

		response, err := c.UpdatePorts(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling send update ports request: %s", err)
		}

		fmt.Printf("update ports request response: %s\n", response.Reply)
	}
}

func (n *Node) UpdatePorts(ctx context.Context, in *DistributedMutualExclusion.UpdatePortsRequest) (*DistributedMutualExclusion.UpdatePortsReply, error) {
	n.ports = strings.Split(in.Ports, " ")

	return &DistributedMutualExclusion.UpdatePortsReply{
		Reply: "OK",
	}, nil
}
