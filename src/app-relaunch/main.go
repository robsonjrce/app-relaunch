package main

import (
  "fmt"
  "os"
  "time"
)

func watchFile(filePath string) error {
  initialStat, err := os.Stat(filePath)
  if err != nil {
    return err
  }

  for {
    stat, err := os.Stat(filePath)
    if err != nil {
      return err
    }

    if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
      break
    }

    time.Sleep(1 * time.Second)
  }

  return nil
}

func main() {
  doneChan := make(chan bool)

  go func(doneChan chan bool) {
    defer func() {
      doneChan <- true
    }()

    err := watchFile("script.sh")
    if err != nil {
      fmt.Println(err)
    }

    fmt.Println("File has been changed")
  }(doneChan)

  <-doneChan
}