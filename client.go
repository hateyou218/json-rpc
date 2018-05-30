package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc/jsonrpc"
	"time"
)

const (
	addr         = "127.0.0.1:1234"
	alpha        = "abcdefghijklmnopqrstuvwxyz"
	digit        = "0123456789"
	key_len      = 2
	val_len      = 20
	worker_count = 10
	max_sleep    = 5000
)

type Args struct {
	Arg1, Arg2 string
}

func GetKey(n int) string {
	b := make([]byte, key_len)
	for i := range b {
		b[i] = alpha[n%len(alpha)]
		n = n / len(alpha)
	}
	return string(b)
}

func GetVal(n, m int) string {
	var b []byte
	if m == 0 {
		b = make([]byte, n)
		for i := range b {
			if n > 1 && i == 0 {
				b[i] = digit[rand.Int()%(len(digit)-1)+1]
			} else {
				b[i] = digit[rand.Int()%len(digit)]
			}
		}
	} else {
		b = make([]byte, n+m+1)
		for i := range b {
			if n > 1 && i == 0 {
				b[i] = digit[rand.Int()%(len(digit)-1)+1]
			} else if i == n {
				b[i] = '.'
			} else {
				b[i] = digit[rand.Int()%len(digit)]
			}
		}
	}
	return string(b)
}

func SendRequest(addr string) (string, error) {
	const (
		NEW = iota
		SET
		DEL
		ADD
		SUB
		MUL
		DIV
		CMD_MAX
	)
	client, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	var reply string
	args := new(Args)
	cmd := rand.Int() % CMD_MAX
	switch cmd {
	case NEW:
		args.Arg1 = GetKey(rand.Int())
		args.Arg2 = GetVal(rand.Int()%val_len+1, rand.Int()%val_len)
		c := jsonrpc.NewClient(client)
		err = c.Call("Computation.New", args, &reply)
		client.Close()
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("Result: Create %s with value %s", args.Arg1, args.Arg2), nil
		}
	case SET:
		args.Arg1 = GetKey(rand.Int())
		args.Arg2 = GetVal(rand.Int()%val_len+1, rand.Int()%val_len)
		c := jsonrpc.NewClient(client)
		err = c.Call("Computation.Set", args, &reply)
		client.Close()
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("Result: Set %s with value %s", args.Arg1, args.Arg2), nil
		}
	case DEL:
		args.Arg1 = GetKey(rand.Int())
		c := jsonrpc.NewClient(client)
		err = c.Call("Computation.Del", args, &reply)
		client.Close()
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("Result: Delete %s", args.Arg1), nil
		}
	case ADD:
		if rand.Int()%10 >= 7 {
			args.Arg1 = GetVal(rand.Int()%val_len+1, rand.Int()%val_len)
		} else {
			args.Arg1 = GetKey(rand.Int())
		}
		if rand.Int()%10 >= 7 {
			args.Arg2 = GetVal(rand.Int()%val_len+1, rand.Int()%val_len)
		} else {
			args.Arg2 = GetKey(rand.Int())
		}
		c := jsonrpc.NewClient(client)
		err = c.Call("Computation.Add", args, &reply)
		client.Close()
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("Result: %s + %s = %s", args.Arg1, args.Arg2, reply), nil
		}
	case SUB:
		if rand.Int()%10 >= 7 {
			args.Arg1 = GetVal(rand.Int()%val_len+1, rand.Int()%val_len)
		} else {
			args.Arg1 = GetKey(rand.Int())
		}
		if rand.Int()%10 >= 7 {
			args.Arg2 = GetVal(rand.Int()%val_len+1, rand.Int()%val_len)
		} else {
			args.Arg2 = GetKey(rand.Int())
		}
		c := jsonrpc.NewClient(client)
		err = c.Call("Computation.Sub", args, &reply)
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("Result: %s - %s = %s", args.Arg1, args.Arg2, reply), nil
		}
	case MUL:
		if rand.Int()%10 >= 7 {
			args.Arg1 = GetVal(rand.Int()%val_len+1, rand.Int()%val_len)
		} else {
			args.Arg1 = GetKey(rand.Int())
		}
		if rand.Int()%10 >= 7 {
			args.Arg2 = GetVal(rand.Int()%val_len+1, rand.Int()%val_len)
		} else {
			args.Arg2 = GetKey(rand.Int())
		}
		c := jsonrpc.NewClient(client)
		err = c.Call("Computation.Mul", args, &reply)
		client.Close()
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("Result: %s * %s = %s", args.Arg1, args.Arg2, reply), nil
		}
	case DIV:
		if rand.Int()%10 >= 7 {
			args.Arg1 = GetVal(rand.Int()%val_len+1, rand.Int()%val_len)
		} else {
			args.Arg1 = GetKey(rand.Int())
		}
		if rand.Int()%10 >= 7 {
			args.Arg2 = GetVal(rand.Int()%val_len+1, rand.Int()%val_len)
		} else {
			args.Arg2 = GetKey(rand.Int())
		}
		c := jsonrpc.NewClient(client)
		err = c.Call("Computation.Div", args, &reply)
		client.Close()
		if err != nil {
			return "", err
		} else {
			return fmt.Sprintf("Result: %s / %s = %s", args.Arg1, args.Arg2, reply), nil
		}
	default:
		client.Close()
		return "", errors.New("Something went wrong")
	}
}

func worker(addr string, interval time.Duration) {
	log.Printf("Starting worker with sleeping interval = %v\n", interval*time.Millisecond)
	for {
		s, err := SendRequest(addr)
		if err != nil {
			log.Printf("Error: %v\n", err)
		} else {
			log.Println(s)
		}
		time.Sleep(interval * time.Millisecond)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 1; i < worker_count; i++ {
		go worker(addr, time.Duration(rand.Int63()%max_sleep+1))
	}
	worker(addr, time.Duration(rand.Int63()%max_sleep+1))
}
