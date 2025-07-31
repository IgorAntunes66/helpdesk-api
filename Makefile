gen:
	protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/ticket.proto 

clean:
	rm pkg/pb/*.go

run:
	go run main.go