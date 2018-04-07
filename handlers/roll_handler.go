package handlers

import (
  "fmt"

  "../utilize/gdice"
)

func Roll(content string, verbose bool) string {
  var answer string

  roller := gdice.New()

  roller.Verbose(verbose)

  if res, ok := roller.RollString( content ); ok {
    answer = res
  }else {
    answer = "Wrong expression!"
  }

  fmt.Println("roll answer:", answer)

  return answer
}
