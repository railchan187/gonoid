package minecraft

import (
  "fmt"
  "github.com/andrewtian/minepong"
  "net"
  "strconv"
  "time"
)


var flg bool = false

type status struct {
  name string

  online string
  max string
}

type server struct {
  name string
  host string

  conn net.Conn
}

func newServer(name string, host string) *server  {
  return &server {
    name: name,
    host: host,
  }
}

func (s *server) connect() error  {

  var err error
  s.conn, err = net.Dial("tcp", s.host)
  if err != nil {
    return err
  }

  return nil
}

func (s *server) disconnect() error {
  return s.conn.Close()
}



func Ping(name string, host string) string  {

  if flg {
    fmt.Println("Ping() start")
    defer fmt.Println("\n")
  }

  var result string

  if stat, ok := getStatus(name, host); ok {
    result = stat.name + "\n" +
            "online: " + stat.online + "/" + stat.max
  }else {
    result = "Server " + stat.name + " is down"
  }

  if flg {
    fmt.Println("result", result)
  }

  return result
}



func getStatus(name string, host string) (*status, bool) {
  if flg {
    fmt.Println("getStatus() start")
    defer fmt.Println("\n\n")
  }

  timeoutcha := make(chan bool, 1)
  connectcha := make(chan bool, 1)

  stat := new(status)
  stat.name = name

  svr := newServer(name, host)

  //create connecttion to server
  go func ()  {
    if err := svr.connect(); err != nil {
      fmt.Println("Connect error")
      connectcha <- false
    }else{
      connectcha <- true
    }
  }()

  //connection waiting
  go func()  {
    time.Sleep(3 * time.Second)
    timeoutcha <- true
  }()

  //handle the chans with results
  select{
  case res := <- connectcha:
    //if connect error
    if res == false {
      return stat, false
    }
    break

  //if connect so long
  case <-timeoutcha:
    return stat, false
  }

  //ping with connection
  pong, err := minepong.Ping(svr.conn, svr.host)
  if err != nil {
    return stat, false
  }

  stat.online = strconv.FormatInt(int64(pong.Players.Online), 10)
  stat.max = strconv.FormatInt(int64(pong.Players.Max), 10)

  return stat, true

}
