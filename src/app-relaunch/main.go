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

func runCommandWithOutputStream(command string) (*exec.Cmd, error) {
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

  return cmd, err

  // err = cmd.Wait()
  // if err != nil {
  //   fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
  //   os.Exit(1)
  // }
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
    if cmd, err := runCommandWithOutputStream(filename); err == nil {
      // terminating program
      //
      done := make(chan error, 1)
      go func() {
        done <- cmd.Wait()
      }()

      select {
        case <-time.After(3 * time.Second):
          if err := cmd.Process.Kill(); err != nil {
            log.Fatal("failed to kill: ", err)
          }
          log.Println("process killed as timeout reached")

        case err := <-done:
          if err != nil {
            log.Printf("process done with error = %v", err)
          } else {
            log.Print("process done gracefully without error")
          }
      }
    }

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