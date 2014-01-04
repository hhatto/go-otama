
all:
	@echo "clean"

test:
	go test

clean:
	rm -rf go-otama.test examples/data/*
