package main

import (
	"net"
	"log"
	"google.golang.org/grpc"
	"go-envoy-poc/proxy/grpcp"
	"google.golang.org/grpc/reflection"
	"fmt"
	"context"
	"io"
	"google.golang.org/grpc/codes"
	"errors"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	lis, err := net.Listen("tcp", ":9933")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serverOption := grpc.UnknownServiceHandler(nil)
	serverOption2 := grpc.StreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		fullMethodName, e1:= grpc.MethodFromServerStream(ss)
		if !e1{
			log.Println("获取方法失败", e1)
			return errors.New("获取方法失败")
		}
		backendConn, e := grpc.DialContext(ss.Context(), "localhost:50051",grpc.WithInsecure())
		if e !=nil{
			log.Println("创建客户端连接失败", e)
			return e
		}
		clientStreamDescForProxying := &grpc.StreamDesc{
			ServerStreams: true,
			ClientStreams: true,
		}
		clientCtx, clientCancel := context.WithCancel(ss.Context())
		clientStream,  err := backendConn.NewStream(clientCtx, clientStreamDescForProxying, fullMethodName, grpc.CallCustomCodec(grpcp.Codec()))
		if err != nil{
			log.Println("aaaaaaaaaa", err)
			return err
		}
		s2cErrChan := forwardServerToClient(ss, clientStream)
		c2sErrChan := forwardClientToServer(clientStream, ss)
		for i := 0; i < 2; i++ {
			select {
			case s2cErr := <-s2cErrChan:
				if s2cErr == io.EOF {
					// this is the happy case where the sender has encountered io.EOF, and won't be sending anymore./
					// the clientStream>serverStream may continue pumping though.
					clientStream.CloseSend()
					break
				} else {
					clientCancel()
					return grpc.Errorf(codes.Internal, "failed proxying s2c: %v", s2cErr)
				}
			case c2sErr := <-c2sErrChan:
				ss.SetTrailer(clientStream.Trailer())
				// c2sErr will contain RPC error from client code. If not io.EOF return the RPC error as server stream error.
				if c2sErr != io.EOF {
					return c2sErr
				}
				return nil
			}
		}
		return nil
	})


	s := grpc.NewServer(serverOption,serverOption2,grpc.CustomCodec(grpcp.Codec()))
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func forwardClientToServer(src grpc.ClientStream, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &grpcp.Frame{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}
			if i == 0 {
				md, err := src.Header()
				if err != nil {
					ret <- err
					break
				}
				if err := dst.SendHeader(md); err != nil {
					ret <- err
					break
				}
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}

func forwardServerToClient(src grpc.ServerStream, dst grpc.ClientStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &grpcp.Frame{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}