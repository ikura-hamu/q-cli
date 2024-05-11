# Go compiler
GO := go

# Output binary name
BINARY := q

DIR := ~/bin

# Build target
build:
	$(GO) build -o $(BINARY) -ldflags "-s -w -X main.version=0.0.1" .

install:
	cp $(BINARY) $(DIR)

# Clean target
clean:
	rm -f $(BINARY)
