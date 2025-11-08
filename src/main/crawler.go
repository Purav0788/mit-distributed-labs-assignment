// package main

// import (
// 	"fmt"
// 	"sync"
// )

// type Fetcher interface {
// 	// Fetch returns the body of URL and
// 	// a slice of URLs found on that page.
// 	Fetch(url string) (body string, urls []string, err error)
// }

// // Crawl uses fetcher to concurrently crawl
// // pages starting with url.
// func Crawl(url string, fetcher Fetcher) {
// 	// TODO: Fill in this function.
// 	//
// 	// This implementation is sequential and doesn't track visited URLs.
// 	// Your mission is to make it concurrent and to fetch each URL only once.
// 	//
// 	// HINTS:
// 	// 1. You will need to keep track of URLs that have been fetched.
// 	//    A map is a good choice for this.
// 	// 2. Since multiple goroutines will be accessing this map, it will need
// 	//    to be protected by a Mutex to prevent data races.
// 	// 3. You will need a way for the main 'Crawl' function to know when all
// 	//    the concurrent crawling is finished. A sync.WaitGroup is the
// 	//    perfect tool for this.

// 	// fmt.Println("Starting crawl...") // You can remove this line.

// 	// // A simple, sequential (and incorrect) implementation to get you started:
// 	// body, urls, err := fetcher.Fetch(url)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// 	return
// 	// }
// 	// fmt.Printf("found: %s %q\n", url, body)
// 	// for _, u := range urls {
// 	// 	// This will fetch URLs repeatedly and sequentially.
// 	// 	// Your job is to fix this!
// 	// }
	
// 	//iterate on the new url, if not iterated already(visited)
// 	//print the url when trying to iterate upon it
// 	//add its urls to the queue if not already visisted before.
// 	//fill the iterator queue ( that is the tasks)
// 	//maintain the list of tasks already done ( visited, if visited dont add to iterator queue)
// 	var mu sync.mutex
// 	visited.insert(url);
// 	var wg sync.WaitGroup
// 	wg.Add(1);
// 	go iterate(url, &wg);
// 	wg.Wait()
// 	return
// }

// func iterate(url string, wg *sync.WaitGroup) {
// 	defer wg.Done();
// 	if(fetcher.find(url) != fetcher.end()){
// 		for range fetcher[url]{
// 				urlchild = fetcher[url]
// 				//this map would have been an internet call, so this would be consuming time
// 				//new mutex.lock()
// 				//this looks wrong, it has to be a global mutex lock I think rather than just new one being defined 
// 				// anew with each iteration of the for loop
// 				// defer mutex.unlock()
// 				mu.lock()
// 				if(visited.find(urlchild) == visited.end()){
// 					mu.unlock()
// 				} else{
// 					visited.insert(urlchild);
// 					fmt.Print("urlfound", url);
// 					wg.Add(1);
// 					go iterate(url, wg);
// 					mu.unlock();
// 				}
			
// 			}
// 		} else{
// 			mu.lock()
// 			visited.insert(urlchild)
// 			fmt.Print("url not found")
// 			mu.unlock()
// 		}
// 		return

// }
// func main() {
// 	Crawl("https://golang.org/", fetcher)
// }

// // --- Fake Fetcher for Testing (Do not modify) ---

// // fakeFetcher is Fetcher that returns canned results.
// type fakeFetcher map[string]*fakeResult

// type fakeResult struct {
// 	body string
// 	urls []string
// }

// func (f fakeFetcher) Fetch(url string) (string, []string, error) {
// 	if res, ok := f[url]; ok {
// 		fmt.Printf("found: %s %q\n", url, res.body)
// 		return res.body, res.urls, nil
// 	}
// 	fmt.Printf("not found: %s\n", url)
// 	return "", nil, fmt.Errorf("not found: %s", url)
// }

// // fetcher is a populated fakeFetcher.
// var fetcher = fakeFetcher{
// 	"https://golang.org/": &fakeResult{
// 		"The Go Programming Language",
// 		[]string{
// 			"https://golang.org/pkg/",
// 			"https://golang.org/cmd/",
// 		},
// 	},
// 	"https://golang.org/pkg/": &fakeResult{
// 		"Packages",
// 		[]string{
// 			"https://golang.org/",
// 			"https://golang.org/cmd/",
// 			"https://golang.org/pkg/fmt/",
// 			"https://golang.org/pkg/os/",
// 		},
// 	},
// 	"https://golang.org/pkg/fmt/": &fakeResult{
// 		"Package fmt",
// 		[]string{
// 			"https://golang.org/",
// 			"https://golang.org/pkg/",
// 		},
// 	},
// 	"https://golang.org/pkg/os/": &fakeResult{
// 		"Package os",
// 		[]string{
// 			"https://golang.org/",
// 			"https://golang.org/pkg/",
// 		},
// 	},
// }