test:
	go test ./... -v
test-server:
	go test -v ./client_test.go ./client.go ./server.go 
