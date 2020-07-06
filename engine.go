package grpc_end

import (
	"context"
	"errors"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type HandleFunc func(c *GRpcContext)
type HandlersChain []HandleFunc

type GRpcEngine struct {
	// funcMap stored all registered HandleFunc, the key composed of controller and action of request.
	funcMap map[string]HandleFunc

	// handlers is an ordered set of HandleFunc, it works like middleware
	handlers HandlersChain

	// a sync pool for GRpcContext
	pool sync.Pool

	// appName, for log
	appName string
}

// NewGRpcEngine returns a GRpcEngine
func NewGRpcEngine(appName string) *GRpcEngine {
	engine := &GRpcEngine{
		funcMap:  make(map[string]HandleFunc),
		handlers: make([]HandleFunc, 0),
		appName:  appName,
	}

	engine.pool.New = func() interface{} {
		return &GRpcContext{}
	}

	return engine
}

// RegisterFunc register a HandleFunc to mapFunc, it will be panic if use a exist key
// i.e Register("user/login", loginAction) mean the request will router to func loginAction
// if request.Controller is 'user' and request.Action is 'login'
func (e *GRpcEngine) RegisterFunc(controller, action string, f HandleFunc) {
	key := controller + "/" + action
	if _, ok := e.funcMap[key]; ok {
		panic("func " + key + " already exists")
	}
	e.funcMap[key] = f
}

// Use attaches a global middleware to the engine. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request.
func (e *GRpcEngine) Use(f ...HandleFunc) {
	e.handlers = append(e.handlers, f...)
}

// newContext returns a GRpcContext for every requests
func (e *GRpcEngine) newContext(ctx context.Context, key string, in *Request) *GRpcContext {
	c := e.pool.Get().(*GRpcContext)
	c.reset()
	c.handlers = append([]HandleFunc{}, e.handlers...)
	c.handlers = append(c.handlers, e.funcMap[key])
	c.ctx = ctx
	c.engine = e
	c.req = in
	c.appName = e.appName
	return c
}

// DoRequest is the implementation of the gRpc protocol.
func (e *GRpcEngine) DoRequest(ctx context.Context, in *Request) (*Response, error) {
	key := in.Controller + "/" + in.Action

	if _, ok := e.funcMap[key]; !ok {
		return nil, errors.New("404")
	}

	c := e.newContext(ctx, key, in)
	defer e.pool.Put(c)
	c.Next()

	return c.resp, nil
}

// Run attaches the funcMap to a gRpc Server and starts tcp listening and serving gRpc requests.
// It is a shortcut for net.Listen and *grpc.Server.Serve.
// Note: this method will execute in a new goroutine.
func (e *GRpcEngine) Run(addr string, opt ...grpc.ServerOption) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	gRpcSvr := grpc.NewServer(opt...)
	RegisterEndServer(gRpcSvr, e)
	go gRpcSvr.Serve(lis)

	return gRpcSvr, nil
}
