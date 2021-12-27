package gdrivesvc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Service struct {
	*drive.Service
}

type FileMeta struct {
	Name     string
	MimeType string
	FolderID string
}

func New(ctx context.Context, client *http.Client) *Service {
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	return &Service{srv}
}

func (s *Service) UploadFile(f *os.File, meta FileMeta, doneCh chan struct{}, errCh chan error) {
	if _, err := s.Files.Create(&drive.File{
		Name:     meta.Name,
		MimeType: meta.MimeType,
		Parents:  []string{meta.FolderID},
	}).Media(f).Do(); err != nil {
		errCh <- fmt.Errorf("could not create file '%s': %s", meta.Name, err)
	}
	doneCh <- struct{}{}
}
