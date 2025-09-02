package main

import (
	"HSECourse/Homework3/models"
	"HSECourse/Homework3/proto"
	"context"
	"errors"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type bankServer struct {
	proto.BankProtoServer
	accounts map[string]*models.Account
	guard    *sync.RWMutex
}

func newBankServer() *bankServer {
	return &bankServer{
		accounts: make(map[string]*models.Account),
		guard:    &sync.RWMutex{},
	}
}

func (s *bankServer) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	if len(req.GetName()) == 0 {
		return nil, errors.New("name is empty")
	}

	s.guard.Lock()
	if _, ok := s.accounts[req.GetName()]; ok {
		s.guard.Unlock()

		return nil, errors.New("account already exists")
	}

	s.accounts[req.GetName()] = &models.Account{
		Name:   req.GetName(),
		Amount: int(req.GetAmount()),
	}

	s.guard.Unlock()

	response := &proto.CreateAccountResponse{Result: "account created"}

	return response, nil
}

func (s *bankServer) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error) {
	if len(req.GetName()) == 0 {
		return nil, errors.New("name is empty")
	}

	s.guard.Lock()

	if _, ok := s.accounts[req.GetName()]; !ok {
		s.guard.Unlock()

		return nil, errors.New("account not found")
	}

	delete(s.accounts, req.GetName())
	s.guard.Unlock()

	response := &proto.DeleteAccountResponse{Result: "account deleted"}

	return response, nil
}

func (s *bankServer) ChangeAccountName(ctx context.Context, req *proto.ChangeAccountNameRequest) (*proto.ChangeAccountNameResponse, error) {
	if len(req.GetName()) == 0 {
		return nil, errors.New("name is empty")
	}

	s.guard.Lock()

	if _, ok := s.accounts[req.GetName()]; !ok {
		s.guard.Unlock()

		return nil, errors.New("account not found")
	}

	s.accounts[req.GetNewName()] = &models.Account{
		Name:   req.GetNewName(),
		Amount: s.accounts[req.GetName()].Amount,
	}
	delete(s.accounts, req.GetName())

	s.guard.Unlock()

	response := &proto.ChangeAccountNameResponse{Result: "account name changed"}

	return response, nil
}

func (s *bankServer) ChangeAccountAmount(ctx context.Context, req *proto.ChangeAccountAmountRequest) (*proto.ChangeAccountAmountResponse, error) {
	if len(req.GetName()) == 0 {
		return nil, errors.New("name is empty")
	}

	s.guard.Lock()

	if _, ok := s.accounts[req.GetName()]; !ok {
		s.guard.Unlock()

		return nil, errors.New("account not found")
	}

	s.accounts[req.GetName()] = &models.Account{
		Name:   req.GetName(),
		Amount: int(req.GetNewAmount()),
	}

	s.guard.Unlock()

	response := &proto.ChangeAccountAmountResponse{Result: "account amount changed"}

	return response, nil
}

func (s *bankServer) GetAccount(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	if len(req.GetName()) == 0 {
		return nil, errors.New("name is empty")
	}

	s.guard.RLock()
	account, ok := s.accounts[req.GetName()]
	s.guard.RUnlock()

	if !ok {
		return nil, errors.New("account not found")
	}

	response := &proto.GetAccountResponse{
		Name:   account.Name,
		Amount: int64(account.Amount),
	}

	return response, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	proto.RegisterBankProtoServer(s, newBankServer())
	if err := s.Serve(listener); err != nil {
		panic(err)
	}
}
