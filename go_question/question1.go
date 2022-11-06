package main

import (
	"fmt"
	"sync"
)

func main() {
	numChan := make(chan struct{})  //打印数字的通知
	charChan := make(chan struct{}) //打印字母的通知
	waitG := sync.WaitGroup{}
	waitG.Add(2)
	go func() {
		i := 1
		for _ = range numChan {
			fmt.Printf("%d%d", i, i+1)
			i = i + 2
			charChan <- struct{}{}
		}
		close(charChan)
		waitG.Done()
	}()
	go func() {
		c := 'A'
		for _ = range charChan {
			if c > 'Z' {
				break
			}
			fmt.Printf("%c%c", c, c+1)
			c = c + 2
			numChan <- struct{}{}
		}
		close(numChan)
		waitG.Done()
	}()
	numChan <- struct{}{}
	waitG.Wait()
}
