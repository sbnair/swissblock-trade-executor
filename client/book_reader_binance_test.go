package client_test

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/omhen/swissblock-trade-executor/v2/client"
	"github.com/omhen/swissblock-trade-executor/v2/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type wsMockServer struct {
	addr     string
	path     string
	messages []string
	server   *http.Server
}

func (ms *wsMockServer) Setup() {
	var upgrader = websocket.Upgrader{}
	http.HandleFunc(ms.path, func(w http.ResponseWriter, r *http.Request) {
		c, _ := upgrader.Upgrade(w, r, nil)
		defer c.Close()
		for _, message := range ms.messages {
			c.WriteMessage(1, []byte(message))
		}
	})
}

func (ms *wsMockServer) SetMessages(messages []string) {
	ms.messages = messages
}

func (ms *wsMockServer) Start() {
	ms.server = &http.Server{Addr: ms.addr, Handler: nil}
	go ms.server.ListenAndServe()
}

func (ms *wsMockServer) Stop() {
	ms.server.Shutdown(context.Background())
}

var _ = Describe("BookReaderBinance", func() {
	var addr = "localhost:8080"
	var path = "/ws/bnbusdt@bookTicker"
	mockServer := wsMockServer{addr: addr, path: path}
	mockServer.Setup()

	BeforeEach(func() {
		mockServer.Start()
		time.Sleep(500 * time.Millisecond) // Wait for the server to be up
	})

	AfterEach(func() {
		mockServer.Stop()
	})

	Describe("Reading orders from the order book", func() {
		Context("from the mock WS Server", func() {
			It("should return three orders", func() {
				messages := []string{
					`{"u":7000459416,"s":"BNBUSDT","b":"449.60000000","B":"9.22000000","a":"449.70000000","A":"221.55800000"}`,
					`{"u":7000459423,"s":"BNBUSDT","b":"449.60000000","B":"5.87500000","a":"449.70000000","A":"221.55800000"}`,
					`{"u":7000459441,"s":"BNBUSDT","b":"449.60000000","B":"5.87500000","a":"449.70000000","A":"223.55800000"}`,
				}
				mockServer.SetMessages(messages)
				reader := client.NewBookReaderBinance("ws://localhost:8080/ws")
				orderChan, errChan := reader.StartReading(context.Background(), "BNBUSDT")
				received := make([]*model.BookItem, 0)
			OuterLoop:
				for {
					select {
					case message := <-orderChan:
						received = append(received, message)
					case err := <-errChan:
						Expect(err.Error()).To(Equal("While reading from stream: websocket: close 1006 (abnormal closure): unexpected EOF"))
						break OuterLoop
					}
				}
				Expect(len(received)).To(Equal(3))
				for i, uid := range []uint{7000459416, 7000459423, 7000459441} {
					Expect(received[i].UpdateID).To(Equal(uid))
				}
			})

			It("should return an error when a message is malformed", func() {
				messages := []string{`malformed json`}
				mockServer.SetMessages(messages)
				reader := client.NewBookReaderBinance("ws://localhost:8080/ws")
				orderChan, errChan := reader.StartReading(context.Background(), "BNBUSDT")
			OuterLoop:
				for {
					select {
					case message := <-orderChan:
						Expect(message).To(BeNil())
					case err := <-errChan:
						Expect(err).To(HaveOccurred())
						break OuterLoop
					}
				}
			})
		})

		Context("from a wrong URL", func() {
			It("should return an error through the error channel", func() {
				reader := client.NewBookReaderBinance("ws://verywrongwebsocket:8080/ws")
				orderChan, errChan := reader.StartReading(context.Background(), "BNBUSDT")
			OuterLoop:
				for {
					select {
					case message := <-orderChan:
						Expect(message).To(BeNil())
					case err := <-errChan:
						Expect(err).To(HaveOccurred())
						break OuterLoop
					}
				}
			})
		})
	})
})
