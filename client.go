package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc/jsonrpc"
)

type Args struct {
	Arg1, Arg2 string
}

func main() {
	for {
		client, err := net.Dial("tcp", "127.0.0.1:1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		var text, reply string
        args := new(Args)
		fmt.Printf("Enter command:\n")
		fmt.Scan(&text)
		switch text {
		case "new":
			fmt.Printf("Enter name:\n")
			fmt.Scan(&args.Arg1)
			fmt.Printf("Enter value:\n")
			fmt.Scan(&args.Arg2)
            fmt.Printf("%v\n", args)
			c := jsonrpc.NewClient(client)
			err = c.Call("Computation.New", args, &reply)
			if err != nil {
                fmt.Println(reply)
				log.Fatal("Error:", err)
			} else {
				fmt.Printf("Result: Create %s with value %s\n", args.Arg1, args.Arg2)
			}
		case "set":
			fmt.Printf("Enter name:\n")
			fmt.Scan(&args.Arg1)
			fmt.Printf("Enter value:\n")
			fmt.Scan(&args.Arg2)
			c := jsonrpc.NewClient(client)
			err = c.Call("Computation.Set", args, &reply)
			if err != nil {
				log.Fatal("Error:", err)
			} else {
				fmt.Printf("Result: Set %s with value %s\n", args.Arg1, args.Arg2)
			}
		case "del":
			fmt.Printf("Enter name:\n")
			fmt.Scan(&args.Arg1)
			c := jsonrpc.NewClient(client)
			err = c.Call("Computation.Del", args, &reply)
			if err != nil {
				log.Fatal("Error:", err)
			} else {
				fmt.Printf("Result: Delete %s\n", args.Arg1)
			}
		case "add":
			fmt.Printf("Enter name1:\n")
			fmt.Scan(&args.Arg1)
			fmt.Printf("Enter name2:\n")
			fmt.Scan(&args.Arg2)
			c := jsonrpc.NewClient(client)
			err = c.Call("Computation.Add", args, &reply)
			if err != nil {
				log.Fatal("arith error:", err)
			} else {
				fmt.Printf("Result: %s + %s = %s\n", args.Arg1, args.Arg2, reply)
			}
		case "sub":
			fmt.Printf("Enter name1:\n")
			fmt.Scan(&args.Arg1)
			fmt.Printf("Enter name2:\n")
			fmt.Scan(&args.Arg2)
			c := jsonrpc.NewClient(client)
			err = c.Call("Computation.Sub", args, &reply)
			if err != nil {
				log.Fatal("Error:", err)
			} else {
				fmt.Printf("Result: %s - %s = %s\n", args.Arg1, args.Arg2, reply)
			}
		case "mul":
			fmt.Printf("Enter name1:\n")
			fmt.Scan(&args.Arg1)
			fmt.Printf("Enter name2:\n")
			fmt.Scan(&args.Arg2)
			c := jsonrpc.NewClient(client)
			err = c.Call("Computation.Mul", args, &reply)
			if err != nil {
				log.Fatal("arith error:", err)
			} else {
				fmt.Printf("Result: %s * %s = %s\n", args.Arg1, args.Arg2, reply)
			}
		case "div":
			fmt.Printf("Enter name1:\n")
			fmt.Scan(&args.Arg1)
			fmt.Printf("Enter name2:\n")
			fmt.Scan(&args.Arg2)
			c := jsonrpc.NewClient(client)
			err = c.Call("Computation.Div", args, &reply)
			if err != nil {
				log.Fatal("Error:", err)
			} else {
				fmt.Printf("Result: %s / %s = %s\n", args.Arg1, args.Arg2, reply)
			}
		case "quit":
			return
		}
		client.Close()
	}
}
