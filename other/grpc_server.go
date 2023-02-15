// package main

// import (
//     "context"
//     "log"
//     "net"
//     "google.golang.org/grpc"
//     pb "path/to/your/proto/package"
// )

// type server struct{}

// func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
//     log.Printf("Received: %v", in)
//     return &pb.HelloReply{Message: "Hello " + in.Name}, nil
// }

// func main() {
//     // 监听TCP端口
//     lis, err := net.Listen("tcp", "localhost:9527")
//     if err != nil {
//         log.Fatalf("failed to listen: %v", err)
//     }

//     // 创建一个新的gRPC服务器
//     s := grpc.NewServer()

//     // 注册你的服务
//     pb.RegisterYourServiceServer(s, &server{})

//     // 开始服务
//     if err := s.Serve(lis); err != nil {
//         log.Fatalf("failed to serve: %v", err)
//     }
// }