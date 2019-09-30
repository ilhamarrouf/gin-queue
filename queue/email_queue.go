package queue

import (
	"fmt"
	"time"
)

type emailQueue struct {
	emailChannel chan string
	workingChannel chan bool
}

func NewEmailQueue() *emailQueue {
	emailChannel := make(chan string, 10000)
	workingChannel := make(chan bool, 10000)

	return &emailQueue{
		emailChannel: emailChannel,
		workingChannel: workingChannel,
	}
}

func (e *emailQueue) Work() {
	for {
		select {
		case eChan := <-e.emailChannel:
			e.workingChannel <- true

			time.Sleep(time.Second * 2)
			fmt.Println(eChan)

			<-e.workingChannel
		}
	}
}

func (e *emailQueue) Size() int {
	return len(e.emailChannel) + len(e.workingChannel)
}

func (e *emailQueue) Enqueue(emailString string)  {
	e.emailChannel <- emailString
}
