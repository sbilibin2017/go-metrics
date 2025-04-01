package engines

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// FileStoragePathGetter interface is used to get the file storage path.
type FileStoragePathGetter interface {
	GetFileStoragePath() string
}

type FileEngine struct {
	*os.File
	once sync.Once
}

func (e *FileEngine) Open(ctx context.Context, fsp FileStoragePathGetter) error {
	var err error
	e.once.Do(func() {
		filePath := fsp.GetFileStoragePath()
		dir := filepath.Dir(filePath)
		if err = os.MkdirAll(dir, 0755); err != nil {
			return
		}
		file, openErr := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
		if openErr != nil {
			err = openErr
			return
		}
		e.File = file
	})
	return err
}

type FileWriterEngine[T any] struct {
	*FileEngine
	mu   sync.Mutex
	once sync.Once
}

func (e *FileWriterEngine[T]) Write(ctx context.Context, data []T) error {
	e.once.Do(func() {
		e.mu.Lock()
		defer e.mu.Unlock()

		if err := e.File.Truncate(0); err != nil {
			return
		}
		if _, err := e.File.Seek(0, 0); err != nil {
			return
		}
		encoder := json.NewEncoder(e.File)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(data); err != nil {
			return
		}
		if err := e.File.Sync(); err != nil {
			return
		}
	})
	return nil
}

type FileGeneratorEngine[T any] struct {
	*FileEngine
	once sync.Once
}

func (e *FileGeneratorEngine[T]) Generate(ctx context.Context) <-chan T {
	ch := make(chan T)
	e.once.Do(func() {
		go func() {
			defer close(ch)
			_, err := e.File.Seek(0, 0)
			if err != nil {
				return
			}
			decoder := json.NewDecoder(e.File)
			var records []T
			if err := decoder.Decode(&records); err != nil {
				return
			}
			for _, record := range records {
				select {
				case ch <- record:
				case <-ctx.Done():
					return
				}
			}
		}()
	})

	return ch
}

func NewFileEngine() *FileEngine {
	return &FileEngine{}
}

func NewFileWriterEngine[T any](fileEngine *FileEngine) *FileWriterEngine[T] {
	return &FileWriterEngine[T]{FileEngine: fileEngine}
}

func NewFileGeneratorEngine[T any](fileEngine *FileEngine) *FileGeneratorEngine[T] {
	return &FileGeneratorEngine[T]{FileEngine: fileEngine}
}
