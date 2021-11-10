package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/antonPalmFolkmann/DISYS_MandatoryAssignment2/DistributedMutualExclusion"
)

type Node struct {
	DistributedMutualExclusion.UnimplementedCriticalSectionServiceServer
	ports                 []string
	port                  string
	leaderPort            string
	queue                 []string //only for leader
	nodeInCriticalSection string   //only for leader
	isInCriticalSection   bool     //for node
}

func main() {

	fmt.Println("Main")

	ports := []string{"1111", "2222", "3333", "4444"}
	fmt.Println("ports created")

	var n1 Node
	var n2 Node
	var n3 Node
	var n4 Node

	fmt.Println("nodes declared")

	nodes := []Node{n1, n2, n3, n4}

	fmt.Println("nodes in array")

	for i, n := range nodes {
		go n.createNode(ports[i], ports)
	}

	fmt.Println("Nodes created")

	n1.startElection()

	go n1.sendQueueUpRequestAfter(100)
	go n2.sendQueueUpRequestAfter(323)
	go n3.sendQueueUpRequestAfter(187)
	go n4.sendQueueUpRequestAfter(245)

	writeStatus(nodes)

}
func writeStatus(nodes []Node) {
	for {
		time.Sleep(3 * time.Second)
		for _, n := range nodes {
			fmt.Printf("%v is in critial section: %v\n", n.port, n.isInCriticalSection)
		}
	}
}

func (n *Node) checkCriticalSection() {
	fmt.Printf("Checking crit: %v\n", n.nodeInCriticalSection)
	for {
		if strings.Compare(n.nodeInCriticalSection, "empty") == 0 {
			n.sendGrantAccessRequest()
		}
	}
}

func (n *Node) sendQueueUpRequestAfter(milliseconds int) {
	for {
		time.Sleep(time.Duration(milliseconds) * time.Millisecond)
		n.sendQueueUpRequest()
	}
}

func (n *Node) createNode(port string, listOfPorts []string) {
	//SERVER
	n.ports = listOfPorts
	n.port = port
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
		n.nodeInCriticalSection = "empty"
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
	go n.startElection()
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

		fmt.Printf("LeaderDeclaration response: %s\n", response.Reply)
	}
	fmt.Printf("checkCriticalSection go routine start")
	go n.checkCriticalSection()
	fmt.Printf("checkCriticalSection go routine started")

}

func (n *Node) LeaderDeclaration(ctx context.Context, in *DistributedMutualExclusion.LeaderRequest) (*DistributedMutualExclusion.LeaderReply, error) {
	n.leaderPort = in.Port

	return &DistributedMutualExclusion.LeaderReply{
		Reply: "OK",
	}, nil
}

func (n *Node) sendQueueUpRequest() {
	if n.leaderPort != "" {
		// Creat a virtual RPC Client Connection on leaderPort
		var conn *grpc.ClientConn
		portString := ":" + n.leaderPort
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

		fmt.Printf("QueueUp response: %s\n", response.Reply)
	}

}

func (n *Node) QueueUp(ctx context.Context, in *DistributedMutualExclusion.CriticalSectionRequest) (*DistributedMutualExclusion.CriticalSectionReply, error) {
	n.queue = append(n.queue, in.Port)

	return &DistributedMutualExclusion.CriticalSectionReply{
		Reply: "OK",
	}, nil

}

func (n *Node) sendGrantAccessRequest() {
	if len(n.queue) != 0 {
		// Creat a virtual RPC Client Connection on the node that can access
		var conn *grpc.ClientConn
		portString := ":" + n.queue[0] //next in queue
		n.nodeInCriticalSection = n.queue[0]
		fmt.Printf("queue: %v", n.queue)
		n.queue = dequeue(n.queue)

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

		fmt.Printf("Grant Access response: %s\n", response.Reply)
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
	go n.doingCriticalStuff()
	return &DistributedMutualExclusion.GrantAccessReply{
		Reply: "OK, access granted",
	}, nil
}

func (n *Node) doingCriticalStuff() {
	time.Sleep(4000 * time.Millisecond)
	n.sendLeaveRequest()
}

func (n *Node) sendLeaveRequest() {
	// Creat a virtual RPC Client Connection on leaderPort
	var conn *grpc.ClientConn
	portString := ":" + n.leaderPort
	conn, err := grpc.Dial(portString, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	c := DistributedMutualExclusion.NewCriticalSectionServiceClient(conn)

	// Send leave request
	message := DistributedMutualExclusion.LeaveRequest{
		Message: "I am leaving the critical section",
	}

	response, err := c.LeaveCriticalSection(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling leave critical section: %s", err)
	}

	n.isInCriticalSection = false

	fmt.Printf("Leave critical section response: %s\n", response.Reply)
}

func (n *Node) LeaveCriticalSection(ctx context.Context, in *DistributedMutualExclusion.LeaveRequest) (*DistributedMutualExclusion.LeaveReply, error) {
	n.nodeInCriticalSection = "empty"
	return &DistributedMutualExclusion.LeaveReply{
		Reply: "OK",
	}, nil
}
