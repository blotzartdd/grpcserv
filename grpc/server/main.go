package main

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"grpcserv/database"
	"grpcserv/proto"
	"log"
	"net"
	"os"
)

type bankServer struct {
	proto.BankProtoServer
	db *database.DB
}

func newBankServer(db *database.DB) *bankServer {
	return &bankServer{
		db: db,
	}
}

func (s *bankServer) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	if len(req.GetName()) == 0 {
		return nil, errors.New("name is empty")
	}

	err := s.db.CreateAccount(req.GetName(), req.GetAmount())
	if err != nil {
		return nil, err
	}

	response := &proto.CreateAccountResponse{Result: "account created"}
	return response, nil
}

func (s *bankServer) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error) {
	if len(req.GetName()) == 0 {
		return nil, errors.New("name is empty")
	}

	err := s.db.DeleteAccount(req.GetName())
	if err != nil {
		return nil, err
	}

	response := &proto.DeleteAccountResponse{Result: "account deleted"}
	return response, nil
}

func (s *bankServer) ChangeAccountName(ctx context.Context, req *proto.ChangeAccountNameRequest) (*proto.ChangeAccountNameResponse, error) {
	if len(req.GetName()) == 0 {
		return nil, errors.New("name is empty")
	}

	if len(req.GetNewName()) == 0 {
		return nil, errors.New("new name is empty")
	}

	err := s.db.ChangeAccountName(req.GetName(), req.GetNewName())
	if err != nil {
		return nil, err
	}

	response := &proto.ChangeAccountNameResponse{Result: "account name changed"}
	return response, nil
}

func (s *bankServer) ChangeAccountAmount(ctx context.Context, req *proto.ChangeAccountAmountRequest) (*proto.ChangeAccountAmountResponse, error) {
	if len(req.GetName()) == 0 {
		return nil, errors.New("name is empty")
	}

	err := s.db.ChangeAccountAmount(req.GetName(), req.GetNewAmount())
	if err != nil {
		return nil, err
	}

	response := &proto.ChangeAccountAmountResponse{Result: "account amount changed"}
	return response, nil
}

func (s *bankServer) GetAccount(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	if len(req.GetName()) == 0 {
		return nil, errors.New("name is empty")
	}

	account, err := s.db.GetAccount(req.GetName())
	if err != nil {
		return nil, err
	}

	response := &proto.GetAccountResponse{
		Name:   account.Name,
		Amount: int64(account.Amount),
	}

	return response, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	dbHost := getEnv("DB_HOST", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "bankdb")

	db, err := database.NewDB(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterBankProtoServer(s, newBankServer(db))

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
