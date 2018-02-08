package main

import (
	"google.golang.org/grpc"
	_ "golang.org/x/net/context"
	pb "gprcstudy5/proto"
	_ "fmt"
	"net"
	"log"
	"google.golang.org/grpc/grpclog"
	"golang.org/x/net/context"
	"fmt"
	"io"
	"time"
	"strconv"
)
const(
	PORT=":8077"
)
type SServer struct{}
//简单模式
func (ss *SServer)StudentAdd(ctx context.Context,in *pb.StudentResponse)(*pb.ResultResponse,error){
	fmt.Println("StudentAdd================简单模式================")
	//简单模式不需要调用recv和send。数据直接通过方法的参数和返回值来接收发送。
	re:=&pb.ResultResponse{1,"noerror"}
	return re,nil
}
//客户端流模式
func (ss *SServer)GetHelloTest(stream pb.SelfManage_GetHelloTestServer)(error){
	fmt.Println("GetHelloTest================客户端模式================")
	//服务器端的代码stream是SelfManage_GetHelloTestServer类型，这种类型只有recv和sendandclose两种方法接收和发送并关闭
	//而客户端跟着相反是send和closeandrecv两种方法。
	//客户端流模式要先接收然后再发送。
	//使用for循环接收数据，因为接收的是流模式的StudentRequest类型，所以每次循环得到的是一个studentRequest类型
	for {
		in, err := stream.Recv()//通过for循环接收客户端传来的流，该方法协程不安全。
		//必须把判断EOF的条件放在nil前面，因为如果获取完数据err会返回EOF类型
		if err == io.EOF {//在读取完所有数据之后会返回EOF来执行该条件。
			fmt.Println("read done")
			//发送（只在读取完数据后发送了一次）
			rp:=&pb.ResultResponse{1000,"nxin"}
			//发送并关闭。该方法的返回值是error类型。但客户端并未得到返回值，只收到sendandclose发送的数据。
			//注意：该方法内部调用的SendMsg方法，但是该方法协程不安全。不同协程可以同时调用会出现问题。
			return stream.SendAndClose(rp)
		}
		if err != nil {
			fmt.Println("ERR IN GetHelloTest", err)
			return err
		}
		//打印出每次接收的数据
		fmt.Println("Recieved information:",in)
	}
}
//服务器端流模式
func (ss *SServer)GetNxinMethod(in *pb.StudentRequest,stream pb.SelfManage_GetNxinMethodServer)error{
	fmt.Println("GetNxinMethod================服务器端模式================")
	fmt.Println(in.Sid)//打印出来客户端参数传来的数据（此处无意义，只是用来判断客户端是否传来数据）
	p:=&pb.ResultResponse{10,"nnnnn"}
	stream.Send(p)//服务器流模式发送数据。此处不用关闭发送，因为没找到关闭的方法。
	return nil
}
//双向模式
func (ss *SServer)GetStudent(stream pb.SelfManage_GetStudentServer)(error){
	for{
		//接收
		in,err:=stream.Recv()
		if err==io.EOF{
			fmt.Println("read done")
			return nil
		}
		if err!=nil{
			grpclog.Fatal("some error in recv",err)
		}
		fmt.Println(in,strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
		//发送
		rp:=&pb.ResultResponse{1000,"nxin"}
		stream.Send(rp)
	}
	return nil
}


func main(){
	lis,err:=net.Listen("tcp",PORT)
	if err!=nil{
		grpclog.Fatal("error occur in listen")
		log.Fatal("testtttt")
	}
	s:=grpc.NewServer()
	pb.RegisterSelfManageServer(s,&SServer{})
	s.Serve(lis)
}