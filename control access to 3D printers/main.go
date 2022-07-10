package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	NUMBER_OF_HACKERS  = 7
	NUMBER_OF_PRINTERS = 3
)

func main() {

	printersCh := make(chan hacker)
	terminate := make(chan struct{}, NUMBER_OF_PRINTERS)
	freePrinter := make(chan struct{}, NUMBER_OF_PRINTERS)

	hackers := []hacker{
		{id: 0}, {id: 1}, {id: 2}, {id: 3}, {id: 4}, {id: 5}, {id: 6},
	}
	var wgHackers sync.WaitGroup
	var wgPrinters sync.WaitGroup
	wgPrinters.Add(NUMBER_OF_PRINTERS)
	for i := 0; i < NUMBER_OF_PRINTERS; i++ {

		go func() {
			defer wgPrinters.Done()
			for {
				select {
				case h := <-printersCh:
					h.using = true
					fmt.Printf("Hacker %v is using the printer\n", h.id)
					<-freePrinter
				case <-terminate:
					fmt.Println("Printer is terminating")
					return
				}

			}

		}()
	}

	wgHackers.Add(NUMBER_OF_HACKERS)
	for i := 0; i < NUMBER_OF_HACKERS; i++ {
		i := i
		go func() {
			defer wgHackers.Done()
			for {
				select {
				// case that our request goes through
				case printersCh <- hackers[i]:
					requestedTime := getUseTime()
					hackers[i].timesUsed += 1
					fmt.Printf("Hacker with id : %v has requests %v seconds time: %v\n", hackers[i].id, requestedTime, hackers[i].timesUsed)
					t := time.NewTimer(time.Duration(requestedTime) * time.Second)
					// Block until the timer is triggered
					<-t.C
					// Now we can release one printer
					fmt.Printf("Hacker %v releasing a printer\n", hackers[i].id)

					hackers[i].using = false
					freePrinter <- struct{}{}
					if hackers[i].timesUsed == 2 {
						fmt.Printf("Hacker with id : %v has used the printer 2 times\n", hackers[i].id)
						return
					}
				case <-time.After(5 * time.Second):
					if !hackers[i].using {
						fmt.Printf("Hacker with id : %v quited\n", hackers[i].id)
						return
					}
				}
			}

		}()
	}

	wgHackers.Wait()
	for i := 0; i < NUMBER_OF_PRINTERS; i++ {
		terminate <- struct{}{}
	}
	wgPrinters.Wait()

}

type hacker struct {
	timesUsed int
	id        int
	using     bool
}

func getUseTime() int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 10
	return rand.Intn(max-min+1) + min

}
