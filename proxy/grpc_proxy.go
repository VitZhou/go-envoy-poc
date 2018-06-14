package proxy

import (
	"google.golang.org/grpc"
	"io"
	"google.golang.org/grpc/codes"
	"golang.org/x/net/context"
	"go-envoy-poc/proxy/encoding"
	"go-envoy-poc/proxy/conn"
	"go-envoy-poc/log"
)

var (
	proxyStreamDesc = &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

func NewGrpcServer() *grpc.Server {
	return grpc.NewServer(grpc.UnknownServiceHandler(nil), streamInterceptor(), grpc.CustomCodec(encoding.Codec()))
}

func streamInterceptor() grpc.ServerOption {
	return grpc.StreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		fullMethodName := info.FullMethod
		backendConn, e := conn.GetConn(ss.Context())
		if e != nil {
			log.Error.Println("获取连接失败", e)
		}
		clientCtx, clientCancel := context.WithCancel(ss.Context())
		clientStream, err := backendConn.NewStream(clientCtx, proxyStreamDesc, fullMethodName, grpc.CallCustomCodec(encoding.Codec()))
		if err != nil {
			log.Error.Println("创建proxyClientStream失败", err)
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
		return err
	})
}

func forwardClientToServer(src grpc.ClientStream, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &encoding.Frame{}
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
		f := &encoding.Frame{}
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
