package main

import (
	"fmt"
	"net/http"

	"github.com/streadway/amqp"
)

func sendMessage(name string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Printf("An error occured: %s\n", err.Error())
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("An error occured: %s\n", err.Error())
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"TestNovalQueue",
		false,
		false,
		false,
		false,
		nil,
	)

	err = ch.Publish(
		"",
		"TestNovalQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(name),
		},
	)

	if err != nil {
		fmt.Printf("Error when publish message: %s\n", err.Error())
	}
	fmt.Println("Success sent message queue")
}

func main() {
	fmt.Println("Message Queue Basic Publiser Rest")

	ch := make(chan string)
	go func() {
		for {
			data := <-ch
			sendMessage(data)
		}
		//
	}()

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			name := r.FormValue("name")
			ch <- name
			w.Write([]byte("Success send message"))
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.ListenAndServe(":8080", nil)
}
