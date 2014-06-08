
all:
	@echo "test"
	@echo "install"
	@echo "clean"

install:
	go get -v -u github.com/hhatto/go-otama

test:
	go test

clean:
	rm -rf go-otama.test examples/data/*
