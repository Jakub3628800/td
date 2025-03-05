package core

import (
	"fmt"
	"log"
	"os/exec"
)

func SendNotification(msg string, silent bool) {
	if silent {
		fmt.Println(msg)
		fmt.Println()
	}
	err := exec.Command("notify-send", msg).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func PauseMusic() {
	execPlayerctl("pause")
}

func PlayMusic() {
	execPlayerctl("play")
}

func execPlayerctl(subcmd string) {
	err := exec.Command("playerctl", subcmd).Run()
	if err != nil {
		log.Fatal(err)
	}
}
