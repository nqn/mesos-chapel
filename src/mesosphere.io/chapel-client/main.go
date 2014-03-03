package main

import (
  "fmt"
  "bufio"
  "net"
  "log"
  "strings"
  "os"
)

func main() {
  if len(os.Args) < 4 {
    fmt.Println("usage <client file> <client> <target> <arguments>")
    os.Exit(1)
  }

  target := os.Args[2]
  selected_hostname := os.Args[3]
  client := ""

  f, _ := os.Open(os.Args[1])
  r := bufio.NewReader(f)
  for {
    line, _, err := r.ReadLine()
    if err != nil {
      break
    }

    str := strings.Trim(string(line), "\r\n")
    split := strings.Split(str, ":")

    hostname := split[0]
    if hostname == selected_hostname {
      client = str
      break
    }
  }

  if client == "" {
    fmt.Println("Agent not found")
    os.Exit(1)
  }

  conn, err := net.Dial("tcp", client)
  if err != nil {
   log.Fatal(err)
  }

  pwd, _ := os.Getwd()
  buf := []byte(pwd + "?" + target + "? mkdir -p " + pwd + ";"  + strings.Join(os.Args[4:], " "))
  _, err = conn.Write(buf)
  if err != nil {
   log.Fatal(err)
  }

  os.Exit(0)
}
