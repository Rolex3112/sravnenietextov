package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	cl "github.com/fatih/color"
)

func main() {
	fmt.Println("сравнение текстов: белый — совпадающие буквы,")
	fmt.Println("красный жирный — разные буквы (на обеих строках),")
	fmt.Println("жёлтый — «лишние» буквы, если один текст длиннее")
	reader := bufio.NewReader(os.Stdin)

	cl.New(cl.FgGreen, cl.Bold).Print("Введите первый текст: ")
	t1, _ := reader.ReadString('\n')
	t1 = strings.TrimRight(t1, "\r\n")

	cl.New(cl.FgGreen, cl.Bold).Print("Введите второй текст: ")
	t2, _ := reader.ReadString('\n')
	t2 = strings.TrimRight(t2, "\r\n")

	r1 := []rune(t1)
	r2 := []rune(t2)

	maxLen := len(r1)
	if len(r2) > maxLen {
		maxLen = len(r2)
	}

	sameStyle := cl.New(cl.FgWhite)
	diffStyle := cl.New(cl.FgRed, cl.Bold, cl.BgHiBlack)
	extraStyle := cl.New(cl.FgYellow, cl.Bold)

	line1 := &strings.Builder{}
	line2 := &strings.Builder{}

	for i := 0; i < maxLen; i++ {
		switch {
		case i >= len(r1):
			line2.WriteString(extraStyle.Sprint(string(r2[i])))
		case i >= len(r2):
			line1.WriteString(extraStyle.Sprint(string(r1[i])))
		case r1[i] == r2[i]:
			line1.WriteString(sameStyle.Sprint(string(r1[i])))
			line2.WriteString(sameStyle.Sprint(string(r2[i])))
		default:
			line1.WriteString(diffStyle.Sprint(string(r1[i])))
			line2.WriteString(diffStyle.Sprint(string(r2[i])))
		}
	}

	fmt.Println()
	fmt.Println(line1.String())
	fmt.Println(line2.String())

	cl.New(cl.FgYellow, cl.Bold).Println("\nи вот шутка:")
	setup, punchline, err := fetchJoke()
	if err != nil {
		cl.New(cl.FgRed).Println("не удалось получить шутку:", err)
		return
	}
	cl.New(cl.FgCyan).Println(setup)
	cl.New(cl.FgMagenta, cl.Bold).Println(punchline)
}

type joke struct {
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

func fetchJoke() (string, string, error) {
	resp, err := http.Get("https://official-joke-api.appspot.com/random_joke")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var j joke
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return "", "", err
	}
	return j.Setup, j.Punchline, nil
}
