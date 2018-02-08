package main

import (
	"google.golang.org/grpc"
	pb "gprcstudy5/proto"
	"fmt"
	//"sync"
	"google.golang.org/grpc/grpclog"
	"context"
	"io"
	"time"
	"strconv"
)
const (
	address = "127.0.0.1:8077"
)
func StudentAdd_Method(client pb.SelfManageClient,info *pb.StudentResponse){
	fmt.Println("StudentAdd_Method================简单模式================")
	//简单模式下直接调用方法就可以传递参数和获得返回值
	re,err:=client.StudentAdd(context.Background(),info)
	if err!=nil{
		fmt.Println("some error occur in studentadd",err)
	}
	fmt.Println("studentadd_method:",re.Affectrows,re.Errorinfo)
}
func GetHelloTest_Method(client pb.SelfManageClient){
	fmt.Println("GetHelloTest_Method================客户端模式================")
	//此处相当于创建了一个GetHelloTest的通道。通过通道发送和接收数据。
	re,err:=client.GetHelloTest(context.Background())
	if err!=nil{
		fmt.Println("some error occur in gethellotest",err)
	}
	//初始化要发送给服务器端的数据，并发送（此处连续发送10次）
	for i:=0;i<10;i++{
		sq:=&pb.StudentRequest{1}
		re.Send(sq)//客户端要先发送数据再接收
	}
	//使用for循环接收数据
	for {
		//先关闭send然后再接收数据。该方法内部调用了CloseSend和RecvMsg方法。但是后者协程不安全
		r, err2 := re.CloseAndRecv()
		if err2==io.EOF{
			fmt.Println("recv done")
			return
		}
		if err2 != nil {
			fmt.Println("some error occur in gethellotest recv")
		}
		fmt.Println(r)
	}
}
func GetNxinMethod_Method(client pb.SelfManageClient){
	fmt.Println("GetNxinMethod_Method================服务器端模式================")
	s:=&pb.StudentRequest{2}
	//向服务器端以参数形式发送数据s，并创建了一个通道，通过通道来获取服务器端返回的流
	re,err:=client.GetNxinMethod(context.Background(),s)
	if err!=nil{
		fmt.Println("some error occur in getnxinmethod")
	}
	for{
		//上面已经以参数的形式发送了数据，此处直接使用Recv接收数据。因是参数发送，所有不用使用CloseAndRecv
		v,e:=re.Recv()
		if e==io.EOF{
			fmt.Println("recv done")
			return
		}
		if e!=nil{
			fmt.Println("e not nil",e)
			return
		}
		fmt.Println("get nxin method recv",v)
	}
}
func GetStudent_Method(client pb.SelfManageClient){
	fmt.Println("GetStudent_Method================双向模式================")
	re,err:=client.GetStudent(context.Background())
	if err!=nil{
		fmt.Println("get student error",err)
	}
	//创建一个通道作为控制协程协作
	waitc:=make(chan struct{})
	//开启协程进行接收信息
	go func() {
		for {
			in, errr := re.Recv()
			if errr == io.EOF {
				fmt.Println("read done ")
				close(waitc)
				return
			}
			if errr != nil {
				fmt.Println("getstudent recv", errr)
			}
			fmt.Println("Got message at point(%s, %d)", in.Errorinfo, in.Affectrows, strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
		}
	}()
	//主进程继续执行发送
	s1:=&pb.StudentRequest{2}
	er:=re.Send(s1)
	if er!=nil{
		fmt.Println("get student send err ",er)
	}
	re.CloseSend()
	<-waitc
}
func main(){
	conn,err:=grpc.Dial(address,grpc.WithInsecure())
	if err!=nil{
		grpclog.Fatal("some error in dial")
	}
	defer conn.Close()

	client:=pb.NewSelfManageClient(conn)

	s:=&pb.StudentResponse{1,2,"wang","1888888","beijing"}
	StudentAdd_Method(client,s)
	GetHelloTest_Method(client)

	GetNxinMethod_Method(client)
	GetStudent_Method(client)
}
