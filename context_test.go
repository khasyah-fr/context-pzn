package main

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextC, "e", "E")

	contextF := context.WithValue(contextE, "f", "F")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)

	fmt.Println(contextF.Value("f"))
	fmt.Println(contextE.Value("c"))
	fmt.Println(contextF.Value("b"))
}

func CreateCounter(ctx context.Context) chan int {
	dest := make(chan int)

	go func() {
		defer close(dest)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				dest <- counter
				counter++
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return dest
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println(n)
		if n == 10 {
			break
		}
	}

	time.Sleep(10 * time.Second)
	cancel()
	time.Sleep(2 * time.Second)

	fmt.Println(runtime.NumGoroutine())
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println(n)
		if n == 10 {
			break
		}
	}

	fmt.Println(runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(7*time.Second))
	defer cancel()

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println(n)
		if n == 10 {
			break
		}
	}

	fmt.Println(runtime.NumGoroutine())
}
