package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var dir string
var usingFile string
var links []string

func main() {

	//Создание лога
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)

	//Прием данных пользователя
	fmt.Print("Выберите файл с ссылками: ")
	var setFile string
	fmt.Scan(&setFile)
	usingFile = setFile

	fmt.Print("Укажите директорию для сохранения: ")
	var selectDir string
	fmt.Scan(&selectDir)
	os.MkdirAll(selectDir, 0777)
	dir = selectDir

	log.Println("Директория для сохранения: " + selectDir)
	log.Println("Файл с ссылками: " + setFile)

	createFile()	// Создание файла
}

func createFile() {

	setLink()	// Чтение ссылок из файла

	for i := 0; i < len(links); i++ {

		file, err := os.Create(setName())	// Создание файла с уникальным именем

		if err != nil {
			log.Fatal("Невозможно создать файл: ", err)
		}

		defer file.Close()

		file.WriteString(getData(i))	// Запись данных в файл
	}

	fmt.Println("Создание файлов завершено успешно.")
	log.Println("Создание файлов завершено успешно.")
}

func getData(i int) string {

	resp, err := http.Get(links[i])

	if err != nil {
		fmt.Println("Ошибка доступа к ссылке, подробности в логе.")
		log.Print("Ссылка не доступна!", err)

		return ""
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

	log.Println("Файл найден. Считываем ссылки...")

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
