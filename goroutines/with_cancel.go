package main

import (
"context"
"fmt"
"sync"
)

func main() {
    wg := &sync.WaitGroup{}

    fmt.Println("--- start ---")

    // contextを生成
    ctx, cancel := context.WithCancel(context.Background())
    ch := make(chan int)
    for i := 0; i < 2; i++ {
        wg.Add(1)
        fmt.Println("i -----------", i)
        go goPrintNumbers(wg, ctx, ch)
    }

    for number := 1; ; {
        fmt.Println("number... ", number)
        // When reaching number to 16, stop goroutine
        if number == 16 {
            cancel() // finish ctx
            break
        }

        ch <- number
        number += 5
    }
    wg.Wait()

    fmt.Println("--- finish ---")
}

func goPrintNumbers(wg *sync.WaitGroup, ctx context.Context, ch chan int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("goroutine finish")
            wg.Done() // wgの数を一つ減らす
            return
        case num := <-ch:
            fmt.Println("num ---- ", num)
            for i := num; i < num+5; i++ {
                fmt.Println(i)
            }
        }
    }
}
