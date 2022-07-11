package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	NUM_OF_HACKERS      = 10
	HACKER_FEE          = 100
	HACKERRANK_EXPENSES = 3000
)

func main() {

	payCh := make(chan int)
	endgame := make(chan struct{}, NUM_OF_HACKERS)
	var wg sync.WaitGroup
	wg.Add(NUM_OF_HACKERS)

	for i := 0; i < NUM_OF_HACKERS; i++ {
		i := i
		go func() {
			defer wg.Done()
			t := time.NewTicker(1 * time.Millisecond)
			dateToPay := GetRandomDate()
			today := 1
			for {
				select {
				case <-t.C:

					if dateToPay == today {
						// we need to pay
						fmt.Printf("Hacker %v: today= %v pays %v\n", i, today, HACKER_FEE)
						payCh <- HACKER_FEE
					}
					today = increaseDate(today)
				case <-endgame:
					fmt.Printf("User %v, stops using hacker rank due to bankruptcy\n", i)
					return
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		income := 0
		t := time.NewTicker(30 * time.Millisecond)
		for {
			select {
			case money := <-payCh:
				income += money
			case <-t.C:
				if income-HACKERRANK_EXPENSES < 0 {
					fmt.Println("bankruptcy !!!!\n ")
					for i := 0; i < NUM_OF_HACKERS; i++ {
						endgame <- struct{}{}
					}

					return
				} else {
					income -= HACKERRANK_EXPENSES
					fmt.Printf("Hackerrank paid %v in fees. Reamining income %v", HACKERRANK_EXPENSES, income)
				}
			}
		}
	}()

	wg.Wait()

}

func GetRandomDate() int {
	max := 30
	min := 1
	return rand.Intn(max-min) + min
}

func increaseDate(d int) int {
	d += 1
	if d > 30 {
		d = 1
	}
	return d
}
