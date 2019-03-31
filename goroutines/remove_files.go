package main

import (
    "fmt"
    "io/ioutil"
    "math/rand"
    "os"
    "sync"
    "time"
)

func main() {

    dirName := "./goroutines/dir"

    createFiles(dirName)

    removeFiles(dirName)

}

func createFiles(dirName string) {
    if !isExist(dirName) {
        os.Mkdir(dirName, 0777)
    }
    for i := 0; i < 10; i++ {
        num := randomNumber()
        fmt.Println(fmt.Sprintf("create file --- %s/file%d.txt", dirName, num))
        path := fmt.Sprintf("%s/file%d.txt", dirName, num)
        file, _ := os.OpenFile(fmt.Sprintf("%s/file%d.txt", dirName, num), os.O_RDWR|os.O_CREATE, 0777)
        defer file.Close()
        fmt.Fprintln(file, path)
    }
}

func isExist(name string) bool {
    _, err := os.Stat(name)
    return os.IsExist(err)
}

func removeFiles(dirName string) {
    wg := &sync.WaitGroup{}
    success := make(chan bool)

    files, err := ioutil.ReadDir(dirName)
    if err == nil {
        for i, file := range files {
            wg.Add(1)
            go removeFile(wg, success, i, file, dirName)
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

    /* Another process */
    // for {
    //     v, ok := <-success
    //     if !ok {
    //         break
    //     }
    //     if v {
    //         fmt.Println("End")
    //         count++
    //     } else {
    //         fmt.Println("Failed")
    //     }
    // }

    if count >= len(files) {
        deleteDir(dirName)
    }
}

func removeFile(wg *sync.WaitGroup, success chan bool, index int, file os.FileInfo, dirName string) {
    path := fmt.Sprintf("%s/%s", dirName, file.Name())
    // imitatePath := fmt.Sprintf("%s/file%d.txt", dirName, index)
    if _, err := os.Stat(path); err == nil {
        fmt.Println("path = ", path)
        if err := os.Remove(path); err != nil {
            // if err := os.Remove(imitatePath); err != nil {
            success <- false
        } else {
            success <- true
        }
        wg.Done()
    }
}

func deleteDir(dirName string) {
    fmt.Println(os.Remove(dirName))
}

func randomNumber() int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(10)
}
