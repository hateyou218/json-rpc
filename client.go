package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc/jsonrpc"
	"time"
    "flag"
)

const (
	alpha        = "abcdefghijklmnopqrstuvwxyz"
	digit        = "0123456789"
)

type Args struct {
	Arg1, Arg2 string
}

var (
    addr string
    port string
	key_len int
	val_len int
    worker_count int
    max_sleep int
)

func init() {
    const (
        default_addr = "127.0.0.1"
        default_port = "1234"
        default_worker_count = 1
        default_max_sleep = 5000
        default_key_len = 2
        default_val_len = 20
        usage_addr = "Target address"
        usage_port = "Target port"
        usage_worker_count = "The number of workers sending requests to the server"
        usage_max_sleep = "Max possible sleeping time (in ms) between two consecutive requests of a worker"
        usage_key_len = "Length of each generated key"
        usage_val_len = "Max possible digits of each generated value"
    )
    flag.StringVar(&addr, "addr", default_addr, usage_addr)
    flag.StringVar(&addr, "a", default_addr, usage_addr)
    flag.StringVar(&port, "port", default_port, usage_port)
    flag.StringVar(&port, "p", default_port, usage_port)
    flag.IntVar(&key_len, "key", default_key_len, usage_key_len)
    flag.IntVar(&key_len, "k", default_key_len, usage_key_len)
    flag.IntVar(&val_len, "val", default_val_len, usage_val_len)
    flag.IntVar(&val_len, "v", default_val_len, usage_val_len)
    flag.IntVar(&worker_count, "worker", default_worker_count, usage_worker_count)
    flag.IntVar(&worker_count, "wc", default_worker_count, usage_worker_count)
    flag.IntVar(&max_sleep, "sleep", default_max_sleep, usage_max_sleep)
    flag.IntVar(&max_sleep, "ms", default_max_sleep, usage_max_sleep)
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

func SendRequest(target string) (string, error) {
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
	client, err := net.Dial("tcp", target)
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

func worker(target string, interval time.Duration) {
	log.Printf("Starting worker with sleeping interval = %v\n", interval*time.Millisecond)
	for {
		s, err := SendRequest(target)
		if err != nil {
			log.Printf("Error: %v\n", err)
		} else {
			log.Println(s)
		}
		time.Sleep(interval * time.Millisecond)
	}
}

func main() {
    flag.Parse()
	rand.Seed(time.Now().UnixNano())
    target := addr + ":" + port
	for i := 1; i < worker_count; i++ {
		go worker(target, time.Duration(rand.Int()%max_sleep+1))
	}
	worker(target, time.Duration(rand.Int()%max_sleep+1))
}
