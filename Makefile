build:
	go build -o rucola-chat rucola/*

default: build

run: build
	./rucola-chat

test: build
	./test.sh
