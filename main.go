package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ilhamarrouf/gin-queue/queue"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	server   *http.Server
	osSignal chan os.Signal
)

func main()  {
	osSignal = make(chan os.Signal, 10000)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("In app queue system")

	emailQueue := queue.NewEmailQueue()

	app := gin.Default()

	// route handler
	app.POST("/auth/register", func(context *gin.Context) {
		emailQueue.Enqueue("Send email to the user.")

		context.JSON(http.StatusCreated, gin.H{
			"data": gin.H{
				"username": "ilhamarrouf",
				"email": "ilham.arrouf@gmail.com",
			},
			"message": "Successfully register",
			"status": http.StatusCreated,
		})
	})

	server = &http.Server{
		Addr: ":8000",
		Handler: app,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Unexpected server error because of: %v\n", err)
		}
	}()

	for i:= 0; i < 10; i++ {
		go emailQueue.Work()
	}

	<-osSignal

	fmt.Println("Terminating server")
	server.Shutdown(context.Background())

	fmt.Println("Terminating email queue")
	for emailQueue.Size() > 0  {
		time.Sleep(time.Millisecond * 500)
	}

	fmt.Println("Complete terminating application")
}
