package main

import (
	"fmt"
	"os"

	. "github.com/logrusorgru/aurora"
)

func ErrorOutput(output string, extra ...string) {
	fmt.Println(Black(" Error ").BgRed(), Faint(output))
	if len(extra) > 0 {
		for _, str := range extra {
			fmt.Println("       ", Faint(str))
		}
	}
	os.Exit(1)
}

func StatusOutput(output string, extra ...string) {
	fmt.Println(Bold("==>"), output)
	if len(extra) > 0 {
		for _, str := range extra {
			fmt.Println("   ", Faint(str))
		}
	}
}
