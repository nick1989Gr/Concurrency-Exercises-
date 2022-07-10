// During a busy hackaton 3D printers are heavily used. Write a simulation of hackers trying to access 3D printers.

// In the hackerspace, there are 3 3D printers. There are 7 hackers that are interested in using the printers.

// If the hacker can't access the printer for more than 5 seconds,
//he gets annoyed and quits the hackaton. Hackers use printers for random interval
//from 1 to 10 seconds and usually they need to use the printer at least twice, because nothing is perfect for the first time.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {

	printersCh := make(chan hacker)
	terminate := make(chan struct{}, 3)
	freeUser := make(chan struct{}, 3)
	// hackersCh := make(chan int, 7)
	hackers := []hacker{
		{id: 0}, {id: 1}, {id: 2}, {id: 3}, {id: 4}, {id: 5}, {id: 6},
	}
	var wgHackers sync.WaitGroup
	var wgPrinters sync.WaitGroup
	wgPrinters.Add(3)
	for i := 0; i < 3; i++ {

		go func() {
			defer wgPrinters.Done()
			for {
				select {
				case h := <-printersCh:
					h.using = true
					fmt.Printf("Hacker %v is using the printer\n", h)
					<-freeUser
				case <-terminate:
					fmt.Println("Printer is terminating")
					return
				}

			}

		}()
	}

	wgHackers.Add(7)
	for i := 0; i < 7; i++ {
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
					freeUser <- struct{}{}
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
	for i := 0; i < 3; i++ {
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
