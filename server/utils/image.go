package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"mime/multipart"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
)

func StoreUploadedImage(c *fiber.Ctx, storageDir string, id string, file *multipart.FileHeader, saveLarge bool) error {
	// Save uploaded image file
	if err := c.SaveFile(file, fmt.Sprintf("storage/%s/original/%s", storageDir, id)); err != nil {
		log.Fatalln(err)
	}

	// Open uploaded image
	originalFile, err := os.Open(fmt.Sprintf("storage/%s/original/%s", storageDir, id))
	if err != nil {
		log.Fatalln(err)
	}
	defer originalFile.Close()
	if err != nil {
		log.Fatalln(err)
	}
	var originalImage image.Image
	originalImage, _, err = image.Decode(originalFile)
	if err != nil {
		if err := os.Remove(fmt.Sprintf("storage/%s/original/%s", storageDir, id)); err != nil {
			log.Fatalln(err)
		}
		return fiber.ErrBadRequest
	}

	// Save small resize
	smallFile, err := os.Create(fmt.Sprintf("storage/%s/small/%s.jpg", storageDir, id))
	if err != nil {
		log.Fatalln(err)
	}
	defer smallFile.Close()
	smallImage := resize.Resize(250, 250, originalImage, resize.Lanczos3)
	err = jpeg.Encode(smallFile, smallImage, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Save medium resize
	mediumFile, err := os.Create(fmt.Sprintf("storage/%s/medium/%s.jpg", storageDir, id))
	if err != nil {
		log.Fatalln(err)
	}
	defer mediumFile.Close()
	mediumImage := resize.Resize(500, 500, originalImage, resize.Lanczos3)
	err = jpeg.Encode(mediumFile, mediumImage, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Save large resize
	if saveLarge {
		largeFile, err := os.Create(fmt.Sprintf("storage/%s/large/%s.jpg", storageDir, id))
		if err != nil {
			log.Fatalln(err)
		}
		defer largeFile.Close()
		largeImage := resize.Resize(1000, 1000, originalImage, resize.Lanczos3)
		err = jpeg.Encode(largeFile, largeImage, nil)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return nil
}
