# arbortest

chain your tests in a graph like manner

## example

run command

> go run -v arborgen/main.go -pkg=example -dir=./example

then start the UI server

> go run -v ./server -port=3000

and run the tests

> go test -v ./example/... -args --arborURL=<http://localhost:3000/data.json>

to check the result go to <http://localhost:3000>

## to install locally the server binary

> make install
