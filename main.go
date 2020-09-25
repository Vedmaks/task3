package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var dir string
var usingFile string

func main() {

	fmt.Print("Set file to use: ")
	var setFile string
	fmt.Scan(&setFile)
	usingFile = setFile

	fmt.Print("Select dir: ")
	var selectDir string
	fmt.Scan(&selectDir)
	os.Mkdir(selectDir, 0777)
	dir = selectDir

	createFile(20)

}

func createFile(num int) {

	for i := 0; i < num; i++ {

		file, err := os.Create(setName())
		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(2)
		}
		defer file.Close()

		file.WriteString(getData(i))
		fmt.Println("Файл создан: " + file.Name())
		fmt.Println("Всего создано:")
		fmt.Println(i + 1)
	}

}

func getData(num int) string {

	resp, err := http.Get(setLink(num))
	if err != nil {
		fmt.Println(err)

	}
	defer resp.Body.Close()

	var data string

	for true {

		bs := make([]byte, 1014)
		n, err := resp.Body.Read(bs)

		if n == 0 || err != nil {
			break
		}

		data += (string(bs[:n]))
	}

	return data
}

func setLink(num int) string {

	file, err := os.Open(usingFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var links string
	data := make([]byte, 64)

	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		links += (string(data[:n]))
	}

	linkNum := strings.Fields(links)

	if len(linkNum) == num+1 {

		os.Exit(1)
	}

	return linkNum[num]

}

func setName() string {
	var fullName string

	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()

	fullName = dir + "/" + str + ".html"

	return fullName
}
