package belajar_golang_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	backgroud := context.Background()
	fmt.Println(backgroud)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "f")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)

	fmt.Println(contextF.Value("f"))
	fmt.Println(contextF.Value("c"))
	fmt.Println(contextF.Value("b"))

	fmt.Println(contextA.Value("b"))
}

func CreateCounterLeak() chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			destination <- counter
			counter++
		}
	}()

	return destination
}

func TestContextLeak(t *testing.T) {
	fmt.Println("Total goroutine start:", runtime.NumGoroutine())

	destination := CreateCounterLeak()

	fmt.Println("Total goroutine after call goroutine:", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)

		if n == 10 {
			break
		}
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total goroutine after cancel signal event:", runtime.NumGoroutine())
}

func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
			}
		}
	}()

	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total goroutine start:", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)

	fmt.Println("Total goroutine after call goroutine:", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)

		if n == 10 {
			break
		}
	}
	cancel() // mengirim sinyal cancel ke context

	time.Sleep(2 * time.Second)

	fmt.Println("Total goroutine after cancel signal event:", runtime.NumGoroutine())
}

func CreateCounterSlow(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++

				time.Sleep(1 * time.Second)
			}
		}
	}()

	return destination
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total goroutine start:", runtime.NumGoroutine())

	parent := context.Background()

	// timeoutTime := 5 * time.Second
	timeoutTime := 10 * time.Second

	ctx, cancel := context.WithTimeout(parent, timeoutTime)
	defer cancel() // keep call cancel for automatic stop go routine less than 5 second of timeout

	destination := CreateCounterSlow(ctx)

	fmt.Println("Total goroutine after call goroutine:", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total goroutine after cancel signal event:", runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total goroutine start:", runtime.NumGoroutine())

	parent := context.Background()

	timeoutTime := 5 * time.Second
	deadlineTime := time.Now().Add(timeoutTime)

	ctx, cancel := context.WithDeadline(parent, deadlineTime)
	defer cancel() // keep call cancel for automatic stop go routine less than 5 second of timeout

	destination := CreateCounterSlow(ctx)

	fmt.Println("Total goroutine after call goroutine:", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total goroutine after cancel signal event:", runtime.NumGoroutine())
}
