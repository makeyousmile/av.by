package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
type Country struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Emoji string `json:"emoji"`
	Code  string `json:"code"`
}

type PhoneNumber struct {
	ID      int     `json:"id"`
	Country Country `json:"country"`
	Number  string  `json:"number"`
}

func main() {
	//getPhones(1)
	url := "https://api.av.by/offers/113588239/phones"

	// Создание HTTP-запроса
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	// Установка заголовка User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ошибка: статус ответа", resp.StatusCode)
		return
	}

	// Чтение ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	// Вывод ответа для проверки
	fmt.Println("Ответ от API:", string(body))

	// Парсинг JSON
	var phoneNumbers []PhoneNumber
	err = json.Unmarshal(body, &phoneNumbers)
	if err != nil {
		fmt.Println("Ошибка при парсинге JSON:", err)
		return
	}

	// Вывод поля number
	for _, phoneNumber := range phoneNumbers {
		fmt.Println("Number:", phoneNumber.Number)
	}

}
