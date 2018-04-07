package gdice

import (
  "fmt"
  "strings"
  "strconv"
  "regexp"
  "time"
  "errors"
  
  "github.com/dgryski/go-pcgr"
  "github.com/Knetic/govaluate"
)

//show test info?
const flg bool = false
const regExpTokenDices string = "\\d+d\\d+|d\\d+"


type Roller struct {
  lastRoll string
  s_dices []string

  //this for verbose output
  verbose bool
  all_dices_results [][]string
}

//create new Roller
func New() *Roller {
  roller := new(Roller)
  roller.verbose = false
  return roller
}

//get lastRoll string
func (r *Roller) LastRoll() string  {
  return r.lastRoll
}

//set verbose to verbose output
func (r *Roller) Verbose(verbose bool) {
  r.verbose = verbose
}

//return dice pairs and numbers
func (r *Roller) parseDiceQueryToArray(diceString string) ([][]int64, []string) {

  s_parts := regexp.MustCompile(regExpTokenDices).Split(diceString, -1)
  s_dices := regexp.MustCompile(regExpTokenDices).FindAllString(diceString, -1);
  r.s_dices = s_dices

  //parse s_dices string to int pairs
  var dices [][]int64
  for i := 0; i < len(s_dices); i++ {
    node := regexp.MustCompile("d").Split(s_dices[i], -1)

    first, err := strconv.ParseInt(strings.TrimSpace(node[0]), 10, 64)
    if err != nil {
      first = 1
    }

    second, err := strconv.ParseInt(strings.TrimSpace(node[1]), 10, 64)

    if err != nil {
      second = 1
    }

    num_node := []int64 {first, second}

    dices = append(dices, num_node)
  }

  if flg{
    //test show
    fmt.Println("func parseDiceQueryToArray():")
    fmt.Println("\tdiceString:", diceString)
    testShow("\ts_parts", s_parts)
    testShow("\ts_dices", s_dices)
    fmt.Println("\tdices", dices)
  }

  return dices, s_parts
}


//concat dices results with other parts and calc as one expression
func (r *Roller) calc(s_dices_results []string, s_parts []string) (interface{}, error)  {

  var reser interface {}

  //concat s_dices_results and s_parts to parse as math expression
  expr := r.diceConc(s_dices_results, s_parts)
  //parse as math expression
  expression, err := govaluate.NewEvaluableExpression(expr);

  if err != nil {
    if flg{
      fmt.Println("func calc():")
      fmt.Println("\texpression parse error!")
      fmt.Println("\terr", err)
    }
    return reser, err
  }

  //calc as math expression
  reser, err = expression.Evaluate(nil);

  if err != nil {
    if flg{
      fmt.Println("func calc():")
      fmt.Println("\texpression parse error!")
      fmt.Println("\terr", err)
    }
    return reser, err
  }

  return reser, nil
}

//concat s_dices_results ant other parts (numbers) to one full expressions for math parser
func (r *Roller) diceConc(s_dices_results []string, s_parts []string) string {
  var res string
  for i, v := range(s_dices_results) {
    res += s_parts[i] + v
  }
  res += s_parts[len(s_parts) - 1]

  if flg {
    fmt.Println("func diceConc():")
    fmt.Println("\ts_dices_results", s_dices_results)
    fmt.Println("\ts_parts", s_parts)
    fmt.Println("\tres", res)
  }

  return res
}

//save string lastRoll
func (r *Roller) saveLastRoll(s_dices_results []string, s_parts []string)  {
  var lastRoll string
  //build lastRoll
  //[3d6: 14] + [2d3: 3] + 15 = 32.0
  for i, v := range(s_dices_results){
    lastRoll += s_parts[i]
    lastRoll += "["
    lastRoll += r.s_dices[i]
    lastRoll += ": "
    if r.verbose == true {
      lastRoll += r.putDiceValues(i)
      lastRoll += "="
      lastRoll += v
    }else {
      lastRoll += v
    }
    lastRoll += "]"
  }
  lastRoll += s_parts[len(s_parts) - 1]
  r.lastRoll = lastRoll
}

//one of the implementations put dice value to lastRoll
func (r *Roller) putDiceValues(index int) string{
  var DiceValues string
  //skip last element for right printing ","
  for _, delm := range r.all_dices_results[index][:len(r.all_dices_results[index]) - 1]{
    DiceValues += string(delm) + "+"
  }
  //put the last element after last ","
  DiceValues += r.all_dices_results[index][len(r.all_dices_results[index]) - 1]

  return DiceValues
}

//make result of roll as a string
func (r *Roller) RollString(diceString string) (string, bool)  {

  var s_res string

  reser, err := r.Roll(diceString)
  if err != nil{
    return "", false
  }

  //convert res interface {} to float64 or bool and to result string
  switch reser.(type)  {
  case float64:
    floata := reser.(float64)
    d := floata - float64(int64(floata))

    if d == 0 {
      s_res = strconv.FormatFloat(reser.(float64),  'f', 0, 64) + " = " + r.lastRoll
    } else {
      s_res = strconv.FormatFloat(reser.(float64),  'f', -1, 64) + " = " + r.lastRoll
    }

  case bool:
    s_res = strconv.FormatBool(reser.(bool))
  }


  if flg {
    fmt.Println("func RollString():")
    fmt.Println("\ts_res", s_res)
  }

  return s_res, true
}

//diceString = "3d6"
func (r *Roller) Roll(diceString string) (interface{}, error)  {
  //parse full string query, return [][]uint32 and []string
  dices, s_parts := r.parseDiceQueryToArray(diceString)

  //results of rolled dices
  var s_dices_results []string

  //nulling the all_dices_results
  r.all_dices_results = nil

  //calc dices rolls and fill s_dices_results
  for i := 0; i < len(dices); i++ {
    dice := dices[i]
    res := r.rollManyDice(dice[0], dice[1])
    s_res := strconv.FormatInt(res, 10)

    s_dices_results = append(s_dices_results, s_res) //save results of rolls
  }

  //calc full math expression
  reser, err := r.calc(s_dices_results, s_parts)
  r.saveLastRoll(s_dices_results, s_parts)

  if err != nil {
    return reser, errors.New("Cannot calc expression")
  }

  return reser, nil
}

//roll many dice
func (r *Roller) rollManyDice(number int64, dice int64) int64 {
  var res int64
  var last_elm_ind int
  if r.verbose == true {
    //add new roll
    r.all_dices_results = append(r.all_dices_results, []string{})
    last_elm_ind = len(r.all_dices_results) - 1
  }

  for i := int64(0); i < number; i++ {
    oneRoll := r.rollOneDice(dice)
    res += oneRoll

    if r.verbose == true {
      //save rolls
      r.all_dices_results[last_elm_ind] = append(r.all_dices_results[last_elm_ind], strconv.FormatInt(oneRoll, 10))
    }
  }
  return res
}

//emulate one dice roll
func (r *Roller) rollOneDice(dice int64) int64 {
  rand := pcgr.New(time.Now().UTC().UnixNano(),time.Now().UTC().UnixNano())
  if dice == 0 {
    return 0
  }
  res := rand.Int63()%6 + 1

  return res
}


//probability density function (PDF)
func Pdf(){

}

//table of chance get X and larger with dice
func ChanceLarger(){

}

//table of chance get X and larger with dice
func ChanceLower(){

}

//show string array
func testShow(prefix string, s []string)  {
  fmt.Print(prefix + ": ")
  for _, v := range(s){
    fmt.Print("[" + v + "]")
  }
  fmt.Println()
}
