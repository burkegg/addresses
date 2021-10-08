package addresses

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDBConfig(t *testing.T) {
	dbConfig, _ := BuildDBConfig("test-host", 8000, "test-user", "test-db", "test-pw", "1.0.0")
	assert.True(t, dbConfig.Host == "test-host")
	assert.True(t, dbConfig.Port == 8000)
}

func TestReadCSV(t *testing.T) {
	dbConfig, _ := BuildDBConfig("test-host", 8000, "test-user", "test-db", "test-pw", "1.0.0")
	res := dbConfig.ImportData("./assets/addressdata.csv")
	assert.True(t, len(res) == 350)
}
