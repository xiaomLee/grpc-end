# grpc_end

grpc_end是一套类似gin框架的grpc框架。设计初衷是为了让开发者像使用gin开发web业务一样，方便的开发grpc业务。

grpc_end主要由GRpcEngine和GRpcContext组成。


---------

#### GRpcEngine：

GRpcEngine是整套框架的引擎，负责启动grpc服务，分发处理grpc请求，管理GRpcContext的生命周期。GRpcEngine支持插件式开发（详见middleware以及example文件）


---------

#### GRpcContext

GRpcContext是一个请求的上下文，每个请求的GRpcContext都是唯一的一份，同时它提供了一些辅助的方法。