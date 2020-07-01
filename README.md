# arbortest

chain your tests in a graph like manner

## example

run command

> go run -v tool/main.go -pkg=example -dir=./example

then start the UI server

> go run -v . -port=3000

and run the tests

> go test -v ./example/... -args --uri=<http://localhost:3000/data.json>

to check the result go to <http://localhost:3000>
