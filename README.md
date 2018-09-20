# firesocket

    通用Socket库
    
### Example


```go

package main

import (
	"github.com/dean2021/firesocket"
	"time"
	"dnsbrute/log"
	"fmt"
)

func main() {

	f := firesocket.New(&firesocket.Options{
		DNSCacheExpire: time.Minute * 5,
		Timeout:        time.Second * 10,
		WriteTimeout:   time.Second * 10,
		ReadTimeout:    time.Second * 10,
	})

	err := f.Connect("tcp", "www.jd.com", "80")

	if err != nil {
		log.Fatal(err)
	}

	f.Write([]byte("xxxxx"))


	b, err := f.ReadN(1000)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	f.Close()
}

```