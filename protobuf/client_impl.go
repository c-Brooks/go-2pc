// Package rpc contains stubs and implementations of
// funcs defined in rpc.proto
package rpc

import (
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Client implements the RPCServiceClient interface defined in rpc.pb.go
// RequestTransactionTry(ctx context.Context, opts ...grpc.CallOption) (RPCService_RequestTransactionTryClient, error)
// RequestTransactionCommit(ctx context.Context, in *CommitReq, opts ...grpc.CallOption) (*CommitResp, error)
type Client struct{}

// RequestTransactionTry pushes a transactionTry request from the master to the slave replicas
func (c *Client) RequestTransactionTry(ctx context.Context, opts ...grpc.CallOption) (RPCService_RequestTransactionTryClient, error) {
	return nil, nil
}

// RequestTransactionCommit pushes a commitRequest from master to slave replicas
func (c *Client) RequestTransactionCommit(ctx context.Context, in *CommitReq, opts ...grpc.CallOption) (*CommitResp, error) {
	return nil, nil
}

// HealthCheck reports status of node
// TODO: loop over clients and ping each one
func (c *Client) HealthCheck(ctx context.Context, in *HealthCheckReq, opts ...grpc.CallOption) (*HealthCheckResp, error) {
	hcr := HealthCheckResp{
		Timestamp:  &google_protobuf.Timestamp{},
		StatusCode: 200,
	}

	return &hcr, nil
}
