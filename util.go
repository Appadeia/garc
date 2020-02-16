package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func PrettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	println(string(s))
}

func CheckErr(err error, extraprint ...string) {
	if err != nil {
		ErrorOutput(err.Error(), extraprint...)
	}
}

func RunCommand(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	CheckErr(err)
	return string(output)
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	randomInt := func(min, max int) int {
		return min + rand.Intn(max-min)
	}

	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func OpenFileInEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func PromptInEditor(template, prompt string) string {
	file, err := ioutil.TempFile(os.TempDir(), "*")
	CheckErr(err)

	if template != "" {
		file.WriteString(template)
	}

	if strings.Contains(prompt, "\n") {
		lines := strings.Split(prompt, "\n")
		for _, line := range lines {
			file.WriteString(fmt.Sprintf("# %s", line))
		}
	} else {
		file.WriteString(fmt.Sprintf("# %s", prompt))
	}

	filename := file.Name()
	defer os.Remove(filename)

	err = file.Close()
	CheckErr(err)

	err = OpenFileInEditor(filename)
	CheckErr(err)

	bytes, err := ioutil.ReadFile(filename)
	CheckErr(err)

	outputString := string(bytes)
	returnString := ""

	if !strings.Contains(outputString, "\n") {
		return outputString
	}

	for _, line := range strings.Split(outputString, "\n") {
		if !strings.HasPrefix(line, "#") {
			returnString = fmt.Sprintf("%s\n%s", returnString, line)
		}
	}

	return strings.TrimSpace(returnString)
}

func PromptInlineAnything(desc string) string {
	prompt := promptui.Prompt{
		Label: desc,
	}

	result, err := prompt.Run()
	CheckErr(err)

	return result
}

func PromptInlineChoice(desc string, choices ...string) string {
	prompt := promptui.Select{
		Label: desc,
		Items: choices,
	}

	_, result, err := prompt.Run()
	CheckErr(err)

	return result
}
