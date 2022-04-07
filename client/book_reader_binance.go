package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/omhen/swissblock-trade-executor/v2/errors"
	"github.com/omhen/swissblock-trade-executor/v2/model"
	log "github.com/sirupsen/logrus"
)

type bookReaderBinance struct {
	quit      chan bool
	streamURL string
}

func (br *bookReaderBinance) StartReading(ctx context.Context, symbol string) (chan *model.BookItem, chan error) {
	receiver := make(chan *model.BookItem)
	errChannel := make(chan error)

	go br.readLoop(ctx, symbol, receiver, errChannel)

	return receiver, errChannel
}

func (br *bookReaderBinance) Stop(ctx context.Context) {
	br.quit <- true
}

func (br *bookReaderBinance) readLoop(ctx context.Context, symbol string, receiver chan *model.BookItem, errChannel chan error) {
	defer close(receiver)
	defer close(errChannel)

	conn, err := br.connectToStream(ctx, symbol)
	if err != nil {
		errChannel <- err
		return
	}

	defer conn.Close()

	for {
		select {
		case <-br.quit:
			return
		default:
			bookItem := new(model.BookItem)
			err = conn.ReadJSON(&bookItem)
			if err != nil {
				log.Error("While reading from stream: ", err)
				errChannel <- errors.Annotate(err, "While reading from stream")
			} else {
				receiver <- bookItem
			}
		}
	}
}

func (br *bookReaderBinance) connectToStream(ctx context.Context, symbol string) (*websocket.Conn, error) {
	url := fmt.Sprintf("%s/%s@bookTicker", br.streamURL, strings.ToLower(symbol))
	log.WithFields(log.Fields{
		"stream": url,
	}).Info("Connecting to stream")
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"stream": url,
		}).Error("Could not connect to stream")

		return nil, errors.Annotate(err, "Could not connect to stream")
	}

	return conn, nil
}

func NewBookReaderBinance(streamURL string) BookReader {
	return &bookReaderBinance{
		streamURL: streamURL,
		quit:      make(chan bool),
	}
}
