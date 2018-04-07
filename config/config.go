package config

import (
  "../utilize"
)

var (
  Token string
  BotPrefix string
  BotName string
)



func ReadConfig() error {

  config, err := utilize.LoadConfig()
  if err != nil {
    return err
  }

  Token, err = utilize.LoadToken()
  if err != nil {
    return err
  }

  BotPrefix = config.BotPrefix
  BotName = config.BotName

  return nil
}
