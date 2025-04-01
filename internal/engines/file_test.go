package engines_test

import (
	"context"
	"go-metrics/internal/engines"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockFileStoragePathGetter struct {
	Path string
}

func (m *MockFileStoragePathGetter) GetFileStoragePath() string {
	return m.Path
}

func TestFileEngine_Open(t *testing.T) {
	tempDir := filepath.Join(t.TempDir(), "testdir")
	filePath := filepath.Join(tempDir, "metrics.json")
	engine := engines.NewFileEngine()
	fsp := &MockFileStoragePathGetter{Path: filePath}
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Fatalf("expected file %s to not exist before opening", filePath)
	}

	if err := engine.Open(context.Background(), fsp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("expected file %s to exist after opening", filePath)
	}
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Fatalf("expected directory %s to exist", tempDir)
	}
}

func TestFileWriterEngine_Write(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	fsp := &MockFileStoragePathGetter{Path: tmpFile.Name()}

	fileEngine := &engines.FileEngine{}
	require.NoError(t, fileEngine.Open(context.Background(), fsp))
	defer fileEngine.File.Close()

	writerEngine := &engines.FileWriterEngine[map[string]interface{}]{FileEngine: fileEngine}

	data := []map[string]interface{}{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}

	err = writerEngine.Write(context.Background(), data)
	require.NoError(t, err)
}

func TestFileGeneratorEngine_Generate(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	fsp := &MockFileStoragePathGetter{Path: tmpFile.Name()}

	fileEngine := &engines.FileEngine{}
	require.NoError(t, fileEngine.Open(context.Background(), fsp))
	defer fileEngine.File.Close()

	writerEngine := &engines.FileWriterEngine[map[string]interface{}]{FileEngine: fileEngine}
	generatorEngine := &engines.FileGeneratorEngine[map[string]interface{}]{FileEngine: fileEngine}

	data := []map[string]interface{}{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}

	err = writerEngine.Write(context.Background(), data)
	require.NoError(t, err)

	// Проверяем содержимое файла после записи
	content, err := os.ReadFile(tmpFile.Name()) // заменили ioutil.ReadFile на os.ReadFile
	require.NoError(t, err)
	t.Logf("File content before reading: %s", string(content))

	// Проверяем позицию файла перед чтением
	pos, err := fileEngine.File.Seek(0, 1)
	require.NoError(t, err)
	t.Logf("File position before reading: %d", pos)

	// Синхронизируем и сбрасываем указатель файла перед чтением
	require.NoError(t, fileEngine.File.Sync())
	_, err = fileEngine.File.Seek(0, 0)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var results []map[string]interface{}
	for record := range generatorEngine.Generate(ctx) {
		t.Logf("Read record: %+v", record)
		results = append(results, record)
	}

	// Преобразуем 'id' в int для корректного сравнения
	for i := range results {
		if id, ok := results[i]["id"].(float64); ok {
			results[i]["id"] = int(id)
		}
	}

	// Теперь сравниваем с ожидаемыми результатами
	assert.Equal(t, data, results)
}
