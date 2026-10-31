package migrations_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/loiuscypher/photoview/api/database/migrations"
	"github.com/loiuscypher/photoview/api/graphql/models"
	"github.com/loiuscypher/photoview/api/test_utils"
)

func TestExifMigration(t *testing.T) {
	db := test_utils.DatabaseTest(t)
	defer db.Exec("DELETE FROM media_exif") // Clean up after test

	// Create test data
	exifEntries := []models.MediaEXIF{
		{GPSLatitude: new(90.1), GPSLongitude: new(180.0)},   // Invalid GPSLatitude
		{GPSLatitude: new(-90.1), GPSLongitude: new(-180.0)}, // Invalid GPSLatitude
		{GPSLatitude: new(90.0), GPSLongitude: new(180.1)},   // Invalid GPSLongitude
		{GPSLatitude: new(-90.0), GPSLongitude: new(-180.1)}, // Invalid GPSLongitude
		{GPSLatitude: new(90.0), GPSLongitude: new(180.0)},   // Valid GPS data
		{GPSLatitude: new(-90.0), GPSLongitude: new(-180.0)}, // Valid GPS data
		{GPSLatitude: new(35.7), GPSLongitude: new(139.7)},   // Valid GPS data, longitude above 90
		{GPSLatitude: new(90.1), GPSLongitude: new(180.1)},   // Invalid GPSLatitude and GPSLongitude
		{GPSLatitude: new(-90.1), GPSLongitude: new(-180.1)}, // Invalid GPSLatitude and GPSLongitude
	}

	// Insert test data
	for _, entry := range exifEntries {
		assert.NoError(t, db.Create(&entry).Error)
	}

	// Run migration
	assert.NoError(t, migrations.MigrateForExifGPSCorrection(db))

	// Validate the results
	var results []models.MediaEXIF
	assert.NoError(t, db.Find(&results).Error)

	validEntries := 0
	for _, entry := range results {
		if entry.GPSLatitude != nil {
			assert.LessOrEqual(t, math.Abs(*entry.GPSLatitude), 90.0, "GPSLatitude should be within [-90, 90]: %+v", entry)
		}
		if entry.GPSLongitude != nil {
			assert.LessOrEqual(t, math.Abs(*entry.GPSLongitude), 180.0, "GPSLongitude should be within [-180, 180]: %+v", entry)
		}
		if entry.GPSLatitude != nil && entry.GPSLongitude != nil {
			validEntries++
		}
	}
	assert.Equal(t, 3, validEntries, "valid GPS data should be kept")
}

//go:fix inline
func floatPtr(f float64) *float64 {
	return new(f)
}
