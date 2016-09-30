package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"flag"
	"os/exec"
	"os"
	"runtime"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}


type Consumerx struct {
	conn	*amqp.Connection
	channel	*amqp.Channel
	tag	string
	done	chan error
}
var (
	uri		= flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchange	= flag.String("exchange", "logs", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType	= flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue		= flag.String("queue", "logs", "Ephemeral AMQP queue name")
	bindingKey	= flag.String("key", "", "AMQP binding key")
	consumerTag	= flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
	wsPort		= flag.Int("ws-port", 23456, "WebSockets port to listen to")
)










var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cls") //Windows example it is untested, but I think its working
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok { //if we defined a clear func for that platform:
		value()  //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func main()   {
	amqpURI:="amqp://guest:guest@localhost:5672/"
	c := &Consumerx{
		conn:		nil,
		channel:	nil,
		tag:		"",
		done:		make(chan error),
	}
	c.conn, _ = amqp.Dial(amqpURI)

	c.channel, _ = c.conn.Channel()

	deliveries, err := c.channel.Consume(
		"pdna",	// name
		c.tag,	// consumerTag,
		true,	// noAck
		false,	// exclusive
		false,	// noLocal
		false,	// noWait
		nil,	// arguments
	)
	if err != nil {
		fmt.Errorf("Queue Consume: %s", err)
	}


	for d := range deliveries {
		log.Printf("got %s", d.Body)
		fmt.Println("I will clean the screen in 2 seconds!")
		time.Sleep(2 * time.Second)
		CallClear()
	}

}