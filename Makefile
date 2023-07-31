BIN_NAME=lambda
local:
	mv main.go main.lambda 
	mv main.local main.go

build:
	GOARCH=amd64 GOOS=windows go build -o ./target/${BIN_NAME}-win main.go

run: build
	./target/${BIN_NAME}-win

clean:
	go clean
	rm ./target/${BIN_NAME}-win
	mv main.go main.local
	mv main.lambda main.go

testing: 
	go test -v ./...