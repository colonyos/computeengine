package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/colonyos/colonies/pkg/client"
	"github.com/colonyos/colonies/pkg/core"
	"github.com/colonyos/colonies/pkg/security/crypto"

	fib "github.com/t-pwk/go-fibonacci"
)

func main() {
	args := os.Args
	unregister := false
	if len(args) > 1 {
		unregister = true
	}

	if unregister {
		colonyPrvKey := os.Getenv("COLONYPRVKEY")
		host := os.Getenv("COLONIES_SERVER_HOST")
		portStr := os.Getenv("COLONIES_SERVER_PORT")

		port, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Println(err)
			return
		}

		client := client.CreateColoniesClient(host, port, true)

		runtimeIDBytes, err := os.ReadFile("/tmp/runtimeid")
		if err != nil {
			fmt.Println(err)
		}

		runtimeID := string(runtimeIDBytes)

		fmt.Println("Unregister Runtime with ID <" + runtimeID + ">")

		err = client.DeleteRuntime(runtimeID, colonyPrvKey)
		if err != nil {
			fmt.Println(err)
		}

		return
	}

	colonyID := os.Getenv("COLONYID")
	colonyPrvKey := os.Getenv("COLONYPRVKEY")
	host := os.Getenv("COLONIES_SERVER_HOST")
	portStr := os.Getenv("COLONIES_SERVER_PORT")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := client.CreateColoniesClient(host, port, true)

	crypto := crypto.CreateCrypto()
	runtimePrvKey, err := crypto.GeneratePrivateKey()
	if err != nil {
		fmt.Println(err)
	}

	runtimeID, err := crypto.GenerateID(runtimePrvKey)
	if err != nil {
		fmt.Println(err)
	}

	err = os.WriteFile("/tmp/runtimeid", []byte(runtimeID), 0644)
	if err != nil {
		fmt.Println(err)
	}

	runtimeType := "fibonacci"
	name := "fibonacci"
	cpu := "AMD Ryzen 9 5950X (32) @ 3.400GHz"
	coresStr := os.Getenv("CORES")
	fmt.Println(coresStr)
	cores, err := strconv.Atoi(coresStr)
	if err != nil {
		fmt.Println(err)
	}

	memStr := os.Getenv("MEM")
	fmt.Println(memStr)
	mem, err := strconv.Atoi(memStr)
	if err != nil {
		fmt.Println(err)
	}

	gpu := ""
	gpus := 0

	fmt.Println(name)
	fmt.Println(runtimeType)

	runtime := core.CreateRuntime(runtimeID, runtimeType, name, colonyID, cpu, cores, mem, gpu, gpus)
	_, err = client.AddRuntime(runtime, colonyPrvKey)
	if err != nil {
		fmt.Println(err)
	}

	err = client.ApproveRuntime(runtimeID, colonyPrvKey)
	if err != nil {
		fmt.Println(err)
	}

	for {
		assignedProcess, err := client.AssignProcess(colonyID, runtimePrvKey)
		if err != nil {
			fmt.Println(err)
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		// Parse env attribute and calculate the given Fibonacci number
		for _, attribute := range assignedProcess.Attributes {
			if attribute.Key == "fibonacciNum" {
				nr, _ := strconv.Atoi(attribute.Value)
				fibonacci := fib.FibonacciBig(uint(nr))

				min := 100   // 0.1 s
				max := 40000 // 40s
				sleepTime := rand.Intn(max-min+1) + min

				fmt.Printf("sleeping for %d\n", sleepTime)

				time.Sleep(time.Duration(sleepTime) * time.Millisecond)

				attribute := core.CreateAttribute(assignedProcess.ID, core.OUT, "result", fibonacci.String())
				client.AddAttribute(attribute, runtimePrvKey)

				// Close the process as Successful
				client.CloseSuccessful(assignedProcess.ID, runtimePrvKey)
			}
		}
	}
}
