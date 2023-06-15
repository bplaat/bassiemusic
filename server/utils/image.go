package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"mime/multipart"
	"os"

	"github.com/bplaat/bassiemusic/core/uuid"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
)

func StoreUploadedImage(c *fiber.Ctx, storageDir string, id uuid.Uuid, file *multipart.FileHeader, saveLarge bool) error {
	idString := id.String()

	// Save uploaded image file
	if err := c.SaveFile(file, fmt.Sprintf("storage/%s/original/%s", storageDir, idString)); err != nil {
		log.Fatalln(err)
	}

	// Open uploaded image
	originalFile, err := os.Open(fmt.Sprintf("storage/%s/original/%s", storageDir, idString))
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
		if err := os.Remove(fmt.Sprintf("storage/%s/original/%s", storageDir, idString)); err != nil {
			log.Fatalln(err)
		}
		return fiber.ErrBadRequest
	}

	// Save small resize
	smallFile, err := os.Create(fmt.Sprintf("storage/%s/small/%s.jpg", storageDir, idString))
	if err != nil {
		log.Fatalln(err)
	}
	defer smallFile.Close()
	smallImage := imaging.Resize(originalImage, 250, 250, imaging.Lanczos)
	err = jpeg.Encode(smallFile, smallImage, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Save medium resize
	mediumFile, err := os.Create(fmt.Sprintf("storage/%s/medium/%s.jpg", storageDir, idString))
	if err != nil {
		log.Fatalln(err)
	}
	defer mediumFile.Close()
	mediumImage := imaging.Resize(originalImage, 500, 500, imaging.Lanczos)
	err = jpeg.Encode(mediumFile, mediumImage, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Save large resize
	if saveLarge {
		largeFile, err := os.Create(fmt.Sprintf("storage/%s/large/%s.jpg", storageDir, idString))
		if err != nil {
			log.Fatalln(err)
		}
		defer largeFile.Close()
		largeImage := imaging.Resize(originalImage, 1000, 1000, imaging.Lanczos)
		err = jpeg.Encode(largeFile, largeImage, nil)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return nil
}
