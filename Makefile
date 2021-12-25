run:
	@cd src/cmd && go run ./*.go

build:
	@cd $(PWD)/src/cmd && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o eatigo .

gen_docs:
	@cp src/cmd/main.go src/
	@cd src && swag init
	@rm -rf src/main.go
	@rm -rf src/cmd/docs
	@mv src/docs src/cmd