

.PHONY: run
run:
	protoc -I book/ book/book.proto --go_out=plugins=grpc:book