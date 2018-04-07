package handlers

import (
  "fmt"
  "strings"

  "../utilize"
  "../utilize/minecraft"
)

var flg bool = false

func PingMinecraftServer(content string) string  {

  if flg {
    fmt.Println("PingMinecraftServer() start")
    defer fmt.Println("\n")
  }

  var (
    name string
    host string
    result string
  )

  content = strings.TrimPrefix(content, " ")

  //if arguments is empty then load server from file by default
  if content == "" {

    servers, err := utilize.ServersLoad();
    if err != nil {
      result = servers[0]
      return result
    }


    //ping all servers
    for _, server := range(servers) {

        serv := strings.SplitN(server, ":", 2)
        name = serv[0]
        host = serv[1]

        result += minecraft.Ping(name, host) + "\n\n"
    }


  } else {
    name = strings.Split(content, ":")[0]
    host = content

    result = minecraft.Ping(name, host)
  }

  if flg {
    fmt.Println("name:", name)
    fmt.Println("content:", content)
  }

  return result
}
