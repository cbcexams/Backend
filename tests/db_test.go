package tests

import (
	"cbc-backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	err := models.InitDB(testDBConn)
	assert.NoError(t, err)

	err = models.TestDatabaseConnection()
	assert.NoError(t, err)
}

func TestEnsureTables(t *testing.T) {
	assert.NoError(t, models.EnsureUsersTable())
	assert.NoError(t, models.EnsureJobsTable())
	assert.NoError(t, models.EnsureUploadsTable())
}
