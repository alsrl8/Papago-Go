package main

import (
	"PapagoGo/api"
	"PapagoGo/powershell"
	"fmt"
)

func main() {
	fmt.Printf("%sPapago-Go Started!%s\n\n\n", powershell.ColorCyan, powershell.ColorReset)
	api.GetUserInputAndTranslate()
}
