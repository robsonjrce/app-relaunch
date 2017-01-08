package main

import (
  "bytes"
  "bufio"
  "fmt"
  "log"
  "os"
  "os/exec"
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

func runCommandWithOutputStream(command string) {
  cmd := exec.Command(command)
  cmdReader, err := cmd.StdoutPipe()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
    os.Exit(1)
  }

  scanner := bufio.NewScanner(cmdReader)
  go func() {
    for scanner.Scan() {
      fmt.Printf("%s | %s\n", command, scanner.Text())
    }
  }()

  err = cmd.Start()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
    os.Exit(1)
  }

  err = cmd.Wait()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
    os.Exit(1)
  }
}

func runCommand(command string) {
  cmd := exec.Command(command)
  // cmd.Stdin = strings.NewReader("some input")
  var out bytes.Buffer
  cmd.Stdout = &out
  err := cmd.Run()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("in all caps: %q\n", out.String())
}

func main() {
  args := os.Args[1:]
  filename := args[0]

  if _, err := os.Stat(filename); os.IsNotExist(err) {
    fmt.Println("file doesn't exist")
  } else {
    runCommandWithOutputStream(filename)
    // doneChan := make(chan bool)

    // go func(doneChan chan bool) {
    //   defer func() {
    //     doneChan <- true
    //   }()

    //   err := watchFile(filename)
    //   if err != nil {
    //     fmt.Println(err)
    //   }

    //   fmt.Println("File has been changed")
    // }(doneChan)

    // <-doneChan
  }
}