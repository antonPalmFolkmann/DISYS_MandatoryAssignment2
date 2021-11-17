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

func main() {

	file, err := os.OpenFile("../info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)

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

		if strings.Compare("/q", input) == 0 {
			n.sendQueueUpRequest()
		} else {
			fmt.Println("Invalid message")
		}
	}

}

func (n *Node) addPort(port string) {
	n.ports = append(n.ports, port)
}

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

func (n *Node) sendQueueUpRequest() {
	if leaderPort != "" {
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

		response, err := c.QueueUp(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling critical section request: %s", err)
		}

		log.Printf("QueueUp response: %s\n", response.Reply)
	}

}

func (n *Node) QueueUp(ctx context.Context, in *DistributedMutualExclusion.CriticalSectionRequest) (*DistributedMutualExclusion.CriticalSectionReply, error) {
	queue = append(queue, in.Port)
	return &DistributedMutualExclusion.CriticalSectionReply{
		Reply: "OK",
	}, nil

}

func (n *Node) checkCriticalSection() {
	for {
		if strings.Compare(nodeInCriticalSection, "empty") == 0 {
			n.sendGrantAccessRequest()
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (n *Node) sendGrantAccessRequest() {
	if len(queue) != 0 {
		// Creat a virtual RPC Client Connection on the node that can access
		var conn *grpc.ClientConn
		port := queue[0]
		portString := ":" + queue[0] //next in queue
		nodeInCriticalSection = queue[0]
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
			Port:    port,
		}

		response, err := c.GrantAccess(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling grant access request: %s", err)
		}

		nodeInCriticalSection = response.Reply

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
	fmt.Println(in.Port + " is now in the critical section")
	log.Println(in.Port + " is now in the critical section")
	go n.doingCriticalStuff(in.Port)
	return &DistributedMutualExclusion.GrantAccessReply{
		Reply: n.port,
	}, nil
}

func (n *Node) doingCriticalStuff(port string) {
	time.Sleep(10000 * time.Millisecond)
	n.sendLeaveCriticalSectionRequest(port)
}

func (n *Node) sendLeaveCriticalSectionRequest(port string) {
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
		Message: n.port + " leaving the critical section",
	}

	response, err := c.LeaveCriticalSection(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling leave critical section: %s", err)
	}

	n.isInCriticalSection = false

	fmt.Printf(port+" is now leaving: \nLeave critical section response: %s\n", response.Reply)
	log.Printf(port+" is now leaving: \nLeave critical section response: %s\n", response.Reply)
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

	log.Printf("Join response: %s\n", response.Ports)
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

		log.Printf("update ports request response: %s\n", response.Reply)
	}
}

func (n *Node) UpdatePorts(ctx context.Context, in *DistributedMutualExclusion.UpdatePortsRequest) (*DistributedMutualExclusion.UpdatePortsReply, error) {
	n.ports = strings.Split(in.Ports, " ")

	return &DistributedMutualExclusion.UpdatePortsReply{
		Reply: "OK",
	}, nil
}

//code for starting election, which we did not end up using
func (n *Node) startElection() {
	sent := false
	recievedResponse := false
	for _, port := range n.ports {
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

	log.Println("Election reply: ", response.Reply)
	return true
}

func (n *Node) Election(ctx context.Context, in *DistributedMutualExclusion.ElectionRequest) (*DistributedMutualExclusion.ElectionReply, error) {
	n.startElection()
	return &DistributedMutualExclusion.ElectionReply{
		Reply: "OK",
	}, nil
}

func (n *Node) sendLeaderRequest() {
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
		message := DistributedMutualExclusion.LeaderRequest{
			Port: n.port,
		}

		response, err := c.LeaderDeclaration(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling LeaderDeclaration: %s", err)
		}
		log.Printf("Leader Declaration response: %v", response.Reply)

	}
	go n.checkCriticalSection()

}

func (n *Node) LeaderDeclaration(ctx context.Context, in *DistributedMutualExclusion.LeaderRequest) (*DistributedMutualExclusion.LeaderReply, error) {
	leaderPort = in.Port

	return &DistributedMutualExclusion.LeaderReply{
		Reply: "OK",
	}, nil
}
