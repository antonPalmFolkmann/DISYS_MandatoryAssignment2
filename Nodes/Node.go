package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/antonPalmFolkmann/DISYS_MandatoryAssignment2/DistributedMutualExclusion"
)

type Node struct {
	DistributedMutualExclusion.UnimplementedCriticalSectionServiceServer
	ports      []string
	port       string
	leaderPort string
}

func main() {

}

func (n *Node) createNode(port string, listOfPorts []string) {
	//SERVER
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

	n.ports = listOfPorts

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

		// Send election request
		message := DistributedMutualExclusion.LeaderRequest{
			Port: n.port,
		}

		response, err := c.LeaderDeclaration(context.Background(), &message)
		if err != nil {
			log.Fatalf("Error when calling LeaderDeclaration: %s", err)
		}

		fmt.Printf("LeaderDeclaration response: %s\n", response.Reply)

	}
}
