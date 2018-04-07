package main

import (
  "fmt"

  "os"
  "os/signal"
	"syscall"

  "./config"
  "./bot"

)


func main()  {

  err := config.ReadConfig()

  if err != nil {
    fmt.Println(err.Error())
    return
  }

  bot.Start()

  sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

  bot.Stop()

}
