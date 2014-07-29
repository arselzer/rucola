build:
	go build -o bin/rucola main.go

default: build

run: build
	./bin/rucola

test: build
	./test.sh

clean:
	rm -r ./db
