# grpcserv
Tiny gRPC server implementing interaction within the banking infrastructure from HSE golang course.

## Running the server

1. Build and start the server:
   ```bash
   docker-compose up --build
   ```

2. The gRPC server will be available on `localhost:8080`

## Commands

You can interact with the server using the client application. Build the client first:

```bash
go build -o client ./grpc/client/main.go
```

Then run various commands:

```bash
# Create an account
./client -host=localhost -port=8080 -cmd=create -name=john -amount=1000

# Get account details
./client -host=localhost -port=8080 -cmd=get -name=john

# Change account amount
./client -host=localhost -port=8080 -cmd=changeAmount -name=john -newAmount=1500

# Change account name
./client -host=localhost -port=8080 -cmd=changeName -name=john -newName=jane

# Delete account
./client -host=localhost -port=8080 -cmd=delete -name=jane
```
