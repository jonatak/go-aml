package aml

import (
	"context"
	"fmt"
	"jonatak/aml/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	proto.UnimplementedPaymentsServer
	tx chan TransactionQuery
}

func (s *GRPCServer) ApproveTransaction(ctx context.Context, request *proto.TransactionRequest) (*proto.TransactionResponse, error) {
	rspChan := make(chan proto.TransactionStatus, 1)

	go func() {
		s.tx <- TransactionQuery{
			Tx:           request.Transaction,
			ResponseChan: rspChan,
		}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context was cancelled: %w", ctx.Err())
	case rsp := <-rspChan:
		return &proto.TransactionResponse{
			Status: rsp,
		}, nil
	}
}

func StartServer(ctx context.Context, port int, tx chan TransactionQuery) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterPaymentsServer(grpcServer, &GRPCServer{
		tx: tx,
	})

	errChan := make(chan error)

	go func() {
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			errChan <- fmt.Errorf("gRPC server failed: %w", err)
		}
	}()

	fmt.Printf("gRPC server started on port: %d\n", port)

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
	case err := <-errChan:
		return err
	}

	return nil
}
