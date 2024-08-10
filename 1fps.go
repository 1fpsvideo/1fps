package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"math/big"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/1fpsvideo/1fps/appconfig"
	"github.com/1fpsvideo/1fps/consoleui"
	"github.com/1fpsvideo/1fps/cursor"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"github.com/kbinani/screenshot"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/image/draw"
)

const (
	SCREENSHOT_PATH = "/tmp/screenshot.jpg"
	KEY_LENGTH      = 10
	SALT_LENGTH     = 16
	IV_LENGTH       = 12
)

var (
	conn              *websocket.Conn
	lastScreenshot    image.Image
	encryptionKey     string
	sessionID         string
	resizedDimensions struct {
		Width  int
		Height int
	}

	consoleUI *consoleui.ConsoleUI
	appConfig *appconfig.AppConfig
)

// log logs an event to the bottom panel
func log(message string) {
	consoleUI.WriteBottom(message)
}

func main() {
	appConfig = appconfig.New()

	// Try to get all the necessary info to start the console app.
	// Do not use initialized UI before we're getting what we need: in case of an error
	// we just need to print error to the console in a non-fancy way.

	var err error
	sessionID, err = createSession()
	if err != nil {
		panic(fmt.Sprintf("Failed to create session: %v", err))
	}

	// Fanciness stars here. Start the console app. All the events below are getting
	// logged with log method only (this way they end up in a log window).
	// App UI starts in its own goroutine.

	consoleUI = consoleui.Start()
	encryptionKey = generateRandomKey(KEY_LENGTH)
	consoleUI.SetUrl(fmt.Sprintf("%s/x/%s#%s", appConfig.Host, sessionID, encryptionKey))

	// Connecting to web socket before we start goroutine to send cursor coodinates.

	for {
		err := connectWebSocket()
		if err == nil {
			break
		}
		log(fmt.Sprintf("WebSocket connection failed: %v. Retrying in 5 seconds...", err))
		time.Sleep(5 * time.Second)
	}
	defer conn.Close()

	// Sending cursor coodinates.

	go sendCursorPosition()

	// Main loop: capture screen, compare, encrypt, send, pause, repeat.

	for {
		img := captureScreen()

		if !imagesEqual(img, lastScreenshot) {
			log("Images are not equal. Uploading new screenshot.")
			encryptedData, err := resizeAndEncryptScreen(img)
			if err != nil {
				log(fmt.Sprintf("Failed to resize and encrypt screenshot: %v", err))
				continue
			}
			for {
				err := uploadEncryptedScreen(encryptedData)
				if err == nil {
					lastScreenshot = img
					break
				}
				log(fmt.Sprintf("Failed to upload screenshot: %v. Retrying...", err))
				time.Sleep(1 * time.Second)
			}
		} else {
			log("Images are equal. Skipping upload.")
		}

		// Sleep for 950ms before the next iteration
		time.Sleep(950 * time.Millisecond)
	}
}

// The rest of the functions remain the same, but replace all printDebug calls with log

func createSession() (string, error) {
	resp, err := http.Post(appConfig.Host+"/v1/api/sessions", "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Status    string `json:"status"`
		SessionID string `json:"session_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Status != "ok" {
		return "", fmt.Errorf("failed to create session")
	}

	return result.SessionID, nil
}

func generateRandomKey(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func connectWebSocket() error {
	var err error
	conn, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf(appConfig.WsUrl, sessionID), nil)
	return err
}

func sendCursorPosition() {
	var lastX, lastY int
	for {
		scaledX, scaledY := cursor.GetCursorPosition(resizedDimensions)

		if scaledX != lastX || scaledY != lastY {
			data := map[string]int{
				"x":  scaledX,
				"y":  scaledY,
				"rw": resizedDimensions.Width,
				"rh": resizedDimensions.Height,
			}
			err := conn.WriteJSON(data)
			if err != nil {
				log(fmt.Sprintf("WebSocket write failed: %v", err))
				for {
					err := connectWebSocket()
					if err == nil {
						break
					}
					log(fmt.Sprintf("WebSocket reconnection failed: %v. Retrying in 5 seconds...", err))
					time.Sleep(5 * time.Second)
				}
			}
			lastX, lastY = scaledX, scaledY
		}

		time.Sleep(70 * time.Millisecond)
	}
}

// captureScreen captures the entire screen and returns the image.
func captureScreen() image.Image {
	for {
		n := screenshot.NumActiveDisplays()
		if n <= 0 {
			log("No active displays found")
			time.Sleep(1 * time.Second)
			continue
		}

		// Capture the first display as an example
		bounds := screenshot.GetDisplayBounds(0)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			log("Failed to capture screen: cannot capture display: locked or switched off, retrying...")
			time.Sleep(1 * time.Second)
			continue
		}

		return img
	}
}

// imagesEqual compares two images pixel by pixel and returns true if they are equal.
// imagesEqual compares two images and returns true if they are equal.
func imagesEqual(img1, img2 image.Image) bool {
	if img1 == nil || img2 == nil {
		return false
	}

	rgba1, ok1 := img1.(*image.RGBA)
	rgba2, ok2 := img2.(*image.RGBA)

	if !ok1 || !ok2 {
		log("Unexpected image format: not RGBA")
		return false
	}

	return bytes.Equal(rgba1.Pix, rgba2.Pix)
}

// resizeAndEncryptScreen resizes the captured screenshot to a fixed width, encodes it as JPEG, and encrypts it.
func resizeAndEncryptScreen(img image.Image) ([]byte, error) {
	// Get the screen width and calculate the target width
	screenWidth, _ := robotgo.GetScreenSize()
	targetWidth := screenWidth
	if targetWidth > 1280 {
		targetWidth = 1280
	}

	// Get the dimensions of the input image
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	var scaledImg image.Image

	// Check if resizing is necessary
	if imgWidth <= targetWidth {
		// No need to resize, use original dimensions
		resizedDimensions.Width = imgWidth
		resizedDimensions.Height = imgHeight
		scaledImg = img
	} else {
		// Resize the image
		resizedDimensions.Width = targetWidth
		resizedDimensions.Height = imgHeight * targetWidth / imgWidth
		scaledImg = image.NewRGBA(image.Rect(0, 0, resizedDimensions.Width, resizedDimensions.Height))
		draw.BiLinear.Scale(scaledImg.(draw.Image), scaledImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	}

	// Encode the image to JPEG
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, scaledImg, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}

	// Encrypt and return the image data
	return encryptData(buf.Bytes())
}

func encryptData(data []byte) ([]byte, error) {
	salt := make([]byte, SALT_LENGTH)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	iv := make([]byte, IV_LENGTH)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	key := pbkdf2.Key([]byte(encryptionKey), salt, 100000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, iv, data, nil)

	encryptedData := make([]byte, 0, len(salt)+len(iv)+len(ciphertext))
	encryptedData = append(encryptedData, salt...)
	encryptedData = append(encryptedData, iv...)
	encryptedData = append(encryptedData, ciphertext...)

	return encryptedData, nil
}

// uploadEncryptedScreen uploads the encrypted screenshot to the server.
func uploadEncryptedScreen(encryptedData []byte) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "screenshot.jpg")
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = part.Write(encryptedData)
	if err != nil {
		return fmt.Errorf("failed to write form file: %v", err)
	}
	err = writer.WriteField("session_id", sessionID)
	if err != nil {
		return fmt.Errorf("failed to write session_id field: %v", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	req, err := http.NewRequest("POST", appConfig.UploadUrl, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status: %s", resp.Status)
	}

	log("Uploaded encrypted screenshot")
	return nil
}
