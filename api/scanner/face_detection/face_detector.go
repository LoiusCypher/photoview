package face_detection

import (
	"gorm.io/gorm"

	"github.com/loiuscypher/photoview/api/graphql/models"
)

type FaceDetector interface {
	ChangeClassifyFaceThreshold(threshold float32)
	CheckFaceGroup(db *gorm.DB, groupID int32)
	DetectFaces(db *gorm.DB, media *models.Media) error
	MergeCategories(sourceID int32, destID int32)
	MergeImageFaces(imageFaceIDs []int, destFaceGroupID int32)
	RecognizeUnlabeledFaces(db *gorm.DB, user *models.User) ([]*models.ImageFace, error)
	ReloadFacesFromDatabase(db *gorm.DB) error
}

var GlobalFaceDetector FaceDetector = nil
