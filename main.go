package main

import (
	"filestore/config"
	"github.com/micro/go-micro"
	"log"
	"time"
	// k8s "github.com/micro/kubernetes/go/micro"
	"filestore/service/account/handler"
	proto "filestore/service/account/proto"
	dbproxy "filestore/service/dbproxy/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
)

func main() {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{config.Consul}
	})
	service := micro.NewService(
		
		micro.Registry(reg),
		micro.Name("go.micro.service.user"),
		micro.RegisterTTL(time.Second*10),//注册服务的过期时间
		micro.RegisterInterval(time.Second*5),//间隔多久再次注册服务
	)
	// 初始化service, 解析命令行参数等
	// 初始化
	service.Init()
	dbproxy.Init(service)
	proto.RegisterUserServiceHandler(service.Server(),new(handler.User))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
