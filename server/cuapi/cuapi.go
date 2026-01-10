package cuapi

import (
	"effective-invention/server/cuapi/cuhandlers"
	"fmt"
)

func ConnectClickUp() {
	err := cuhandlers.AuthClickUp()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
