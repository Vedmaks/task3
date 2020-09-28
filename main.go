package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var dir string
var usingFile string
var links []string

func main() {

	t0 := time.Now()

	//Создание лога
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)

	flag.StringVar(&usingFile, "f", usingFile, "file")
	flag.StringVar(&dir, "d", dir, "dir")

	flag.Parse()

	err = os.MkdirAll(dir, 0777)

	if err != nil {
		log.Fatal("Ошибка директории: ", err)
	}

	log.Println("Директория для сохранения: " + dir)
	log.Println("Файл с ссылками: " + usingFile)

	setLink() // Чтение ссылок из файла

	for i := 0; i < len(links); i++ {

		go createFile(i) // Создание файла
	}

	var input string
	fmt.Scanln(&input)
	fmt.Println("Создание файлов завершено.")
	log.Println("Создание файлов завершено.")
	t1 := time.Now()
	fmt.Printf("Elapsed time: %v", t1.Sub(t0))
}

func createFile(i int) {

	var name string
	var fullName string

	name = strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					links[i], ".", "_"), "http://", ""), "https://", ""), "/", "_")

	fullName = dir + "/" + name + ".html"

	file, err := os.Create(fullName)

	if err != nil {
		log.Println("Невозможно создать файл: ", err)
	}

	defer file.Close()

	data := getData(links[i])

	file.WriteString(data) // Запись данных в файл

	if data == "" {
		log.Println("Запрашиваемая страница недопустимого типа: " + links[i] + " Строка: " + fmt.Sprint(i+1))
		fmt.Println("Запрашиваемая страница недопустимого типа:  " + fmt.Sprint(i+1) + " из " + fmt.Sprint(len(links)) + ".")
	} else if data == "err" {
		fmt.Println("Ошибка сервера: " + fmt.Sprint(i+1) + " из " + fmt.Sprint(len(links)) + ".")
	} else {
		fmt.Println("Создан файл: " + fmt.Sprint(i+1) + " из " + fmt.Sprint(len(links)) + ".")
		//log.Println("Создан файл: " + file.Name() + "	" + fmt.Sprint(i+1) + " из " + fmt.Sprint(len(links)) + ".")
	}
}

func getData(link string) string {

	resp, err := http.Get(link)

	if err != nil {
		log.Print("Ошибка сервера: ", err)

		return "err"
	}

	defer resp.Body.Close()

	var data string

	for {

		bs := make([]byte, 1024)
		n, err := resp.Body.Read(bs)

		if n == 0 || err != nil {
			break
		}

		data += (string(bs[:n]))
	}

	return data
}

func setLink() {

	file, err := os.Open(usingFile)

	if err != nil {
		fmt.Println("Ошибка открытия файла, подробности в логе.")
		log.Fatal("Ошибка открытия файла: ", err)
	}

	defer file.Close()

	log.Println("Файл найден.")

	var linkFile string
	data := make([]byte, 64)

	for {
		n, err := file.Read(data)
		if err == io.EOF {
			log.Println("Считывание ссылок завершено.")
			break
		}
		linkFile += (string(data[:n]))
	}

	links = strings.Fields(linkFile)
}
