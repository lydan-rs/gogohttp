package main

import (
	"http-protocol/internal/request"
	"http-protocol/internal/response"
	"http-protocol/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

const UNSUPPORTED_METHOD = "<html>" + "\n" +
	"<head>" + "\n" +
	"<title>400 Bad Request</title>" + "\n" +
	"</head>" + "\n" +
	"<body>" + "\n" +
	"<h1>Bad Request</h1>" + "\n" +
	"<p>Unsupported Method.</p>" + "\n" +
	"</body>" + "\n" +
	"</html>" + "\n"

const BAD_REQUEST = "<html>" + "\n" +
	"<head>" + "\n" +
	"<title>400 Bad Request</title>" + "\n" +
	"</head>" + "\n" +
	"<body>" + "\n" +
	"<h1>Bad Request</h1>" + "\n" +
	"<p>Your request honestly kinda sucked.</p>" + "\n" +
	"</body>" + "\n" +
	"</html>" + "\n"

const SERVER_ERROR = "<html>" + "\n" +
	"<head>" + "\n" +
	"<title>500 Internal Server Error</title>" + "\n" +
	"</head>" + "\n" +
	"<body>" + "\n" +
	"<h1>Internal Server Error</h1>" + "\n" +
	"<p>Okay, you know what? This one is on me.</p>" + "\n" +
	"</body>" + "\n" +
	"</html>" + "\n"

const SUCCESS = "<html>" + "\n" +
	"<head>" + "\n" +
	"<title>200 OK</title>" + "\n" +
	"</head>" + "\n" +
	"<body>" + "\n" +
	"<h1>Success!</h1>" + "\n" +
	"<p>Your request was an absolute banger.</p>" + "\n" +
	"</body>" + "\n" +
	"</html>" + "\n"

func MainHandler(w response.Writer, req *request.Request) {
	if req.RequestLine.Method != "GET" {
		w.WriteStatusLine(response.SC_BAD_REQUEST)
		w.WriteHTML([]byte(UNSUPPORTED_METHOD))
		return
	}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.SC_BAD_REQUEST)
		w.WriteHTML([]byte(BAD_REQUEST))

	case "/myproblem":
		w.WriteStatusLine(response.SC_INTERNAL_SERVER_ERROR)
		w.WriteHTML([]byte(SERVER_ERROR))

	default:
		w.WriteStatusLine(response.SC_OK)
		w.WriteHTML([]byte(SUCCESS))

	}

}

func main() {
	server, err := server.Serve(port, MainHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
