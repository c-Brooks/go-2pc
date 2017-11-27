// Package rpc contains stubs and implementations of
// funcs defined in rpc.proto
package rpc

import context "golang.org/x/net/context"

// Server implements the RPCServiceServer interface defined in rpc.pb.go
// RequestTransactionTry(RPCService_RequestTransactionTryServer) error
// RequestTransactionCommit(context.Context, *CommitReq) (*CommitResp, error)
type Server struct{}

// RequestTransactionTry pushes a transactionTry request from the master to the slave replicas
func (s Server) RequestTransactionTry(RPCService_RequestTransactionTryServer) error {
	return nil
}

// RequestTransactionCommit pushes a commitRequest from master to slave replicas
func (s Server) RequestTransactionCommit(context.Context, *CommitReq) (*CommitResp, error) {
	return nil, nil
}
