package main

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

type Args struct {
	Arg1, Arg2 string
}

type Computation struct {
	m map[string]*big.Float
	l chan int
}

func (c *Computation) New(args *Args, reply *string) error {
	log.Printf("%v\n", args)
	_, _, err := big.ParseFloat(args.Arg1, 10, 30, big.ToNearestEven)
	if err == nil {
		return errors.New("Value cannot be a key")
	}
	f, _, err := big.ParseFloat(args.Arg2, 10, 30, big.ToNearestEven)
	if err != nil {
		return err
	}
	select {
	case <-c.l:
		if _, ok := c.m[args.Arg1]; ok {
			c.l <- 0
			return errors.New(fmt.Sprintf("New failed because key %s already exists", args.Arg1))
		}
		c.m[args.Arg1] = f
		c.l <- 0
	case <-time.After(10 * time.Millisecond):
		return errors.New("New failed because of time out")
	}
	return nil
}

func (c *Computation) Set(args *Args, reply *string) error {
	log.Printf("%v\n", args)
	f, _, err := big.ParseFloat(args.Arg2, 10, 30, big.ToNearestEven)
	if err != nil {
		return err
	}
	select {
	case <-c.l:
		if _, ok := c.m[args.Arg1]; !ok {
			c.l <- 0
			return errors.New(fmt.Sprintf("Set failed because key %s doesn't exist", args.Arg1))
		}
		c.m[args.Arg1] = f
		c.l <- 0
	case <-time.After(10 * time.Millisecond):
		return errors.New("Set failed because of time out")
	}
	return nil
}

func (c *Computation) Del(args *Args, reply *string) error {
	log.Printf("%v\n", args)
	select {
	case <-c.l:
		if _, ok := c.m[args.Arg1]; !ok {
			c.l <- 0
			return errors.New(fmt.Sprintf("Del failed because key %s doesn't exist", args.Arg1))
		}
		delete(c.m, args.Arg1)
		c.l <- 0
	case <-time.After(10 * time.Millisecond):
		return errors.New("Del failed because of time out")
	}
	return nil
}

func (c *Computation) Add(args *Args, reply *string) error {
	log.Printf("%v\n", args)
	f1, _, err1 := big.ParseFloat(args.Arg1, 10, 30, big.ToNearestEven)
	f2, _, err2 := big.ParseFloat(args.Arg2, 10, 30, big.ToNearestEven)
	if err1 == nil && err2 == nil {
		*reply = new(big.Float).Add(f1, f2).Text('f', 30)
	} else {
		select {
		case <-c.l:
			if err1 != nil {
				if _, ok := c.m[args.Arg1]; !ok {
					c.l <- 0
					return errors.New(fmt.Sprintf("Add failed because key %s doesn't exist", args.Arg1))
				}
				f1 = c.m[args.Arg1]
			}
			if err2 != nil {
				if _, ok := c.m[args.Arg2]; !ok {
					c.l <- 0
					return errors.New(fmt.Sprintf("Add failed because key %s doesn't exist", args.Arg2))
				}
				f2 = c.m[args.Arg2]
			}
			c.l <- 0
		case <-time.After(10 * time.Millisecond):
			return errors.New("Add failed because of time out")
		}
		*reply = new(big.Float).Add(f1, f2).Text('f', 30)
	}
	return nil
}

func (c *Computation) Sub(args *Args, reply *string) error {
	log.Printf("%v\n", args)
	f1, _, err1 := big.ParseFloat(args.Arg1, 10, 30, big.ToNearestEven)
	f2, _, err2 := big.ParseFloat(args.Arg2, 10, 30, big.ToNearestEven)
	if err1 == nil && err2 == nil {
		*reply = new(big.Float).Sub(f1, f2).Text('f', 30)
	} else {
		select {
		case <-c.l:
			if err1 != nil {
				if _, ok := c.m[args.Arg1]; !ok {
					c.l <- 0
					return errors.New(fmt.Sprintf("Add failed because key %s doesn't exist", args.Arg1))
				}
				f1 = c.m[args.Arg1]
			}
			if err2 != nil {
				if _, ok := c.m[args.Arg2]; !ok {
					c.l <- 0
					return errors.New(fmt.Sprintf("Add failed because key %s doesn't exist", args.Arg2))
				}
				f2 = c.m[args.Arg2]
			}
			c.l <- 0
		case <-time.After(10 * time.Millisecond):
			return errors.New("Sub failed because of time out")
		}
		*reply = new(big.Float).Sub(f1, f2).Text('f', 30)
	}
	return nil
}

func (c *Computation) Mul(args *Args, reply *string) error {
	log.Printf("%v\n", args)
	f1, _, err1 := big.ParseFloat(args.Arg1, 10, 30, big.ToNearestEven)
	f2, _, err2 := big.ParseFloat(args.Arg2, 10, 30, big.ToNearestEven)
	if err1 == nil && err2 == nil {
		*reply = new(big.Float).Mul(f1, f2).Text('f', 30)
	} else {
		select {
		case <-c.l:
			if err1 != nil {
				if _, ok := c.m[args.Arg1]; !ok {
					c.l <- 0
					return errors.New(fmt.Sprintf("Add failed because key %s doesn't exist", args.Arg1))
				}
				f1 = c.m[args.Arg1]
			}
			if err2 != nil {
				if _, ok := c.m[args.Arg2]; !ok {
					c.l <- 0
					return errors.New(fmt.Sprintf("Add failed because key %s doesn't exist", args.Arg2))
				}
				f2 = c.m[args.Arg2]
			}
			c.l <- 0
		case <-time.After(10 * time.Millisecond):
			return errors.New("Mul failed because of time out")
		}
		*reply = new(big.Float).Mul(f1, f2).Text('f', 30)
	}
	return nil
}

func (c *Computation) Div(args *Args, reply *string) error {
	log.Printf("%v\n", args)
	f1, _, err1 := big.ParseFloat(args.Arg1, 10, 30, big.ToNearestEven)
	f2, _, err2 := big.ParseFloat(args.Arg2, 10, 30, big.ToNearestEven)
	if err1 == nil && err2 == nil {
		*reply = new(big.Float).Quo(f1, f2).Text('f', 30)
	} else {
		select {
		case <-c.l:
			if err1 != nil {
				if _, ok := c.m[args.Arg1]; !ok {
					c.l <- 0
					return errors.New(fmt.Sprintf("Add failed because key %s doesn't exist", args.Arg1))
				}
				f1 = c.m[args.Arg1]
			}
			if err2 != nil {
				if _, ok := c.m[args.Arg2]; !ok {
					c.l <- 0
					return errors.New(fmt.Sprintf("Add failed because key %s doesn't exist", args.Arg2))
				}
				f2 = c.m[args.Arg2]
			}
			c.l <- 0
		case <-time.After(10 * time.Millisecond):
			return errors.New("Div failed because of time out")
		}
		*reply = new(big.Float).Quo(f1, f2).Text('f', 30)
	}
	return nil
}

func main() {
	c := new(Computation)
	c.m = make(map[string]*big.Float)
	c.l = make(chan int, 1)
	c.l <- 0
	server := rpc.NewServer()
	server.Register(c)
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("Listen error:", e)
	}
	log.Printf("Server starting\n")
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("Accept error: " + err.Error())
		} else {
			log.Printf("New connection established\n")
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}
