package clif

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lazybark/go-helpers/cli/clf"
)

// ActionList is an object to simplify user-CLI interact. It promts available actions and waits for
// input.
type ActionList struct {
	Message           string
	Actions           []Action
	InputLineMessage  string
	WrongInputMessage string
	CanGoBack         bool
	BackKey           string
	CanExit           bool
	ExitKey           string
	CaseSensitive     bool
}

// Action represents specific option for a user to choose
type Action struct {
	Key  string
	Text string
}

var (
	//EmptyAction should be returned in case no actions need to be returned in specific case.
	EmptyAction = Action{}
)

// Promt tells user available options
func (al *ActionList) Promt() {
	fmt.Print(al.Message)

	if len(al.Actions) == 0 {
		if al.CanExit || al.CanGoBack {
			fmt.Print(" or \n")
		}
	} else {
		fmt.Print(": \n")
	}

	for _, act := range al.Actions {
		fmt.Printf("- %s %s\n", clf.Green(act.Key), act.Text)
	}

	if al.CanGoBack {
		fmt.Printf("- %s to go back\n", clf.Yellow(al.BackKey))
	}
	if al.CanExit {
		fmt.Printf("- %s to exit\n", clf.Red(al.ExitKey))
	}
	fmt.Print(al.InputLineMessage)
}

// AwaitCommand awaits for user input and returns option picked by a user.
// Returns EmptyAction in case user picked something wrong.
func (al *ActionList) AwaitCommand() Action {
	bk := al.BackKey
	ek := al.ExitKey
	if al.WrongInputMessage == "" {
		al.InputLineMessage = "unknown option"
	}
	if !al.CaseSensitive {
		bk = strings.ToLower(al.BackKey)
		ek = strings.ToLower(al.ExitKey)
	}

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()

		if text == "" {
			continue
		}

		if !al.CaseSensitive {
			text = strings.ToLower(text)
		}

		if text == bk && al.CanGoBack {
			return Action{Key: al.BackKey}
		}
		if text == ek && al.CanExit {
			return Action{Key: al.ExitKey}
		}
		for _, opt := range al.Actions {
			if text == strings.ToLower(opt.Key) {
				return opt
			}
		}

		fmt.Printf("%s: %s\n\n", text, clf.Red(al.WrongInputMessage))
		return Action{}
	}
}

// AwaitInput returns action picked by user OR string that user has entered. String is parsed via
// parseFunc. If parseFunc returns error, user will see error message and will be asked for input again
func (al *ActionList) AwaitInput(parseFunc func(string) error) (Action, string) {
	bk := al.BackKey
	ek := al.ExitKey
	if !al.CaseSensitive {
		bk = strings.ToLower(al.BackKey)
		ek = strings.ToLower(al.ExitKey)
	}

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()

		if text == "" {
			continue
		}

		if !al.CaseSensitive {
			text = strings.ToLower(text)
		}

		if text == bk && al.CanGoBack {
			return Action{Key: al.BackKey}, ""
		}
		if text == ek && al.CanExit {
			return Action{Key: al.ExitKey}, ""
		}
		if err := parseFunc(text); err != nil {
			fmt.Printf("%s\n\n", clf.Red(err))
			return Action{}, ""
		}
		return Action{}, text
	}
}
