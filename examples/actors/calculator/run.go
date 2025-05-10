// With Full Vibes (Github Copilot using Claude 3.7 Sonnet)
// prompt:
// create an example of a simple calculator where :
// - an actor receives an expression
// - the expression is read from the user like the echo example
// - supported operations are parentesis addition and multiplication
// - delegates spawned children for operation priority
// - actors uses the request/reply pattern to bubble up the results
// - and a list of examples of expressions to test with the related expected result
package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/morphy76/lang-actor/pkg/builders"
	"github.com/morphy76/lang-actor/pkg/framework"
)

// CalcState represents the state of a calculator actor
type calcState struct {
	// The expression being processed at each level
	expression string
	// The result of evaluation
	result float64
	// For parentheses expression tracking
	isParenthesis bool
	// Original expression containing the parentheses
	containingExpr string
	// Position of the opening parenthesis in the containing expression
	openPos int
	// Position of the closing parenthesis in the containing expression
	closePos int
}

// MessageType defines types of calculator messages
type messageType int

const (
	// Request to evaluate an expression
	messageTypeEvaluate messageType = iota
	// Result returned after evaluation
	messageTypeResult
	// Internal message to exit the program
	messageTypeExit
)

// CalcMessage represents messages sent between calculator actors
type calcMessage struct {
	sender         url.URL
	mexType        messageType
	expression     string    // Expression to evaluate
	result         float64   // Result of evaluation
	exitChan       chan bool // Channel to signal program exit
	isParenthesis  bool      // Whether this is a parenthetical expression
	containingExpr string    // Original expression containing the parentheses
	openPos        int       // Position of opening parenthesis
	closePos       int       // Position of closing parenthesis
}

// Sender implements the framework.Message interface
func (m calcMessage) Sender() url.URL {
	return m.sender
}

// Mutation implements the framework.Message interface
func (m calcMessage) Mutation() bool {
	return m.mexType == messageTypeResult
}

// calculatorFn is the main processing function for calculator actors
func calculatorFn(
	msg framework.Message,
	self framework.Actor[calcState],
) (calcState, error) {
	// Type assert to our specific message type
	useMsg := msg.(calcMessage)

	switch useMsg.mexType {
	case messageTypeEvaluate:
		return evaluateExpression(useMsg, self)
	case messageTypeResult:
		return handleResult(useMsg, self)
	case messageTypeExit:
		fmt.Println("Shutting down calculator...")
		useMsg.exitChan <- true
		return self.State(), nil
	default:
		return self.State(), fmt.Errorf("unknown message type")
	}
}

