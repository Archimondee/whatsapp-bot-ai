package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/disintegration/imaging"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

const (
	letterBytes string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateFileName(extension string) string {
	timestamp := time.Now().Unix()
	named := fmt.Sprintf("%d%s", timestamp, RandomStrings(10))
	HashMD5 := Md5Hash(named)
	return fmt.Sprintf("%s%s", HashMD5, extension)
}

func RandomStrings(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	randomString := string(b)
	return randomString
}

func Md5Hash(message string) string {
	md5Hash := md5.Sum([]byte(message))
	md5HashString := hex.EncodeToString(md5Hash[:])
	return md5HashString
}

func OpusEncode(audio string) string {
	ogg_name := strings.ReplaceAll(audio, ".mp3", ".ogg")
	cmd := exec.Command("ffmpeg", "-i", audio, "-c:a", "libopus", "-b:a", "64k", ogg_name)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return ogg_name
}

func ResizeImage(imagePath string, thumbnailPath string, h int, w int, q int) string {
	img, err := imaging.Open(imagePath)
	if err != nil {
		fmt.Println("failed to open image:")
	}
	img = imaging.Fit(img, w, h, imaging.Lanczos)
	err = imaging.Save(img, thumbnailPath, imaging.JPEGQuality(q))
	if err != nil {
		fmt.Println("failed to save thumbnail")
	}
	return thumbnailPath

}
