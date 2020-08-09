# arbortest

Chain your tests in a graph like manner.

## example

to generate the arbor test file run command

> go run -v arborgen/main.go -pkg=example -dir=./example

then compile the UI files from the `/web` folder

> yarn build

then start the UI server

> go run broker.go server.go -port=3000

and run the tests

> go test -v ./example/... -args --arborURL=<http://localhost:3000/data/>

to check the result go to <http://localhost:3000>

## to install the UI server locally

> make install-server

then you can start the UI server: `arbortest -port=3000`

## to install the generator locally

> make install-gen

then you can run it: `arborgen -pkg=example_test -dir ./testdir`
