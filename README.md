# etcd-leaser
基于etcd 的服务注册和服务发现
1. 安装golang开发环境
2. 安装etcd组件
3. 克隆代码 git clone https://github.com/sendy-109/etcd-leaser.git
4. go mod init leaser
5. go build newleaser.go
6. 由于etcdV3和grpc3.0 不兼容 如果编译的时候报错:undefined: balancer.PickOptions 
   需要在go.mod末尾添加一行    replace google.golang.org/grpc => google.golang.org/grpc v1.26.0    指定grpc版本
7. go build watchleaser.go
