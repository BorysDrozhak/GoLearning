package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	// "net"
	"github.com/bogdanovich/dns_resolver"
	"os"
	"regexp"
	"sync"
	"time"
)

var limit chan bool

func unresolved(e string, rc *int) {
	fmt.Fprintln(os.Stderr, e)
	*rc = 255 // set error code in case we do not resolve at least once.
}

func panicIf(e error) {
	if e != nil {
		panic(e)
	}
}

func resolve_file(resolver_host string, fileName string, wg *sync.WaitGroup, ch chan string, rc *int) {
	// fun get a file and then resolver all A records there

	// set resolvers
	resolver := dns_resolver.New([]string{resolver_host})
	resolver.RetryTimes = 5

	file, e := os.Open(fileName) // open file
	panicIf(e)
	defer file.Close()

	Readfile := bufio.NewReader(file)

	// go read lines!
	for {

		s, _, e := Readfile.ReadLine()

		if e != nil {
			if e == io.EOF {
				break
			} else {
				panic(e)
			}
		}

		wg.Add(1) // increment
		go func(s string) {
			// v, e := net.LookupHost(s)
			limit <- true
			v, e := resolver.LookupHost(s)
			if e != nil {
				unresolved(s+" "+e.Error(), rc)
			} else {
				ch <- fmt.Sprint(s, " -> ", v)
			}
			<-limit
			wg.Done()
		}(string(s))
	}
	wg.Done() // decrement for  external
}

func main() {

	var rc int = 0
	begin := time.Now()

	var lists_directory, resolver_host, prefix string
	var limitation int
	flag.StringVar(&lists_directory, "d", "./A_lists", "directory of files with A zones")
	flag.StringVar(&resolver_host, "resolver", "8.8.8.8", "a resolver")
	flag.StringVar(&prefix, "DC", "", `prefix for choosing DC`) // use as prefix for choosing files from dirs
	flag.IntVar(&limitation, "limit", 100, `limitation for conc lookups`)

	flag.Parse()

	prefix = prefix + ".*"

	var wg sync.WaitGroup
	ch := make(chan string, limitation)
	limit = make(chan bool, limitation)
	done := make(chan bool) // logic channel

	// read channel
	go func() {
		for r := range ch {
			fmt.Println(r)
		}
		done <- true
	}()

	// read files in the dir
	dir, e := os.Open(lists_directory)
	panicIf(e)
	fileNames, e := dir.Readdirnames(-1)
	panicIf(e)
	defer dir.Close()

	rex, err := regexp.Compile(prefix) // compile
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	// go through all of them
	for _, fileName := range fileNames {

		// will use only ones we want by matching prefix in the FileName
		if m := rex.Match([]byte(fileName)); !m {
			continue
		}
		// the routine magic
		wg.Add(1)
		go resolve_file(resolver_host, lists_directory+"/"+fileName, &wg, ch, &rc)
	}

	wg.Wait()
	close(ch)
	<-done
	fmt.Println("it took: ", time.Since(begin))
	os.Exit(rc)
}