// evaluateExpression handles the evaluation of expressions
func evaluateExpression(msg calcMessage, self framework.Actor[calcState]) (calcState, error) {
	expression := strings.TrimSpace(msg.expression)

	// Update state with the current expression and parenthesis info
	state := calcState{
		expression:     expression,
		isParenthesis:  msg.isParenthesis,
		containingExpr: msg.containingExpr,
		openPos:        msg.openPos,
		closePos:       msg.closePos,
	}

	if expression == "" {
		if parent, found := self.GetParent(); found {
			self.Send(calcMessage{
				sender:         self.Address(),
				mexType:        messageTypeResult,
				result:         0,
				isParenthesis:  state.isParenthesis,
				containingExpr: state.containingExpr,
				openPos:        state.openPos,
				closePos:       state.closePos,
			}, parent)
		} else {
			fmt.Printf("Result: 0\n")
		}
		return state, nil
	}

	// Step 1: Check if this is a simple number
	if value, err := strconv.ParseFloat(expression, 64); err == nil {
		state.result = value

		if parent, found := self.GetParent(); found {
			// Send result back to parent
			self.Send(calcMessage{
				sender:         self.Address(),
				mexType:        messageTypeResult,
				result:         value,
				isParenthesis:  state.isParenthesis,
				containingExpr: state.containingExpr,
				openPos:        state.openPos,
				closePos:       state.closePos,
			}, parent)
		} else {
			// If we're the root actor, print the result
			fmt.Printf("Result: %g\n", value)
		}

		return state, nil
	}

	// Step 2: Check for parentheses first (highest priority)
	if openIdx := strings.LastIndex(expression, "("); openIdx >= 0 {
		// Find matching closing parenthesis
		openCount := 1
		closeIdx := -1

		for i := openIdx + 1; i < len(expression); i++ {
			if expression[i] == '(' {
				openCount++
			} else if expression[i] == ')' {
				openCount--
				if openCount == 0 {
					closeIdx = i
					break
				}
			}
		}

		if closeIdx == -1 {
			return state, fmt.Errorf("unbalanced parentheses in expression: %s", expression)
		}

		// Extract the expression inside the innermost parentheses
		innerExpr := expression[openIdx+1 : closeIdx]

		// Create a child actor to evaluate the inner expression
		childActor, err := builders.SpawnChild(self, calculatorFn, calcState{})
		if err != nil {
			return state, fmt.Errorf("failed to create child actor: %w", err)
		}

		// Send the inner expression to the child actor
		err = childActor.Deliver(calcMessage{
			sender:         self.Address(),
			mexType:        messageTypeEvaluate,
			expression:     innerExpr,
			isParenthesis:  true,
			containingExpr: expression,
			openPos:        openIdx,
			closePos:       closeIdx,
		})
		if err != nil {
			return state, fmt.Errorf("failed to deliver message: %w", err)
		}

		return state, nil
	}

	// Step 3: No more parentheses, handle operations based on precedence
	// First handle multiplication (higher precedence)
	if strings.Contains(expression, "*") {
		return handleMultiplication(expression, self, state)
	}

	// Then handle addition (lower precedence)
	if strings.Contains(expression, "+") {
		return handleAddition(expression, self, state)
	}

	// If we can't handle the expression, return an error
	return state, fmt.Errorf("invalid expression: %s", expression)
}

// handleResult processes results coming back from child actors
func handleResult(msg calcMessage, self framework.Actor[calcState]) (calcState, error) {
	// Update state with the result
	state := self.State()
	state.result = msg.result

	// If this is a result from a parenthesis evaluation
	if msg.isParenthesis && msg.containingExpr != "" {
		// Replace the parenthesized expression with its value
		resultStr := fmt.Sprintf("%g", msg.result)
		newExpr := msg.containingExpr[:msg.openPos] + resultStr + msg.containingExpr[msg.closePos+1:]

		// Re-evaluate the expression with the substituted value
		return state, self.Deliver(calcMessage{
			sender:         self.Address(),
			mexType:        messageTypeEvaluate,
			expression:     newExpr,
			isParenthesis:  self.State().isParenthesis,
			containingExpr: self.State().containingExpr,
			openPos:        self.State().openPos,
			closePos:       self.State().closePos,
		})
	}

	// Otherwise just pass the result up to our parent
	if parent, found := self.GetParent(); found {
		self.Send(calcMessage{
			sender:         self.Address(),
			mexType:        messageTypeResult,
			result:         state.result,
			isParenthesis:  self.State().isParenthesis,
			containingExpr: self.State().containingExpr,
			openPos:        self.State().openPos,
			closePos:       self.State().closePos,
		}, parent)
	} else {
		// If we're the root actor, print the result
		fmt.Printf("Result: %g\n", state.result)
	}

	return state, nil
}

