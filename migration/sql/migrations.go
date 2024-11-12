package migrations

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestVersionOrder(t *testing.T) {
	commitRefName := os.Getenv("CI_COMMIT_REF_NAME")
	if commitRefName != "master" &&
		!strings.HasPrefix(strings.ToLower(commitRefName), "release") {
		log.Info("Skip migration files check on development branch")
		return
	}

	sourceUrl := "file://."
	sourceDrv, err := source.Open(sourceUrl)
	assert.NoError(t, err, "failed open migration source", err)
	currVersion, err := sourceDrv.First()
	assert.NoError(t, err, "failed get first migration", err)
	for {
		nextVersion, err := sourceDrv.Next(currVersion)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) { //Next will return os.ErrNotExist for last version
				assert.Fail(t, "failed get next migration", err)
			}
			break
		}
		assert.Equal(t, currVersion+1, nextVersion,
			"expected:%d, actual: %d", currVersion+1, nextVersion)
		currVersion = nextVersion
	}
}
