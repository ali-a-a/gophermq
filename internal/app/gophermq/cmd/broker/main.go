package broker

import (
	"errors"
	"fmt"
	"github.com/ali-a-a/gophermq/internal/app/gophermq/broker"
	"github.com/ali-a-a/gophermq/internal/app/gophermq/broker/handler"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ali-a-a/gophermq/config"
	"github.com/ali-a-a/gophermq/pkg/router"
)

func main(cfg config.Config) {
	mq := broker.NewGopherMQ(broker.MaxPending(cfg.Broker.MaxPending))

	h := handler.NewHandler(mq)

	e := router.New()

	e.GET("/ready", func(c echo.Context) error { return c.NoContent(http.StatusNoContent) })

	api := e.Group("/api")

	api.POST("/publish", h.Publish)
	api.POST("/publish/async", h.PublishAsync)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", cfg.Broker.Port)); !errors.Is(err, http.ErrServerClosed) && err != nil {
			e.Logger.Fatal(err.Error())
		}
	}()

	for i := 0; i < len(cfg.Broker.Subjects); i++ {
		_, err := mq.Subscribe(cfg.Broker.Subjects[i], func(event broker.Event) error {
			logrus.Infof("subject: %s, data: %v", event.Subject(), string(event.Data()))

			return nil
		})
		if err != nil {
			return
		}
	}

	logrus.Info("broker is ready!")

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		logrus.Info("signal received: ", sig)
		done <- true
	}()

	<-done
}

// Register registers broker command for gophermq binary.
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "broker",
			Short: "Run broker component",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
