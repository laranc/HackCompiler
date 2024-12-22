INPUT =

default:
	go build -o ./build/HackCompiler ./cmd/.

run:
	go run ./cmd/. $(INPUT)
