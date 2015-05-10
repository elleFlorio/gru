package autonomic

import (
	"fmt"

	"github.com/elleFlorio/gru/action"
)

func execute(act action.Action) {
	//Execute stuff
	fmt.Println("I'm executing...")
	act.Execute()
}
