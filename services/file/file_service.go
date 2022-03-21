package file_service

import (
	"context"
	"log"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

type FileService struct {
	l      *log.Logger
	client *cloudinary.Cloudinary
}

func connectToCloudinary(l *log.Logger) *cloudinary.Cloudinary {
	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUD_NAME"), os.Getenv("API_KEY"), os.Getenv("API_SECRET"))
	if err != nil {
		l.Fatalln("Error connecting to cloudinary service")
		return nil
	}

	return cld

}

func NewFileService(l *log.Logger) IFileService {
	cloudinaryClient := connectToCloudinary(l)

	return FileService{
		l:      l,
		client: cloudinaryClient,
	}
}

func (iS FileService) UploadFile(ctx context.Context, file multipart.File) (string, error) {
	res, err := iS.client.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		iS.l.Println(err)
		return "", err
	}

	return res.URL, nil
}
