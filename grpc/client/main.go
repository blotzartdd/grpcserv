package main

import (
	"HSECourse/Homework3/proto"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Command struct {
	Port      int
	Host      string
	Cmd       string
	Name      string
	NewName   string
	Amount    int
	NewAmount int
}

func (c *Command) do() error {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", c.Host, c.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = conn.Close()
	}()

	switch c.Cmd {
	case "create":
		return c.create(conn)
	case "delete":
		return c.delete(conn)
	case "changeName":
		return c.changeName(conn)
	case "changeAmount":
		return c.changeAmount(conn)
	case "get":
		return c.get(conn)
	default:
		return fmt.Errorf("unknown command: %s", c.Cmd)
	}
}

func (c *Command) create(conn *grpc.ClientConn) error {
	client := proto.NewBankProtoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := &proto.CreateAccountRequest{
		Name:   c.Name,
		Amount: int64(c.Amount),
	}

	_, err := client.CreateAccount(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (c *Command) delete(conn *grpc.ClientConn) error {
	client := proto.NewBankProtoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := &proto.DeleteAccountRequest{
		Name: c.Name,
	}

	_, err := client.DeleteAccount(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (c *Command) changeName(conn *grpc.ClientConn) error {
	client := proto.NewBankProtoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := &proto.ChangeAccountNameRequest{
		Name:    c.Name,
		NewName: c.NewName,
	}

	_, err := client.ChangeAccountName(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (c *Command) changeAmount(conn *grpc.ClientConn) error {
	client := proto.NewBankProtoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := &proto.ChangeAccountAmountRequest{
		Name:      c.Name,
		NewAmount: int64(c.NewAmount),
	}

	_, err := client.ChangeAccountAmount(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (c *Command) get(conn *grpc.ClientConn) error {
	client := proto.NewBankProtoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request := &proto.GetAccountRequest{
		Name: c.Name,
	}

	response, err := client.GetAccount(ctx, request)
	if err != nil {
		return err
	}

	fmt.Printf("Account name: %s, amount: %d\n", response.Name, response.Amount)
	return nil
}

func main() {
	portVal := flag.Int("port", 8080, "server port")
	hostVal := flag.String("host", "0.0.0.0", "server host")
	cmdVal := flag.String("cmd", "", "command to execute")
	nameVal := flag.String("name", "", "name of account")
	newNameVal := flag.String("newName", "", "new name of account")
	amountVal := flag.Int("amount", 0, "amount of account")
	newAmountVal := flag.Int("newAmount", 0, "new amount of account")

	flag.Parse()

	cmd := Command{
		Port:      *portVal,
		Host:      *hostVal,
		Cmd:       *cmdVal,
		Name:      *nameVal,
		NewName:   *newNameVal,
		Amount:    *amountVal,
		NewAmount: *newAmountVal,
	}

	if err := cmd.do(); err != nil {
		panic(err)
	}
}
