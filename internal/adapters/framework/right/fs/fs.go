package fs

import (
	"log"
	"os"

	"github.com/eventials/go-tus"
)

func (da Adapter) SaveBinary(binary string) (string, error) {
	f, err := os.Open(binary)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	upload, err := tus.NewUploadFromFile(f)
	if err != nil {
		return "", err
	}
	// create the uploader.
	uploader, err := da.fs.CreateUpload(upload)
	if err != nil {
		return "", err
	}
	// start the uploading process.
	err = uploader.Upload()
	if err != nil {
		return "", err
	}
	
	return uploader.Url(),nil
}

// Adapter implements the DbPort interface
type Adapter struct {
	fs *tus.Client
}

// NewAdapter creates a new Adapter
func NewAdapter(server string) (*Adapter, error) {
	// connect
	client, err := tus.NewClient(server, nil)
	if err != nil {
		log.Fatalf("db ping failure: %v", err)
	}

	return &Adapter{fs: client}, nil
}
