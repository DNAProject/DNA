package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"
	"testing"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

type ArithAddResp struct {
	Id     interface{} `json:"id"`
	Result Reply       `json:"result"`
	Error  interface{} `json:"error"`
}

func (t *Arith) Add(args *Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

func (t *Arith) Mul(args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func (t *Arith) Div(args *Args, reply *Reply) error {
	print("Run to divide\n")
	if args.B == 0 {
		print("Run to divide by 0\n")
		return errors.New("divide by zero")
	}
	reply.C = args.A / args.B
	return nil
}

func (t *Arith) Error(args *Args, reply *Reply) error {
	panic("ERROR")
}

func init() {
	rpc.Register(new(Arith))
}

func TestServerNoParams(t *testing.T) {
	cli, srv := net.Pipe()
	defer cli.Close()
	go rpc.ServeConn(srv)
	dec := json.NewDecoder(cli)

	fmt.Fprintf(cli, `{"method": "Arith.Add", "id": "123"}`)
	var resp ArithAddResp
	if err := dec.Decode(&resp); err != nil {
		t.Fatalf("Decode after no params: %s", err)
	}
	if resp.Error == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestServerEmptyMessage(t *testing.T) {
	cli, srv := net.Pipe()
	defer cli.Close()
	go rpc.ServeConn(srv)
	dec := json.NewDecoder(cli)

	fmt.Fprintf(cli, "{}")
	var resp ArithAddResp
	if err := dec.Decode(&resp); err != nil {
		t.Fatalf("Decode after empty: %s", err)
	}
	if resp.Error == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestServer() {
	cli, srv := net.Pipe()
	defer cli.Close()
	go rpc.ServeConn(srv)
	dec := json.NewDecoder(cli)

	print("Just for testing server")

	// Send hand-coded requests to server, parse responses.
	for i := 0; i < 10; i++ {
		fmt.Fprint(cli, `{"method": "Arith.Add", "id": "\u%04d", "params": [{"A": %d, "B": %d}]}`, i, i, i+1)
		var resp ArithAddResp
		err := dec.Decode(&resp)
		if err != nil {
			print("Decode: %s", err)
		}
		if resp.Error != nil {
			print("resp.Error: %s", resp.Error)
		}
		if resp.Id.(string) != string(i) {
			print("resp: bad id %q want %q", resp.Id.(string), string(i))
		}
		if resp.Result.C != 2*i+1 {
			print("resp: bad result: %d+%d=%d", i, i+1, resp.Result.C)
		}
	}
}

func TestClient() {
	// Assume server is okay (TestServer is above).
	// Test client against server.
	print("testing at client 1 \n")
	cli, srv := net.Pipe()
	go rpc.ServeConn(srv)

	client := rpc.NewClient(cli)
	defer client.Close()

	// Synchronous calls
	args := &Args{7, 8}
	reply := new(Reply)
	err := client.Call("Arith.Add", args, reply)
	if err != nil {
		print("Add: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A+args.B {
		print("Add: got %d expected %d", reply.C, args.A+args.B)
	}
	print("testing at client 2 \n")

	args = &Args{7, 8}
	reply = new(Reply)
	err = client.Call("Arith.Mul", args, reply)
	if err != nil {
		print("Mul: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A*args.B {
		print("Mul: got %d expected %d", reply.C, args.A*args.B)
	}

	fmt.Print("testing at client 3 result is ", reply.C, "\n")
	// Out of order.
	args = &Args{7, 8}
	mulReply := new(Reply)
	mulCall := client.Go("Arith.Mul", args, mulReply, nil)
	addReply := new(Reply)
	addCall := client.Go("Arith.Add", args, addReply, nil)

	addCall = <-addCall.Done
	if addCall.Error != nil {
		print("Add: expected no error but got string %q", addCall.Error.Error())
	}
	if addReply.C != args.A+args.B {
		print("Add: got %d expected %d", addReply.C, args.A+args.B)
	}

	mulCall = <-mulCall.Done
	if mulCall.Error != nil {
		print("Mul: expected no error but got string %q", mulCall.Error.Error())
	}
	if mulReply.C != args.A*args.B {
		print("Mul: got %d expected %d", mulReply.C, args.A*args.B)
	}

	print("testing at client 4 \n")
	// Error test
	args = &Args{7, 0}
	reply = new(Reply)
	err = client.Call("Arith.Div", args, reply)
	// expect an error: zero divide
	if err == nil {
		print("Div: expected error")
	} else if err.Error() != "divide by zero" {
		print("Div: expected divide by zero error; got", err)
	}
	print("testing at client 5 \n")
}

func TestMalformedInput(t *testing.T) {
	cli, srv := net.Pipe()
	go cli.Write([]byte(`{id:1}`)) // invalid json
	rpc.ServeConn(srv)             // must return, not loop
}

func TestMalformedOutput(t *testing.T) {
	cli, srv := net.Pipe()
	go srv.Write([]byte(`{"id":0,"result":null,"error":null}`))
	go ioutil.ReadAll(srv)

	client := rpc.NewClient(cli)
	defer client.Close()

	args := &Args{7, 8}
	reply := new(Reply)
	err := client.Call("Arith.Add", args, reply)
	if err == nil {
		t.Error("expected error")
	}
}

func TestServerErrorHasNullResult(t *testing.T) {
	var out bytes.Buffer
	sc := jsonrpc.NewServerCodec(struct {
		io.Reader
		io.Writer
		io.Closer
	}{
		Reader: strings.NewReader(`{"method": "Arith.Add", "id": "123", "params": []}`),
		Writer: &out,
		Closer: ioutil.NopCloser(nil),
	})
	r := new(rpc.Request)
	if err := sc.ReadRequestHeader(r); err != nil {
		t.Fatal(err)
	}
	const valueText = "the value we don't want to see"
	const errorText = "some error"
	err := sc.WriteResponse(&rpc.Response{
		ServiceMethod: "Method",
		Seq:           1,
		Error:         errorText,
	}, valueText)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), errorText) {
		t.Fatalf("Response didn't contain expected error %q: %s", errorText, &out)
	}
	if strings.Contains(out.String(), valueText) {
		t.Errorf("Response contains both an error and value: %s", &out)
	}
}

func TestUnexpectedError(t *testing.T) {
	cli, srv := myPipe()
	go cli.PipeWriter.CloseWithError(errors.New("unexpected error!")) // reader will get this error
	rpc.ServeConn(srv)                                                // must return, not loop
}

// Copied from package net.
func myPipe() (*pipe, *pipe) {
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()

	return &pipe{r1, w2}, &pipe{r2, w1}
}

type pipe struct {
	*io.PipeReader
	*io.PipeWriter
}

type pipeAddr int

func (pipeAddr) Network() string {
	return "pipe"
}

func (pipeAddr) String() string {
	return "pipe"
}

func (p *pipe) Close() error {
	err := p.PipeReader.Close()
	err1 := p.PipeWriter.Close()
	if err == nil {
		err = err1
	}
	return err
}

func (p *pipe) LocalAddr() net.Addr {
	return pipeAddr(0)
}

func (p *pipe) RemoteAddr() net.Addr {
	return pipeAddr(0)
}

func (p *pipe) SetTimeout(nsec int64) error {
	return errors.New("net.Pipe does not support timeouts")
}

func (p *pipe) SetReadTimeout(nsec int64) error {
	return errors.New("net.Pipe does not support timeouts")
}

func (p *pipe) SetWriteTimeout(nsec int64) error {
	return errors.New("net.Pipe does not support timeouts")
}

func main() {

	print("Just for testing\n")
	//TestServer()
	TestClient()
}
