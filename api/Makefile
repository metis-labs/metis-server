GOPATH=$(shell go env GOPATH)

GO_TOOLS=google.golang.org/protobuf/cmd/protoc-gen-go \
google.golang.org/grpc/cmd/protoc-gen-go-grpc

PB_OUT=metis.pb.go

INCLUDES= -I=.

all: $(PB_OUT)

tools:
	go get $(GO_TOOLS)

$(PB_OUT): tools metis.proto
	protoc ./metis.proto \
	$(INCLUDES) \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative

clean:
	rm -f $(PB_OUT)

.PHONY: clean, tools
