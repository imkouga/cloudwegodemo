.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
	go install github.com/cloudwego/hertz/cmd/hz@latest
	

.PHONY: common
common:	
	protoc 	\
		--proto_path=. \
		--proto_path=./pkg/pb \
		--go_out=paths=source_relative:. \
		$(shell find pkg -name *.proto) \
		

HZ:=cloudwegodemo
.PHONY: $(HZ)
cloudwegodemo:
	hz new --handler_dir cmd/$(HZ)/internal/service \
		--router_dir cmd/$(HZ)/internal/router \
		--model_dir api \
		-force \
		-I idl \
		-I pkg \
		-idl $(shell find idl/$(HZ) -name *.proto) \
	
# protoc 	\
	# 	--proto_path=. \
	# 	--go_out=paths=source_relative:. \
	# 	$(shell find cmd/$(HZ) -name *.proto) \

	protoc 	\
		--proto_path=. \
		--go_out=paths=source_relative:. \
		$(shell find internal -name *.proto) \
		
	rm -rf main.go router_gen.go router.go
			