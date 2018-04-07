package utilize

import (
  "fmt"
  "io/ioutil"
  "encoding/json"
  "strings"
  "errors"
)

func LoadToken(params... string) (string, error) {
  var path string
  var token string

  if len(params) == 0 {
    path = "./token"
  } else {
    path = params[0]
  }

  b_token, err := ioutil.ReadFile(path) //just bytes

  if err != nil {
    fmt.Println(err)
    token = ""
    return token, errors.New("Cannot read token file")
  }

  token = string(b_token) //convert bytes to string
  token = strings.Trim(token, "\n") //remove "\n" in token at end

  return token, nil
}

type configStruct struct {
  BotPrefix string `json:"BotPrefix"`
  BotName   string `json:"BotName"`
}

func LoadConfig() (*configStruct, error)  {
    var config *configStruct

    fmt.Println("Reading from config file...")

    file, err := ioutil.ReadFile("./config.json")

    if err != nil {
      fmt.Println(err.Error())
      return config, err
    }

    err = json.Unmarshal(file, &config)

    if err != nil {
      fmt.Println(err.Error())
      return config, err
    }

    return config, nil

}


type configServers struct {
   Servers []string `json:"Servers"`
}

func ServersLoad() ([]string, error) {
  var config *configServers

  b_servers, err := ioutil.ReadFile("./servers.json")

  if err != nil {
    result := []string{"Cannot read the servers.json file!"}
    return result, errors.New(result[0])
  }

  err = json.Unmarshal(b_servers, &config)

  if err != nil {
    result := []string{"Cannot parse the servers.json file!"}
    return result, errors.New(result[0])
  }

  return config.Servers, nil

}
