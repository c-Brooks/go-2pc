package cluster

import (
	"google.golang.org/grpc"
)

// Master is the cluster manager
type Master struct {
	grpc.Server
}

// AssembleQuorum assembles a quorum of nodes for a read or write op
func (m *Master) AssembleQuorum() {

}

// RequestTransactionStatus is a TODO
func (m *Master) RequestTransactionStatus() {
	// Loop thru clients
	// Push a XactReq to each client
	// Wait for confirmation from quorum
}

// GRPC SERVER
// Slaves connect to this server.
// Transaction requests are pushed to the clients
// and commit/abort messages are returned
// func (m *Master) startGrpcServer(port int) *grpc.Server {
// 	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
// 	if err != nil {
// 		logrus.Fatalf("master could not start gRPC server: %v", err)
// 	}

// 	grpcServer := grpc.NewServer()
// 	s := pb.Server{}
// 	pb.RegisterRPCServiceServer(grpcServer, s)

// 	grpcServer.Serve(conn)
// 	return grpcServer
// }
