build:
	go build -o rucola-chat rucola/*

default: build

run:
	go run rucola/*

test: build
	./test.sh
