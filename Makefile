GOCMD=go
GOBUILD=$(GOCMD) build
GOENV=$(GOCMD) env
FLAGS=-trimpath

build:
	env GOARCH=amd64 GOOS=linux $(GOBUILD) $(FLAGS) -o chat
