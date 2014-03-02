package main

import (
  "fmt"
  "net"
  "strings"
  "os"
  "os/exec"
  "io"
  "flag"
  "strconv"
  "log"
)

func main() {
  port := flag.Uint64("port", 5000, "Agent port to listen for scheduler")
  flag.Parse()

  portStr := strconv.FormatUint(*port, 10)

  ln, err := net.Listen("tcp", ":" + portStr)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Listening on port: " + portStr)

  conn, err := ln.Accept()
  if err != nil {
    log.Fatal(err)
  }
  buf := make([]byte, 1024)
  n, err := conn.Read(buf)
  if err != nil {
    log.Fatal(err)
  }

  cmds := []string { "sh", "-c", string(buf[:n]) }

  fmt.Println("Cmd: " + strings.Join(cmds, " "))

  cmd := exec.Command(cmds[0], cmds[1:]...)
  cmd.Env = os.Environ()

  stdout, _ := cmd.StdoutPipe()
  go io.Copy(os.Stdout, stdout)

  stderr, _ := cmd.StderrPipe()
  go io.Copy(os.Stderr, stderr)

  if err := cmd.Run(); err != nil {
    log.Fatal(err)
  }
  cmd.Wait()
}
