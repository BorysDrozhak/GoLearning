package main

import (
	"io/ioutil"
	"log"
	"sync"
)

func writeComplete()  {
	log.Println("wow")
	fileName := "/Users/bdrozhak/.bashrc"
	resultFile := "/Users/bdrozhak/TEST"
	_, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		//os.Exit(255)
	} else {
		log.Println("I did it")
		ioutil.WriteFile(resultFile,[]byte("I did it"), 777)
		//fmt.Printf("file content: %s", r)
	}
	//wg.Done()
}

func writeCompleteSync(wg *sync.WaitGroup)  {
	writeComplete()
	wg.Done()
}


func main() {
	//ch := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go writeComplete(&wg)

	// SHOULD i really pass the wg to the func?

//go func(file string) {
	//	//fmt.Println(file)
	//	ioutil.WriteFile(resultFile,[]byte("I did it"), 777)
	//	ch <- true
	//}(resultFile)
	//_ = <- ch
	wg.Wait()
}

