include kripto-bot.env
export $(shell sed 's/=.*//' kripto-bot.env)

all: build test

deps:
	go get -t -v ./...

build: deps
	$(PACKAGE)

run:
	./kshell.sh
