/**
* @Author: HongBo Fu
* @Date: 2019/10/17 13:49
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var ns = flag.String("n", "", "-n 100 | -n 1000 | ...")
var cs = flag.String("c", "", "-c 10 | -c 100 | ...")
var url = flag.String("url", "", "-url http://localhost:8005/")
var ps = flag.String("p", "", "true|false")

func main() {
	flag.Parse()

	n, err := strconv.Atoi(*ns)
	if err != nil {
		fmt.Println("cant parse n", err, *ns)
		os.Exit(1)
	}
	c, err := strconv.Atoi(*cs)
	if err != nil {
		fmt.Println("cant parse c", err, *cs)
		os.Exit(1)
	}
	if *url == "" {
		fmt.Println("please input url")
		os.Exit(1)
	}
	fmt.Println(n, c, *url)

	tasks := make(chan bool, n)
	for i := 0; i < n; i++ {
		tasks <- true
	}

	begin := time.Now()
	cnt := 0
	cntChan := make(chan bool)

	for i := 0; i < c; i++ {
		go func() {
			for len(tasks) > 0 {
				<-tasks
				res, err := http.Get(*url)
				if err != nil {
					panic(err)
				} else {
					defer res.Body.Close()
					content, err := ioutil.ReadAll(res.Body)
					if err != nil {
						panic(err)
					} else {
						if *ps == "true" {
							fmt.Println(string(content))
						}
					}
				}

				cntChan <- true
			}
		}()
	}

	for {
		b := false
		select {
		case <-cntChan:
			cnt += 1
			if cnt%(n/10) == 0 {
				b = true
				fmt.Printf("finished %d request \n", cnt)
			}
			if cnt == n {
				if !b {
					fmt.Printf("finished %d request \n", cnt)
				}
				goto END
			}
		}
	}

END:
	ts := time.Now().Sub(begin).Seconds()
	fmt.Printf("%.3f requests/second \n", float64(cnt)/ts)
}
