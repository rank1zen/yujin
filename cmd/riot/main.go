package main

import (
	"log"
	"os"

	"github.com/KnutZuidema/golio"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rank1zen/yujin/internal/riot"
)

func main() {
	client := golio.NewClient(os.Getenv("RIOT_API_KEY"))
	q := riot.NewSummonerQ(client.Riot.LoL.Summoner)
	_ = q

	conn, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		log.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	queue, err := ch.QueueDeclare(
		"summoner_renewal",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received: %s", d.Body)
			log.Printf("Received: %s", d.ContentType)
		}
	}()

	log.Printf("Listening [*]")
	<-forever

}
