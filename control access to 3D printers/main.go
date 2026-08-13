package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func main() {
	numPrinters := 3
	numhackers := 7
	hackers := make([]*hacker, numhackers)
	for i := range numhackers {
		hackers[i] = newHacker(i)
	}

	printers := make(chan struct{}, numPrinters)

	var wg sync.WaitGroup
	for i := range numhackers {
		wg.Add(1)
		go hackerFun(&wg, hackers[i], printers)
	}

	wg.Wait()

}

func hackerFun(wg *sync.WaitGroup, h *hacker, printer chan struct{}) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case printer <- struct{}{}:
			h.currentUsage++
			fmt.Printf("hacker %d on his %d usage , target %d\n", h.id, h.currentUsage, h.goalUsage)
			t := usePrinterTime()
			time.Sleep(t)
			// release printer
			<-printer
			// check if hacker is done
			ticker.Reset(5 * time.Second)
			if h.goalUsage == h.currentUsage {
				fmt.Printf("hacker %d finished with total %d usages\n", h.id, h.goalUsage)
				return
			}
		case <-ticker.C:
			fmt.Printf("hacker %d gave up\n", h.id)
			return
		}
	}

}

type hacker struct {
	id           int
	currentUsage int
	goalUsage    int
}

func newHacker(id int) *hacker {
	return &hacker{
		id:           id,
		currentUsage: 0,
		goalUsage:    hackerNumberOfUses(),
	}
}

func hackerNumberOfUses() int {
	number := 2 + rand.IntN(2)
	return number
}

func usePrinterTime() time.Duration {
	return time.Duration(1+rand.IntN(10)) * time.Second // 1..10s
}
