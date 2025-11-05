package scanner

import (
	"context"
	"log"
	"os"
	"path"

	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/media_encoding"
	"github.com/photoview/photoview/api/scanner/scanner_cache"
	"github.com/photoview/photoview/api/scanner/scanner_task"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var ProcessSingleMediaFunc = ProcessSingleMedia

func ScanMedia(tx *gorm.DB, mediaPath string, albumId int, cache *scanner_cache.AlbumScannerCache) (*models.Media, bool, error) {
	mediaName := path.Base(mediaPath)

	// Check if media already exists
	{
		var media []*models.Media

		result := tx.Where("path_hash = ?", models.MD5Hash(mediaPath)).Find(&media)

		if result.Error != nil {
			return nil, false, errors.Wrap(result.Error, "scan media fetch from database")
		}

		if result.RowsAffected > 0 {
			// log.Printf("Media already scanned: %s\n", mediaPath)
			stat, err := os.Stat(mediaPath)
			if err != nil {
				// log.Printf("Media file should be there, but could not get stat: %s\n", mediaPath)
				return nil, false, err
			}
			// log.Printf("Media mtime diff: %s\n", stat.ModTime().Sub(media[0].DateShot).Abs().String())
			// Time threshold may be necessary to be increased for non unix systems
			if stat.ModTime().Sub(media[0].DateShot).Abs().Milliseconds() < 1 {
				// log.Printf("Media mtime did not change: %s\n", mediaPath)
				return media[0], false, nil
			}
			// log.Printf("Media mtime changed: %s\n", stat.ModTime().Sub(media[0].DateShot).Abs().String())
			// if err := tx.Where("ID = ?", media[0].ID).Delete(models.Media{}).Error; err != nil {
			//if err := tx.Unscoped().Model(models.MediaEXIF{}).Where("ID = ?",media[0].ExifID).Delete(media).Error; err != nil {
			if err := tx.Model(models.MediaEXIF{}).Where("ID = ?",media[0].ExifID).Delete(media).Error; err != nil {
				// log.Printf("Media deletion from DB failed: %s\n", mediaPath)
				return nil, false, err
			}
			log.Printf("Media.ID delete for update : %d\n", media[0].ID)
		}
	}

	log.Printf("Scanning media: %s\n", mediaPath)

	mediaType, err := cache.GetMediaType(mediaPath)
	if err != nil {
		return nil, false, errors.Wrap(err, "could determine if media was photo or video")
	}

	var mediaTypeText models.MediaType

	if mediaType.IsVideo() {
		mediaTypeText = models.MediaTypeVideo
	} else {
		mediaTypeText = models.MediaTypePhoto
	}

	stat, err := os.Stat(mediaPath)
	if err != nil {
		return nil, false, err
	}

	media := models.Media{
		Title:    mediaName,
		Path:     mediaPath,
		AlbumID:  albumId,
		Type:     mediaTypeText,
		DateShot: stat.ModTime(),
	}

	if err := tx.Create(&media).Error; err != nil {
		return nil, false, errors.Wrap(err, "could not insert media into database")
	}

	return &media, true, nil
}

// ProcessSingleMedia processes a single media, might be used to reprocess media with corrupted cache
// Function waits for processing to finish before returning.
func ProcessSingleMedia(ctx context.Context, db *gorm.DB, media *models.Media) error {
	albumCache := scanner_cache.MakeAlbumCache()

	var album models.Album
	if err := db.Model(media).Association("Album").Find(&album); err != nil {
		return err
	}

	mediaData := media_encoding.NewEncodeMediaData(media)

	taskContext := scanner_task.NewTaskContext(ctx, db, &album, albumCache)
	if err := scanMedia(taskContext, media, &mediaData, 0, 1); err != nil {
		return errors.Wrap(err, "single media scan")
	}

	return nil
}
