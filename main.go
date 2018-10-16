package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/skip2/go-qrcode"
	"image/color"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

const landingUrl = "https://qr.stendahl.me/mRgTTRaywA"

var randomRunes = []rune("123456789")

func randomString(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = randomRunes[rand.Intn(len(randomRunes))]
	}

	return string(b)
}

func handleFinalRequest(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, string(getQrCode("https://www.youtube.com/watch?v=dQw4w9WgXcQ")))
}

func handleLandingRequest(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, string(getRandomQrCode(0)))
}

func handleRandomRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seed, err := strconv.Atoi(vars["random"])

	if err != nil {
		panic(err)
	}

	handleRequest(w, string(getRandomQrCode(int64(seed))))
}

func handleRequest(w http.ResponseWriter, qrCode string) {
	hasher := sha1.New()
	hasher.Write([]byte(qrCode))
	hash := hasher.Sum(nil)

	w.Header().Set("ETag", fmt.Sprintf("%x\n", hash))
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=31536000")

	fmt.Fprint(w, qrCode)
}

func getRandomQrCode(seed int64) []byte {
	rand.Seed(seed)

	return getQrCode("https://qr.stendahl.me/" + randomString(10))
}

func getQrCode(url string) []byte {
	qr, err := qrcode.New(url, qrcode.Highest)

	if err != nil {
		panic(err)
	}

	qr.ForegroundColor = color.RGBA{
		R: uint8(rand.Intn(100)),
		G: uint8(rand.Intn(100)),
		B: uint8(rand.Intn(100)),
		A: 255,
	}

	png, err := qr.PNG(512)

	if err != nil {
		panic(err)
	}

	return png
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handleLandingRequest)
	r.HandleFunc("/1234567890", handleFinalRequest)
	r.HandleFunc("/{random:[0-9a-zA-Z]{10}}", handleRandomRequest)
	log.Fatal(http.ListenAndServe("127.0.0.1:8283", r))
}