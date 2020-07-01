# arbortest

chain your tests in a graph like manner

## example

run command

> go run -v tool/main.go -pkg=example -dir=./example

then

> go test -v ./example/... -args --uri=<http://localhost:3000/data.json>
