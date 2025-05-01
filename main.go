package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

func main() {
  cmd := exec.Command("/bin/bash")
  pty, err := pty.Start(cmd)
  if err != nil {
    panic(err)
  }
  defer pty.Close()

  reader := bufio.NewReader(pty)
  go func() {
   for {
     buf := make([]byte, 1024)
     n, err := reader.Read(buf)
     if err != nil {
       if err == io.EOF {
         return
       }
       panic(err)
     }
     os.Stdout.Write(buf[:n])
   } 
  }()

  go func() {
    scanner := bufio.NewScanner(os.Stdin) 
    for scanner.Scan(){
      input := scanner.Text()
      pty.Write([]byte(input + "\n"))
    }
  }()

  time.Sleep(10 * time.Second)
}
