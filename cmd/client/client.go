package main

import (
	"context"
	"flag"
	"fmt"
	"jonatak/aml/proto"
	"log"
	"math/rand/v2"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func parse(datetime string) time.Time {
	val, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		log.Fatalf("Invalid datetime in test %v.", err)
	}
	return val
}

type TxWrapper struct {
	Expected []proto.TransactionStatus
	Tx       []proto.Transaction
}

func getTx() TxWrapper {
	return TxWrapper{
		Tx: []proto.Transaction{
			{Id: "1", UserId: "u1", Amount: 5000, Timestamp: timestamppb.New(parse("2025-07-01T01:00:00Z"))},
			{Id: "2", UserId: "u1", Amount: 4000, Timestamp: timestamppb.New(parse("2025-07-02T00:30:00Z"))},
			{Id: "3", UserId: "u1", Amount: 4000, Timestamp: timestamppb.New(parse("2025-07-03T00:00:00Z"))},
			{Id: "4", UserId: "u1", Amount: 7000, Timestamp: timestamppb.New(parse("2025-07-03T12:00:00Z"))},
			{Id: "5", UserId: "u1", Amount: 7000, Timestamp: timestamppb.New(parse("2025-07-01T23:59:00Z"))},
		},
		Expected: []proto.TransactionStatus{
			proto.TransactionStatus_APPROVED,
			proto.TransactionStatus_APPROVED,
			proto.TransactionStatus_APPROVED,
			proto.TransactionStatus_MAX_AMOUNT_REACH,
			proto.TransactionStatus_INVALID_DATE,
		},
	}
}

func main() {
	port := flag.Int("p", 8080, "GRpc port to connect to.")
	flag.Parse()

	conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%d", *port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("Error connecting to gRPC server: ", err.Error())
	}

	defer conn.Close()

	// create the stream
	client := proto.NewPaymentsClient(conn)
	txs := getTx()

	for i := range txs.Tx {
		rsp, err := client.ApproveTransaction(context.Background(), &proto.TransactionRequest{
			Transaction: &txs.Tx[i],
		})
		if err != nil {
			log.Fatalf("some error: %v", err)
		}
		if rsp.Status != txs.Expected[i] {
			fmt.Printf("received unexpected status for tx id: %s\n", txs.Tx[i].GetId())
		} else {
			fmt.Printf("worked for tx id: %s\n", txs.Tx[i].GetId())
		}
	}

	i := 0

	startTime := time.Now().AddDate(-2, 0, 0)
	for {

		tx := &proto.Transaction{
			Id:        fmt.Sprintf("%d", i),
			UserId:    fmt.Sprintf("s%d", rand.IntN(1000)),
			Amount:    rand.Float64() * (10000.0),
			Timestamp: timestamppb.New(startTime.Add(time.Duration(i) * time.Hour)),
		}
		rsp, err := client.ApproveTransaction(context.Background(), &proto.TransactionRequest{
			Transaction: tx,
		})
		i++
		if err != nil {
			log.Fatalf("some error: %v", err)
		}
		fmt.Printf("response: %v\n", rsp.Status)
	}

}
