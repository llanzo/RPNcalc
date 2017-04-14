package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var validNumber = regexp.MustCompile(`^[0-9]+\.?[0-9]*`)
var operands = [2]float64{}
var runningTotal, displayValue float64
var operandCount int
var hasRunningTotal = false
var display string

func parseInput(input string) string {
	var isValidOperation = func() (bool, string) {
		switch {
		case operandCount > len(operands):
			return false, "Error, exactly " + strconv.Itoa(len(operands)) + " operands per operation"
		case (operandCount == 0):
			return false, "Error, insufficient operands for operation"
		case (operandCount == 2 || (operandCount == 1 && hasRunningTotal)):
			return true, ""
		default:
			return false, "Unrecognized state"
		}
	}
	var isAcceptingOperands = func(stringValue string) (bool, string) {
		matched := validNumber.MatchString(stringValue)
		switch {
		case operandCount >= len(operands):
			return false, "Error, exactly " + strconv.Itoa(len(operands)) + " operands per operation"
		case !matched:
			return false, "Error, token " + stringValue + " is not a valid operand"
		default:
			return true, ""
		}
	}
	var ensureMultipleOperands = func() {
		if operandCount == 1 && hasRunningTotal {
			operands[1] = operands[0]
			operands[0] = runningTotal
		}
	}
	var updateCaculatorRegisters = func(calculatedValue float64) string {
		runningTotal = calculatedValue
		hasRunningTotal = true
		displayValue = runningTotal
		operandCount = 0
		operands = [2]float64{}
		return strconv.FormatFloat(displayValue, 'f', 1, 32)
	}
	var insertAnotherOperand = func(tokenValue string) string {
		enabled, unacceptableExplaination := isAcceptingOperands(tokenValue)
		if !enabled {
			return unacceptableExplaination
		}
		numberToken, err := strconv.ParseFloat(tokenValue, 64)
		if err != nil {
			return "Error parsing " + tokenValue
		}
		operands[operandCount] = numberToken
		operandCount++
		return "result: " + strconv.FormatFloat(numberToken, 'f', 1, 32)
	}
	var add = func() string {
		ensureMultipleOperands()
		return updateCaculatorRegisters(operands[0] + operands[1])
	}
	var subtract = func() string {
		ensureMultipleOperands()
		return updateCaculatorRegisters(operands[0] - operands[1])
	}
	var multiply = func() string {
		ensureMultipleOperands()
		return updateCaculatorRegisters(operands[0] * operands[1])
	}
	var divide = func() string {
		ensureMultipleOperands()
		return updateCaculatorRegisters(operands[0] / operands[1])
	}

	var tokens = strings.Split(input, " ")
	for _, token := range tokens {
		var operable, inoperableExplaination = isValidOperation()
		if operable {
			switch {
			case token == `+`:
				display = ("result: " + add())
			case token == "-":
				display = ("result: " + subtract())
			case token == "*":
				display = ("result: " + multiply())
			case token == "/":
				display = ("result: " + divide())
			default:
				display = insertAnotherOperand(token)
			}

		} else if token == `+` || token == `-` || token == `*` || token == `/` {
			display = inoperableExplaination
		} else {
			display = insertAnotherOperand(token)
		}
	}
	return display
}

var holderChan = make(chan string)

func localParser(ch chan string) {
	for result := range ch {
		holderChan <- parseInput(result)
	}
}

func main() {
	ch := make(chan string)
	go localParser(ch)
	scanner := bufio.NewScanner(os.Stdin)
	var text string

	for {
		fmt.Print("input expression (q to quit) > ")
		scanner.Scan()
		text = scanner.Text()
		if text == "q" {
			break
		}
		if len(text) > 0 && text[0] == '\u0004' {
			break
		}
		ch <- text
		var showValue = <-holderChan
		println(showValue)
	}
}