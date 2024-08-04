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
	"os"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"github.com/kbinani/screenshot"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/image/draw"
)

const (
	REMOTE          = "localhost:4567"
	HOST            = "http://" + REMOTE
	WS_URL          = "ws://" + REMOTE + "/x/%s/ws"
	UPLOAD_URL      = HOST + "/upload"
	SCREENSHOT_PATH = "/tmp/screenshot.jpg"
	KEY_LENGTH      = 10
	SALT_LENGTH     = 16
	IV_LENGTH       = 12
)

var (
	conn           *websocket.Conn
	lastScreenshot image.Image
	encryptionKey  string
	sessionID      string
)

func main() {
	var err error
	sessionID, err = createSession()
	if err != nil {
		printDebug(fmt.Sprintf("Failed to create session: %v", err))
		os.Exit(1)
	}

	encryptionKey = generateRandomKey(KEY_LENGTH)
	fmt.Printf("Access the website at: %s/x/%s#%s\n", HOST, sessionID, encryptionKey)

	for {
		err := connectWebSocket()
		if err == nil {
			break
		}
		printDebug(fmt.Sprintf("WebSocket connection failed: %v. Retrying in 5 seconds...", err))
		time.Sleep(5 * time.Second)
	}
	defer conn.Close()

	go sendCursorPosition()

	for {
		img := captureScreen()

		if !imagesEqual(img, lastScreenshot) {
			printDebug("Images are not equal. Uploading new screenshot.")
			encryptedData, err := resizeAndEncryptScreen(img)
			if err != nil {
				printDebug(fmt.Sprintf("Failed to resize and encrypt screenshot: %v", err))
				continue
			}
			for {
				err := uploadEncryptedScreen(encryptedData)
				if err == nil {
					lastScreenshot = img
					break
				}
				printDebug(fmt.Sprintf("Failed to upload screenshot: %v. Retrying...", err))
				time.Sleep(1 * time.Second)
			}
		} else {
			printDebug("Images are equal. Skipping upload.")
		}

		// Sleep for 950ms before the next iteration
		time.Sleep(950 * time.Millisecond)
	}
}

func createSession() (string, error) {
	resp, err := http.Post(HOST+"/v1/api/sessions", "application/json", nil)
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
	conn, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf(WS_URL, sessionID), nil)
	return err
}

func sendCursorPosition() {
	var lastX, lastY int
	for {
		x, y := robotgo.GetMousePos()
		screenWidth, screenHeight := robotgo.GetScreenSize()
		scaledX := int(float64(x) * 1280.0 / float64(screenWidth))
		scaledY := int(float64(y) * 720.0 / float64(screenHeight))

		if scaledX != lastX || scaledY != lastY {
			data := map[string]int{
				"x": scaledX,
				"y": scaledY,
			}
			err := conn.WriteJSON(data)
			if err != nil {
				printDebug(fmt.Sprintf("WebSocket write failed: %v", err))
				for {
					err := connectWebSocket()
					if err == nil {
						break
					}
					printDebug(fmt.Sprintf("WebSocket reconnection failed: %v. Retrying in 5 seconds...", err))
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
			printDebug("No active displays found")
			time.Sleep(1 * time.Second)
			continue
		}

		// Capture the first display as an example
		bounds := screenshot.GetDisplayBounds(0)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			printDebug("Failed to capture screen: cannot capture display: locked or switched off, retrying...")
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
		printDebug("Unexpected image format: not RGBA")
		return false
	}

	return bytes.Equal(rgba1.Pix, rgba2.Pix)
}

// resizeAndEncryptScreen resizes the captured screenshot to a fixed width, encodes it as JPEG, and encrypts it.
func resizeAndEncryptScreen(img image.Image) ([]byte, error) {
	newWidth := 1280
	newHeight := img.Bounds().Dy() * newWidth / img.Bounds().Dx()
	scaledImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.BiLinear.Scale(scaledImg, scaledImg.Bounds(), img, img.Bounds(), draw.Over, nil)

	var buf bytes.Buffer
	err := jpeg.Encode(&buf, scaledImg, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}

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

	req, err := http.NewRequest("POST", UPLOAD_URL, body)
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

	printDebug("Uploaded encrypted screenshot")
	return nil
}

func printDebug(message string) {
	currentTime := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s\n", currentTime, message)
}
