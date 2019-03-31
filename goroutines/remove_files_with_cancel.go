package main

import (
    "context"
    "fmt"
    "io/ioutil"
    "math/rand"
    "os"
    "sync"
    "time"
)

func main() {

    dirName := "./goroutines/dir"

    createFs(dirName)

    removeFilesWithCancel(dirName)

}

type CancelContext struct {
    ctx    context.Context
    cancel context.CancelFunc
}

func removeFilesWithCancel(dirName string) {
    ctx, cancel := context.WithCancel(context.Background())
    cancelCtx := &CancelContext{ctx, cancel}
    wg := &sync.WaitGroup{}
    success := make(chan bool)

    files, err := ioutil.ReadDir(dirName)
    if err == nil {
        for i, file := range files {
            wg.Add(1)
            go removeFileWithCacel(wg, success, cancelCtx, i, file, dirName)
        }
    }

    go func() {
        wg.Wait()
        close(success)
    }()

    count := 0
    for i := range files {
        fmt.Println("for: ", i)
        select {
        case <-cancelCtx.ctx.Done():
            fmt.Println("Canceled since an error occured")
            break
        case v, ok := <-success:
            if !ok {
                break
            }
            if v {
                fmt.Println("End")
                count++
            } else {
                fmt.Println("Failed")
            }
        }
    }
    if count >= len(files) {
        deleteDr(dirName)
    }
}

func removeFileWithCacel(wg *sync.WaitGroup, success chan bool, cancelCtx *CancelContext, index int, file os.FileInfo, dirName string) {
    path := fmt.Sprintf("%s/%s", dirName, file.Name())
    // imitatePath := fmt.Sprintf("%s/file%d.txt", dirName, index)
    if _, err := os.Stat(path); err == nil {
        fmt.Println("path = ", path)
        if err := os.Remove(path); err != nil {
            cancelCtx.cancel()
            success <- false
        } else {
            success <- true
        }
        wg.Done()
    }
}

func createFs(dirName string) {
    if !exists(dirName) {
        os.Mkdir(dirName, 0777)
    }
    for i := 0; i < 10; i++ {
        num := randomNum()
        fmt.Println(fmt.Sprintf("create file --- %s/file%d.txt", dirName, num))
        path := fmt.Sprintf("%s/file%d.txt", dirName, num)
        file, _ := os.OpenFile(fmt.Sprintf("%s/file%d.txt", dirName, num), os.O_RDWR|os.O_CREATE, 0777)
        defer file.Close()
        fmt.Fprintln(file, path)
    }
}

func exists(name string) bool {
    _, err := os.Stat(name)
    return os.IsExist(err)
}

func deleteDr(dirName string) {
    fmt.Println(os.Remove(dirName))
}

func randomNum() int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(10)
}
