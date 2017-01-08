package main

import (
  "bytes"
  "bufio"
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
    log.Println(os.Stderr, "Error creating StdoutPipe for Cmd", err)
    os.Exit(1)
  }

  scanner := bufio.NewScanner(cmdReader)
  go func() {
    for scanner.Scan() {
      log.Println(scanner.Text())
    }
  }()

  err = cmd.Start()
  if err != nil {
    log.Fatal(os.Stderr, "Error starting Cmd", err)
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
  log.Printf("in all caps: ", out.String())
}

func main() {
  args := os.Args[1:]
  filename := args[0]

  if _, err := os.Stat(filename); os.IsNotExist(err) {
    log.Fatal("file doesn't exist")
  } else {
    for {
      if cmd, err := runCommandWithOutputStream(filename); err == nil {
        // terminating program
        //
        killChan := make(chan bool)
        go func(doneChan chan bool) {
          defer func() {
            doneChan <- true
          }()

          err := watchFile(filename)
          if err != nil {
            log.Fatal(err)
          }

          log.Println("file has changed")
        }(killChan)

        doneChan := make(chan error, 1)
        go func() {
          doneChan <- cmd.Wait()
        }()

        select {
          case <-killChan:
            if err := cmd.Process.Kill(); err != nil {
              log.Fatal("failed to kill: ", err)
            }
            log.Println("")
            log.Println("=========================================")
            log.Println("process killed as new file was identified")
            log.Println("=========================================")
            log.Println("")

          case err := <-doneChan:
            if err != nil {
              log.Printf("process done with error = %v", err)
            } else {
              log.Print("process done gracefully without error")
            }
        }
      }
    }
  }
}