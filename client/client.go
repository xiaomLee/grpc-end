package client

import (
	"context"
	"time"

	"github.com/xiaomLee/grpc-end"

	"google.golang.org/grpc"
)

var (
	gRpcMapPool *MapPool
)

func InitClient(dial DialFunc) error {
	if dial == nil {
		dial = defaultDialFunc
	}
	gRpcMapPool = NewMapPool(dial, 0, 5*time.Minute)

	// 从本地配置 初始化server列表

	return nil
}

func defaultDialFunc(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return grpc.DialContext(ctx, addr, grpc.WithBlock(), grpc.WithInsecure())
}

func CallEndApi(serverAddress, controller, action string, params map[string]string) ([]byte, error) {
	// Make new request
	request := &grpc_end.Request{
		Controller: controller,
		Action:     action,
		Params:     params,
	}

	// Get a gRpc conn from pool
	pool := gRpcMapPool.GetPool(serverAddress)
	conn, err := pool.Get()
	if err != nil {
		return nil, err
	}

	// Set the request timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Do gRpc request
	resp, err := grpc_end.NewEndClient(conn.GetConn()).DoRequest(ctx, request)
	if err != nil {
		pool.DelErrorClient(conn)
		return nil, err
	}
	_ = pool.Put(conn)

	return resp.Data, nil
}
