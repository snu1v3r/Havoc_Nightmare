package main

import (
	"Havoc/cmd"
	"Havoc/pkg/colors"
	"fmt"
)
import "Havoc/pkg/logger"

func main() {
	startMenu()

	err := cmd.HavocCmd.Execute()
	if err != nil {
		logger.Error("Failed to execute havoc")
		return
	}
}

func startMenu() {
	fmt.Print(colors.Red("              _______           _______  _______ \n    │\\     /│(  ___  )│\\     /│(  ___  )(  ____ \\\n    │ )   ( ││ (   ) ││ )   ( ││ (   ) ││ (    \\/\n    │ (___) ││ (___) ││ │   │ ││ │   │ ││ │      \n    │  ___  ││  ___  │( (   ) )│ │   │ ││ │      \n    │ (   ) ││ (   ) │ \\ \\_/ / │ │   │ ││ │      \n    │ )   ( ││ )   ( │  \\   /  │ (___) ││ (____/\\\n    │/     \\││/     \\│   \\_/   (_______)(_______/"))
	fmt.Printf(" by %v\n\n", colors.BoldBlue("@C5pider"))
	fmt.Println("  	", colors.Red("pwn"), "and", colors.Blue("elevate"), "until it's done\n")
}
