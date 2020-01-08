# Lab #2 - Simple TCP server/client Golang
----

## Usage

```bash
$ export GOPATH=$(pwd)
$ go get all
$ go build main
$ ./main -m server -p 9000 -n 50 # start server
$ ./main -m client --host localhost -p 9000 # start client
$ ./main --help # see help
```

## In docker

```bash
$ docker build -t maximskripnik/lab2 .
$ docker run --rm --net host maximskripnik/lab2 -m server -p 9000 # start server
$ docker run -ti --rm --net host maximskripnik/lab2 -m client --host localhost -p 9000 # start client
```
