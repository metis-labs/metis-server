GOPATH=$(shell go env GOPATH)

GO_TOOLS=github.com/gogo/protobuf/protoc-gen-gofast

PB_OUT=metis.pb.go

INCLUDES= -I=. \
	-I=$(GOPATH)/src \
	-I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf

all: $(PB_OUT)

.PHONY: tools
tools:
	go get $(GO_TOOLS)

$(PB_OUT): tools metis.proto
	protoc metis.proto \
	$(INCLUDES) \
	--gofast_out=plugins=grpc:.

.PHONY: clean
clean:
	rm -f $(PB_OUT)
