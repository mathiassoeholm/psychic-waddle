package main

import (
	"fmt"
	"sync"
)

type Person struct {
	Name string
	Age  int
}

type Nameable interface {
	talk()
	changeName(name string)
}

func (p Person) talk() {
	fmt.Println("Hi, my name is ", p.Name)
}

func (p *Person) changeName(newName string) {
	p.Name = newName
}

func main() {
	p := Person{Name: "John", Age: 30}
	channel := make(chan Person, 1)
	channel <- p

	go func() {
		for {
			channel <- p
		}
	}()

	i := 0
	for {
		if i == 10 {
			break
		}
		i++

		select {
		case p, ok := <-channel:
			if !ok {
				panic("Channel closed")
			}
			p.talk()
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { p.talk(); wg.Done() }()
	wg.Wait()
	p.talk()
	p.changeName("Jenny")
	p.talk()
}
