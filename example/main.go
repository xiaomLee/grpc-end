package main

//go:generate protoc -I ../ grpc_end.proto --go_out=plugins=grpc:../

import (
	grpc_end "grpc-end"
	"grpc-end/middleware"
	"math"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

func main() {
	ser, err := newGRpcEngine().Run(":1234", grpc.MaxRecvMsgSize(math.MaxInt32))
	if err != nil {
		panic(err.Error())
	}
	println("server start success")

	// hold here and deal request...
	// time.Sleep(time.Second * 100)
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, os.Kill)
	select {
	case <-exit:
		ser.GracefulStop()
		println("server stop success")
	}

}

func newGRpcEngine() *grpc_end.GRpcEngine {
	engine := grpc_end.NewGRpcEngine("MyAppName")
	engine.RegisterFunc("hello", "world", sayHi)

	engine.Use(middleware.Recover)
	engine.Use(middleware.Logger)

	return engine
}

func sayHi(c *grpc_end.GRpcContext) {
	name := c.StringParamDefault("name", "")
	c.SuccessResponse("hi " + name)
}