// handleMultiplication handles expressions with multiplication
func handleMultiplication(expression string, self framework.Actor[calcState], state calcState) (calcState, error) {
	// Split the expression by multiplication
	parts := strings.Split(expression, "*")

	// Try to evaluate each part
	values := make([]float64, 0, len(parts))
	allSimple := true

	for _, part := range parts {
		part = strings.TrimSpace(part)

		// Try to parse as simple number
		if value, err := strconv.ParseFloat(part, 64); err == nil {
			values = append(values, value)
		} else if strings.Contains(part, "+") {
			// Contains addition, spawn child actor to evaluate
			allSimple = false

			childActor, err := builders.SpawnChild(self, calculatorFn, calcState{})
			if err != nil {
				return state, fmt.Errorf("failed to create child actor: %w", err)
			}

			err = childActor.Deliver(calcMessage{
				sender:     self.Address(),
				mexType:    messageTypeEvaluate,
				expression: part,
			})
			if err != nil {
				return state, fmt.Errorf("failed to deliver message: %w", err)
			}

			// Child actors will send results back, we'll handle in handleResult
			return state, nil
		} else {
			return state, fmt.Errorf("invalid term in multiplication: %s", part)
		}
	}

	// If all parts were simple numbers, calculate the result
	if allSimple {
		result := 1.0
		for _, v := range values {
			result *= v
		}

		state.result = result

		// Send result to parent
		if parent, found := self.GetParent(); found {
			self.Send(calcMessage{
				sender:         self.Address(),
				mexType:        messageTypeResult,
				result:         result,
				isParenthesis:  state.isParenthesis,
				containingExpr: state.containingExpr,
				openPos:        state.openPos,
				closePos:       state.closePos,
			}, parent)
		} else {
			fmt.Printf("Result: %g\n", result)
		}
	}

	return state, nil
}

// handleAddition handles expressions with addition
func handleAddition(expression string, self framework.Actor[calcState], state calcState) (calcState, error) {
	// Split the expression by addition
	parts := strings.Split(expression, "+")

	// Try to evaluate each part
	values := make([]float64, 0, len(parts))
	allSimple := true

	for _, part := range parts {
		part = strings.TrimSpace(part)

		// Try to parse as simple number
		if value, err := strconv.ParseFloat(part, 64); err == nil {
			values = append(values, value)
		} else if strings.Contains(part, "*") {
			// Contains multiplication, spawn child actor to evaluate
			allSimple = false

			childActor, err := builders.SpawnChild(self, calculatorFn, calcState{})
			if err != nil {
				return state, fmt.Errorf("failed to create child actor: %w", err)
			}

			err = childActor.Deliver(calcMessage{
				sender:     self.Address(),
				mexType:    messageTypeEvaluate,
				expression: part,
			})
			if err != nil {
				return state, fmt.Errorf("failed to deliver message: %w", err)
			}

			// Child actors will send results back, we'll handle in handleResult
			return state, nil
		} else {
			return state, fmt.Errorf("invalid term in addition: %s", part)
		}
	}

	// If all parts were simple numbers, calculate the result
	if allSimple {
		result := 0.0
		for _, v := range values {
			result += v
		}

		state.result = result

		// Send result to parent
		if parent, found := self.GetParent(); found {
			self.Send(calcMessage{
				sender:         self.Address(),
				mexType:        messageTypeResult,
				result:         result,
				isParenthesis:  state.isParenthesis,
				containingExpr: state.containingExpr,
				openPos:        state.openPos,
				closePos:       state.closePos,
			}, parent)
		} else {
			fmt.Printf("Result: %g\n", result)
		}
	}

	return state, nil
}

func main() {
	// Create an exit channel for graceful shutdown
	exitChan := make(chan bool)

	// Create the main calculator actor
	calcURL, _ := url.Parse("actor://calculator")
	calcActor, err := builders.NewTransientActor(*calcURL, calculatorFn, calcState{})
	if err != nil {
		fmt.Println("Error creating calculator actor:", err)
		return
	}

	// Ensure we stop the actor when done
	defer func() {
		done, _ := calcActor.Stop()
		<-done
	}()

	// Create a reader for user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Calculator with Actors")
	fmt.Println("Supported operations: parentheses (), addition +, multiplication *")
	fmt.Println("Enter 'exit' to quit.")
	fmt.Println("-----------------------------------------")

	// Main input loop
	for {
		fmt.Print("Enter expression: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			// Send exit message to calculator
			calcActor.Deliver(calcMessage{
				sender:   *calcURL,
				mexType:  messageTypeExit,
				exitChan: exitChan,
			})
			break
		}

		// Send the input expression to the calculator
		calcActor.Deliver(calcMessage{
			sender:     *calcURL,
			mexType:    messageTypeEvaluate,
			expression: input,
			exitChan:   exitChan,
		})
	}

	// Wait for exit signal
	<-exitChan
	fmt.Println("Calculator stopped.")
}
