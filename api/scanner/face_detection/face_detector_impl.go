//go:build !no_face_detection

package face_detection

import (
	"log"
	"math"
	"slices"
	"sort"
	"sync"

	"github.com/Kagami/go-face"
	"github.com/loiuscypher/photoview/api/graphql/models"
	"github.com/loiuscypher/photoview/api/scanner/media_encoding"
	"github.com/loiuscypher/photoview/api/utils"
	"github.com/pkg/errors"
	"gopkg.in/gographics/imagick.v3/imagick"
	"gorm.io/gorm"
)

type faceDetector struct {
	mutex             sync.Mutex
	rec               *face.Recognizer
	faceDescriptors   []face.Descriptor
	faceGroupIDs      []int32
	imageFaceIDs      []int
	classifyThreshold float32
}

func InitializeFaceDetector(db *gorm.DB) error {
	if utils.EnvDisableFaceRecognition.GetBool() {
		log.Printf("Face detection disabled (%s=1)\n", utils.EnvDisableFaceRecognition.GetName())
		return nil
	}

	log.Println("Initializing face detector")

	rec, err := face.NewRecognizer(utils.FaceRecognitionModelsPath())
	if err != nil {
		return errors.Wrap(err, "initialize facedetect recognizer")
	}

	faceDescriptors, faceGroupIDs, imageFaceIDs, err := getSamplesFromDatabase(db)
	if err != nil {
		return errors.Wrap(err, "get face detection samples from database")
	}
	rec.SetSamples(faceDescriptors, faceGroupIDs)

	var classifyThreshold float32
	site_info, err := models.GetSiteInfo(db)
	if err != nil {
		return errors.Wrap(err, "get classify face threshold from database")
	}
	classifyThreshold = float32(site_info.ClassifyFaceThreshold)

	GlobalFaceDetector = &faceDetector{
		rec:               rec,
		faceDescriptors:   faceDescriptors,
		faceGroupIDs:      faceGroupIDs,
		imageFaceIDs:      imageFaceIDs,
		classifyThreshold: classifyThreshold,
	}
	log.Println("InitializeFaceDetector: threshold, ", classifyThreshold)

	return nil
}

func getSamplesFromDatabase(db *gorm.DB) (samples []face.Descriptor, faceGroupIDs []int32, imageFaceIDs []int, err error) {

	var imageFaces []*models.ImageFace

	if err = db.Find(&imageFaces).Error; err != nil {
		return
	}

	samples = make([]face.Descriptor, len(imageFaces))
	faceGroupIDs = make([]int32, len(imageFaces))
	imageFaceIDs = make([]int, len(imageFaces))

	for i, imgFace := range imageFaces {
		samples[i] = face.Descriptor(imgFace.Descriptor)
		faceGroupIDs[i] = int32(imgFace.FaceGroupID)
		imageFaceIDs[i] = imgFace.ID
	}
	//log.Println("getSamplesFromDatabase: len(fd.descriptors):", len(samples))

	return
}

func setImageFaceSubGroup(db *gorm.DB, imageFaceID int, subGroup int) (err error) {

	if err := db.
			Where("id = ?", imageFaceID).
			Model(&models.ImageFace{}).
			UpdateColumn("subgroup", subGroup).
			Error; err != nil {
		return err
	}
	return nil
}

func setImageFaceGroup(db *gorm.DB, imageFaceID int, faceGroup int) (err error) {

	if err := db.
			Where("id = ?", imageFaceID).
			Model(&models.ImageFace{}).
			UpdateColumn("face_group_id", faceGroup).
			UpdateColumn("confirmed", false).
			UpdateColumn("subgroup", 0).
			Error; err != nil {
		return err
	}
	return nil
}

// ReloadFacesFromDatabase replaces the in-memory face descriptors with the ones in the database
func (fd *faceDetector) ReloadFacesFromDatabase(db *gorm.DB) error {
	faceDescriptors, faceGroupIDs, imageFaceIDs, err := getSamplesFromDatabase(db)
	if err != nil {
		return err
	}

	fd.mutex.Lock()
	defer fd.mutex.Unlock()

	fd.faceDescriptors = faceDescriptors
	fd.faceGroupIDs = faceGroupIDs
	fd.imageFaceIDs = imageFaceIDs
	fd.rec.SetSamples(faceDescriptors, faceGroupIDs)

	return nil
}

// ChangeClassifyFaceThreshold
func (fd *faceDetector) ChangeClassifyFaceThreshold(classifyThreshold float32) {
	fd.classifyThreshold = classifyThreshold
	log.Println("ChangeClassifyFaceThreshold: threshold, ", classifyThreshold)
}

// ReDetectFaces finds the faces in the given image and saves them to the database
func (fd *faceDetector) ReDetectFaces(db *gorm.DB, media *models.Media) error {
	if err := db.Model(media).Preload("MediaURL").First(&media).Error; err != nil {
		return err
	}

	log.Println("  ", len(media.Faces), " Faces exist already for: ", media.ID)
	return nil
}

// DetectFaces finds the faces in the given image and saves them to the database
func (fd *faceDetector) DetectFaces(db *gorm.DB, media *models.Media) error {
	return mediaDetectFaces(db, fd, media, 0)
}

// mediaDetectFaces finds the faces in the given image and saves them to the database
func mediaDetectFaces(db *gorm.DB, fd *faceDetector, media *models.Media, detection int) error {
	//if err := db.Model(media).Preload("MediaURL").Preload("Faces").First(&media).Error; err != nil {
	if err := db.Model(media).Preload("MediaURL").First(&media).Error; err != nil {
		return err
	}

	//if media.Faces == nil {
	//if err := db.Model(media).Where("Detection = ?", detection).Association("Faces").Find(&media.Faces); err != nil {
	//return err
	//}
	//}
	//log.Println("  ", len(media.Faces)," Faces exist already for: ", media.ID)

	var oldFaces []models.ImageFace
	if err := db.Where("media_id = ? AND detection = ?", media.ID, detection).Find(&oldFaces).Error; err != nil {
		return errors.Wrap(err, "error classify face")
	}
	log.Println("  existing Faces", len(oldFaces))
	for oIdx := range len( oldFaces) {
		if oldFaces[oIdx].FaceGroup == nil {
			var faceGroup models.FaceGroup
			if err := db.Model(&oldFaces[oIdx]).Association("FaceGroup").Find(&faceGroup); err != nil {
				return err
			}
			oldFaces[oIdx].FaceGroup = &faceGroup
		}
	}

	var faces []face.Face
	var img_width float64
	var img_height float64
	for _, url := range media.MediaURL {
		var thumbnailURL *models.MediaURL
		thumbnailURL = &url
		thumbnailURL.Media = media
		if thumbnailURL == nil {
			// return errors.New()
			log.Println("thumbnail url is missing: skip entry")
			continue
		}

		cachedPath, err := thumbnailURL.CachedPath()
		if err != nil {
			// return err
			log.Printf("Error: Access cacched thumbnail path %s\n", err)
			continue
		}

		if url.Purpose == models.MediaOriginal {
			// var orient int64 = 1
			// log.Printf("  Original URL found %s\n", cachedPath)
			if media.Exif != nil {
				if media.Exif.Orientation != nil {
					log.Printf("   Orientation %d  %s\n", *media.Exif.Orientation, cachedPath)
				}
			}

			mw := imagick.NewMagickWand()
			defer mw.Destroy()
			if err = mw.ReadImage(cachedPath); err != nil {
				log.Printf("Err: ReadImage %s %s\n", err, cachedPath)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			mt := mw.Clone()
			defer mt.Destroy()

			mt.ThumbnailImage(mt.GetImageWidth(),mt.GetImageHeight())
			depth := mt.GetImageDepth()
			colSp := mt.GetImageColorspace()
			itype := mt.GetImageType()
			log.Printf("   ImageColorspace %d  Imagedepth %d  ImageType %d %s\n", colSp, depth, itype, "THUMBNAIL" )
			if err = mt.AutoOrientImage(); err != nil {
				log.Printf("Err: AutoOrientImage %s\n", err)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			var blob []byte
			if blob, err = mt.GetImageBlob(); err != nil {
				log.Printf("Err: GetImageBlob %s\n", err)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			fd.mutex.Lock()
			log.Printf("RecognizeFile MediaOriginal %s\n", cachedPath)
			if faces, err = fd.rec.Recognize(blob); err != nil {
				fd.mutex.Unlock()
				log.Printf("Err: Recognize %s\n", err)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			fd.mutex.Unlock()
			return fd.reClassifyMediaFaces(db, oldFaces, faces, media, float64(mt.GetImageWidth()), float64(mt.GetImageHeight()), detection)

if 1 == 0 {
			depth := mw.GetImageDepth()
			colSp := mw.GetImageColorspace()
			itype := mw.GetImageType()
			log.Printf("   ImageColorspace %d  Imagedepth %d  ImageType %d %s\n", colSp, depth, itype, cachedPath )
			if err = mw.AutoOrientImage(); err != nil {
				log.Printf("Err: AutoOrientImage %s\n", err)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			if err = mw.SetImageFormat("JPEG"); err != nil {
				log.Printf("Err: SetImageFormat %s\n", err)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			depth = mw.GetImageDepth()
			colSp = mw.GetImageColorspace()
			itype = mw.GetImageType()
			if colSp == imagick.COLORSPACE_GRAY {
				log.Printf("   GRAY dep %d typ %d %s\n", depth, itype, cachedPath )
				cw := mw.Clone()
				//if err = mw.AddImage(cw.Clone()); err != nil {
				//	log.Printf("Err: AddImage 1 %s\n", err)
				//	return errors.Wrap(err, "error read faces")
				//	continue
				//}
				if err = mw.AddImage(cw.Clone()); err != nil {
					log.Printf("Err: AddImage 2 %s\n", err)
					return errors.Wrap(err, "error read faces")
					continue
				}
				if err = mw.AddImage(cw.Clone()); err != nil {
					log.Printf("Err: AddImage 3 %s\n", err)
					return errors.Wrap(err, "error read faces")
					continue
				}
				mw = mw.CombineImages( imagick.COLORSPACE_SRGB)
				colSp = mw.GetImageColorspace()
				depth = mw.GetImageDepth()
				itype = mw.GetImageType()
				log.Printf("   NEW  ColorSpace col %d dep %d typ %d\n", colSp, depth, itype )
				if err = mw.SetImageType( imagick.IMAGE_TYPE_TRUE_COLOR); err != nil {
					log.Printf("Err: SetImageType %s\n", err)
					return errors.Wrap(err, "error read faces")
					continue
				}
				colSp = mw.GetImageColorspace()
				depth = mw.GetImageDepth()
				itype = mw.GetImageType()
				log.Printf("   NOW  ColorSpace col %d dep %d typ %d\n", colSp, depth, itype )
			}
			if colSp != 23 {
				log.Printf("   ColorSpace %d  Depth %d\n", colSp, depth)
				if err = mw.SetImageColorspace( imagick.COLORSPACE_SRGB); err != nil {
				// if err = mw.SetImageColorspace(23); err != nil {
					log.Printf("Err: SetImageColorspace %s\n", err)
					return errors.Wrap(err, "error read faces")
					continue
				}
				// depthn := mw.GetImageDepth()
				// colSpn := mw.GetImageColorspace()
				//log.Printf("   New ColorSpace %d  Depth %d\n", colSpn, depthn)
				}
			var blob []byte
			if blob, err = mw.GetImageBlob(); err != nil {
				log.Printf("Err: GetImageBlob %s\n", err)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			fd.mutex.Lock()
			img_width = float64(mw.GetImageWidth())
			img_height = float64(mw.GetImageHeight())
			log.Printf("RecognizeFile MediaOriginal %s\n", cachedPath)
			if faces, err = fd.rec.Recognize(blob); err != nil {
				fd.mutex.Unlock()
				log.Printf("Err: Recognize %s\n", err)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			fd.mutex.Unlock()
			break
		}
}

		if url.Purpose == models.PhotoThumbnail {
			fd.mutex.Lock()
			log.Printf("RecognizeFile Thumbnail %s\n", cachedPath)
			faces, err = fd.rec.RecognizeFile(cachedPath)
			fd.mutex.Unlock()
			if err != nil {
				return errors.Wrap(err, "error read faces")
			}
			dimension, err := media_encoding.GetPhotoDimensions(cachedPath)
			img_width = float64(dimension.Width)
			img_height = float64(dimension.Height)
			if err != nil {
				return errors.Wrap(err, "error get dimensions")
			}
			continue
		}

		// if url.Purpose == models.PhotoHighRes {
		// thumbnailURL = &url
		// thumbnailURL.Media = media
		// log.Println("  Thumbnail URL found")
		// break
		// }
	}

	return fd.reClassifyMediaFaces(db, oldFaces, faces, media, float64(img_width), float64(img_height), detection)
}

func (fd *faceDetector) reClassifyMediaFaces(db *gorm.DB, oldFaces []models.ImageFace, faces []face.Face, media *models.Media, width float64, height float64, detection int) error {

	log.Println("  ", len(faces), " Faces found")
	matches := make([][]bool, len(faces))
	matchedOld := make([]int, len(oldFaces))
	matchedNew := make([]int, len(faces))
	for nIdx, nFace_ := range faces {
		matches[nIdx] = make([]bool, len(oldFaces))
		for oIdx, oFace_ := range oldFaces {
			eu_dist := euclidean_distance( nFace_.Descriptor, face.Descriptor(oFace_.Descriptor))
			matches[nIdx][oIdx] = eu_dist == 0
			if eu_dist == 0 {
				matchedOld[oIdx] += 1
				matchedNew[nIdx] += 1
				// oldFaces[oIdx] = oldFaces[len(oldFaces)-1]
				// oldFaces = oldFaces[:len(oldFaces)-1]
				//log.Println("   Keep identical face: EuD", eu_dist)
				// break
			}
		}
	}
	for nIdx, _ := range faces {
		log.Println("   matches: ", matches[nIdx])
	}
	log.Println("   matchedOld: ", matchedOld)
	if slices.Contains( matchedOld, 0) {
		//log.Println("   matchedOld: ", matchedOld)
		for oIdx, oFace_ := range oldFaces {
			if matchedOld[oIdx] == 0 {
				for nIdx, nFace_ := range faces {
					eu_dist := euclidean_distance( nFace_.Descriptor, face.Descriptor(oFace_.Descriptor))
					// log.Println("   Eucidean Dist: ", eu_dist)
					o_span_X := oFace_.Rectangle.MaxX - oFace_.Rectangle.MinX
					d_MinX := (float64(nFace_.Rectangle.Min.X) / width) - oFace_.Rectangle.MinX
					d_MaxX := (float64(nFace_.Rectangle.Max.X) / width) - oFace_.Rectangle.MaxX
					// log.Println("   DiffRectangle: X Span ", o_span_X, " Diff ", d_MinX, d_MaxX)
					o_span_Y := oFace_.Rectangle.MaxY - oFace_.Rectangle.MinY
					d_MinY := (float64(nFace_.Rectangle.Min.Y) / height) - oFace_.Rectangle.MinY
					d_MaxY := (float64(nFace_.Rectangle.Max.Y) / height) - oFace_.Rectangle.MaxY
					 log.Println("   DiffRectangle: Y Span ", o_span_Y, " Diff ", d_MinY, d_MaxY)
					p_MinX := 100 * max(d_MinX,-d_MinX) / o_span_X
					p_MaxX := 100 * max(d_MaxX,-d_MaxX) / o_span_X
					p_MinY := 100 * max(d_MinY,-d_MinY) / o_span_Y
					p_MaxY := 100 * max(d_MaxY,-d_MaxY) / o_span_Y
					log.Println("   DiffPercent: X ", p_MinX, p_MaxX, "  Y ", p_MinY, p_MaxY)
					p_MovX :=  100 * max(d_MaxX+d_MinX,-d_MaxX-d_MinX) / o_span_X / 2.0
					p_MovY :=  100 * max(d_MaxY+d_MinY,-d_MaxY-d_MinY) / o_span_Y / 2.0
					log.Println("   MovePercent: X ", p_MovX, "  Y ", p_MovY,)
					max_perc_diff := float64(40.0)
					max_perc_move := float64(15.0)
					if p_MovX < max_perc_move && p_MovY < max_perc_move && p_MinX < max_perc_diff && p_MaxX < max_perc_diff && p_MinY < max_perc_diff && p_MaxY < max_perc_diff {
						log.Println("   Eucidean Dist: ", eu_dist)
						log.Println("delete: ", oIdx, "  same space as:", nIdx)
						if err := db.Where("id = ?", oFace_.ID).Delete(models.ImageFace{}).Error; err != nil {
							return errors.Wrap(err, "error deleting face")
						}
						break // we deleted this unmatched one, go to next unmatched
					}
				}
			}
		}
	}
	log.Println("   matchedNew: ", matchedNew)
	if slices.Contains( matchedNew, 0) {
		log.Println("   matchedNew: ", matchedNew)
	}
	for nIdx, nFace := range faces {
		if matchedNew[nIdx] == 0 {
			log.Println("NEW face found: ", nIdx)
			if err := fd.classifyFace(db, &nFace, media, width, height, detection); err != nil {
				return errors.Wrap(err, "error classify face")
			}
		}
		if matchedNew[nIdx] > 1 {
			var keep int = -1
			for oIdx, oFace_ := range oldFaces {
				if matches[nIdx][oIdx] {
					log.Println("  Old FaceID: ", oFace_.ID, "  FaceGroup: ", oFace_.FaceGroupID, oFace_.FaceGroup != nil)
					log.Println("   Rectangle: ", oFace_.Rectangle.MinX, oFace_.Rectangle.MinY, oFace_.Rectangle.MaxX, oFace_.Rectangle.MaxY)
					log.Println("   Detection: ", oFace_.Detection, "  Confirmed: ", oFace_.Confirmed, "  Subgroup: ", oFace_.Subgroup)
					if oFace_.FaceGroup.Label == nil {
						log.Println("   Faceoup ID: ", oFace_.FaceGroup.ID, " Label: ", oFace_.FaceGroup.Label, " oIdx: ", oIdx)
					} else {
						log.Println("   Faceoup ID: ", oFace_.FaceGroup.ID, " Label: ", *oFace_.FaceGroup.Label, " oIdx: ", oIdx)
					}
					if keep < 0 {
						keep = oIdx
						log.Println("keep 1: ", keep)
					} else {
						if oFace_.FaceGroup.Label == nil {
							log.Println("delete 1: ", oIdx, "  keep:", keep)
							if err := db.Where("id = ?", oFace_.ID).Delete(models.ImageFace{}).Error; err != nil {
								return errors.Wrap(err, "error deletim#ng face")
							}
						} else {
							log.Println("check keep 1: ", keep, oldFaces[keep].FaceGroup == nil)
							if oldFaces[keep].FaceGroup.Label == nil || !oFace_.Confirmed {
								log.Println("delete 2: ", keep, "  keep:", oIdx)
								if err := db.Where("id = ?", oldFaces[keep].ID).Delete(models.ImageFace{}).Error; err != nil {
									return errors.Wrap(err, "error deleting face")
								}
								keep = oIdx
							} else {
								log.Println("delete 3: ", oIdx, "  keep:", keep)
								if err := db.Where("id = ?", oFace_.ID).Delete(models.ImageFace{}).Error; err != nil {
									return errors.Wrap(err, "error deleting face")
								}
							}
						}
					}
				}
			}
		}
	}

	fd.ReloadFacesFromDatabase(db)

	return nil
}

func (fd *faceDetector) classifyFace(db *gorm.DB, face *face.Face, media *models.Media, width float64, height float64, detection int) error {
	fd.mutex.Lock()
	defer fd.mutex.Unlock()

	log.Println("classifyFace threshold: ", fd.classifyThreshold, "len(fd.descriptors):", len(fd.faceDescriptors))
	match := int32(fd.rec.ClassifyThreshold(face.Descriptor, fd.classifyThreshold))

	imageFace := models.ImageFace{
		MediaID:    media.ID,
		Descriptor: models.FaceDescriptor(face.Descriptor),
		Rectangle: models.FaceRectangle{
			// Converts a pixel absolute rectangle to a relative FaceRectangle.
			MinX: float64(face.Rectangle.Min.X) / width,
			MaxX: float64(face.Rectangle.Max.X) / width,
			MinY: float64(face.Rectangle.Min.Y) / height,
			MaxY: float64(face.Rectangle.Max.Y) / height,
		},
		Detection: detection,
	}
	//log.Printf("    Face region: %f-%f:%f-%f", imageFace.Rectangle.MinX, imageFace.Rectangle.MaxX, imageFace.Rectangle.MinY, imageFace.Rectangle.MaxX)

	var faceGroup models.FaceGroup

	// If no match add it new to samples
	if match < 0 {
		log.Println("     No match, assigning new face")

		faceGroup = models.FaceGroup{
			ImageFaces: []models.ImageFace{imageFace},
		}

		if err := db.Create(&faceGroup).Error; err != nil {
			return err
		}

	} else {
		log.Println("     Found match", match)

		if err := db.First(&faceGroup, int(match)).Error; err != nil {
			return err
		}

		if err := db.Model(&faceGroup).Association("ImageFaces").Append(&imageFace); err != nil {
			return err
		}
	}

	fd.faceDescriptors = append(fd.faceDescriptors, face.Descriptor)
	fd.faceGroupIDs = append(fd.faceGroupIDs, int32(faceGroup.ID))
	fd.imageFaceIDs = append(fd.imageFaceIDs, imageFace.ID)

	fd.rec.SetSamples(fd.faceDescriptors, fd.faceGroupIDs)
	return nil
}

func (fd *faceDetector) MergeCategories(sourceID int32, destID int32) {
	fd.mutex.Lock()
	defer fd.mutex.Unlock()

	for i := range fd.faceGroupIDs {
		if fd.faceGroupIDs[i] == sourceID {
			fd.faceGroupIDs[i] = destID
		}
	}
}

func (fd *faceDetector) MergeImageFaces(imageFaceIDs []int, destFaceGroupID int32) {
	fd.mutex.Lock()
	defer fd.mutex.Unlock()

	for i := range fd.faceGroupIDs {
		imageFaceID := fd.imageFaceIDs[i]

		if slices.Contains(imageFaceIDs, imageFaceID) {
			fd.faceGroupIDs[i] = destFaceGroupID
		}
	}
}

func (fd *faceDetector) RecognizeUnlabeledFaces(tx *gorm.DB, user *models.User) ([]*models.ImageFace, error) {
	unrecognizedDescriptors := make([]face.Descriptor, 0)
	unrecognizedFaceGroupIDs := make([]int32, 0)
	unrecognizedImageFaceIDs := make([]int, 0)

	newFaceGroupIDs := make([]int32, 0)
	newDescriptors := make([]face.Descriptor, 0)
	newImageFaceIDs := make([]int, 0)

	var unlabeledFaceGroups []*models.FaceGroup

	err := tx.
		Joins("JOIN image_faces ON image_faces.face_group_id = face_groups.id").
		Joins("JOIN media ON image_faces.media_id = media.id").
		Where("face_groups.label IS NULL").
		Where("media.album_id IN (?)",
			tx.Select("album_id").Table("user_albums").Where("user_id = ?", user.ID),
		).
		Find(&unlabeledFaceGroups).Error

	if err != nil {
		return nil, err
	}

	fd.mutex.Lock()
	defer fd.mutex.Unlock()

	for i := range fd.faceDescriptors {
		descriptor := fd.faceDescriptors[i]
		faceGroupID := fd.faceGroupIDs[i]
		imageFaceID := fd.imageFaceIDs[i]

		isUnlabeled := false
		for _, unlabeledFaceGroup := range unlabeledFaceGroups {
			if faceGroupID == int32(unlabeledFaceGroup.ID) {
				isUnlabeled = true
				continue
			}
		}

		if isUnlabeled {
			unrecognizedFaceGroupIDs = append(unrecognizedFaceGroupIDs, faceGroupID)
			unrecognizedDescriptors = append(unrecognizedDescriptors, descriptor)
			unrecognizedImageFaceIDs = append(unrecognizedImageFaceIDs, imageFaceID)
		} else {
			newFaceGroupIDs = append(newFaceGroupIDs, faceGroupID)
			newDescriptors = append(newDescriptors, descriptor)
			newImageFaceIDs = append(newImageFaceIDs, imageFaceID)
		}
	}

	fd.faceGroupIDs = newFaceGroupIDs
	fd.faceDescriptors = newDescriptors
	fd.imageFaceIDs = newImageFaceIDs
	fd.rec.SetSamples(fd.faceDescriptors, fd.faceGroupIDs)

	updatedImageFaces := make([]*models.ImageFace, 0)

	for i := range unrecognizedDescriptors {
		descriptor := unrecognizedDescriptors[i]
		faceGroupID := unrecognizedFaceGroupIDs[i]
		imageFaceID := unrecognizedImageFaceIDs[i]

		match := int32(fd.rec.ClassifyThreshold(descriptor, fd.classifyThreshold))

		if match < 0 {
			// still no match, we can readd it to the list
			fd.faceGroupIDs = append(fd.faceGroupIDs, faceGroupID)
			fd.faceDescriptors = append(fd.faceDescriptors, descriptor)
			fd.imageFaceIDs = append(fd.imageFaceIDs, imageFaceID)
		} else {
			// found new match, update the database
			var imageFace models.ImageFace
			if err := tx.Model(&models.ImageFace{}).First(imageFace, imageFaceID).Error; err != nil {
				return nil, err
			}

			if err := tx.Model(&imageFace).Update("face_group_id", int(faceGroupID)).Error; err != nil {
				return nil, err
			}

			updatedImageFaces = append(updatedImageFaces, &imageFace)

			fd.faceGroupIDs = append(fd.faceGroupIDs, match)
			fd.faceDescriptors = append(fd.faceDescriptors, descriptor)
			fd.imageFaceIDs = append(fd.imageFaceIDs, imageFaceID)
		}
		fd.rec.SetSamples(fd.faceDescriptors, fd.faceGroupIDs)
	}

	return updatedImageFaces, nil
}

func (fd *faceDetector) findFaceGroupCandidates(faceGroup int, checkLng int) {

	fd.mutex.Lock()
	defer fd.mutex.Unlock()

	k := len(fd.faceDescriptors) - checkLng
	log.Printf("findFaceGroupCandidates Id %d Lng %d start %d chkLng %d\n", faceGroup, len(fd.faceDescriptors), k, checkLng)

	descriptor := fd.faceDescriptors[k]
	faceGroupID := fd.faceGroupIDs[k]
	imageFaceID := fd.imageFaceIDs[k]
	fd.faceDescriptors = slices.Delete(fd.faceDescriptors, k, k+1)
	fd.faceGroupIDs = slices.Delete(fd.faceGroupIDs, k, k+1)
	fd.imageFaceIDs = slices.Delete(fd.imageFaceIDs, k, k+1)
	for j := range checkLng {
		i := j + k
		// log.Printf("findFaceGroupCandidates index %d\n", i)
		fd.rec.SetSamples(fd.faceDescriptors, fd.faceGroupIDs)
		match := fd.rec.ClassifyThreshold(descriptor, fd.classifyThreshold)
		if match == faceGroup {
			log.Printf("Face %d matches desired face group %d \n", imageFaceID, match)
		} else {
			if int32(match) == faceGroupID {
				log.Printf("Face %d confirms former group %d \n", imageFaceID, match)
			} else {
				log.Printf("Face %d matches completely different group %d(%d) \n", imageFaceID, match, faceGroupID)
			}
		}
		t_descriptor := fd.faceDescriptors[i]
		t_faceGroupID := fd.faceGroupIDs[i]
		t_imageFaceID := fd.imageFaceIDs[i]
		fd.faceDescriptors[i] = descriptor
		fd.faceGroupIDs[i] = faceGroupID
		fd.imageFaceIDs[i] = imageFaceID
		descriptor = t_descriptor
		faceGroupID = t_faceGroupID
		imageFaceID = t_imageFaceID
	}
	fd.faceDescriptors = append(fd.faceDescriptors, descriptor)
	fd.faceGroupIDs = append(fd.faceGroupIDs, faceGroupID)
	fd.imageFaceIDs = append(fd.imageFaceIDs, imageFaceID)
	fd.rec.SetSamples(fd.faceDescriptors, fd.faceGroupIDs)
}

func (fd *faceDetector) ConfirmedFaceGroups(db *gorm.DB) ([]int, error) {
	var groupIDs []int
	if err := db.
		Model(models.ImageFace{}).
		// Debug().
		Where("confirmed = TRUE").
		Group("face_group_id").
		Order("COUNT(face_group_id)").
		Pluck("face_group_id", &groupIDs).
		Error; err != nil {

		return make([]int, 0), err
	}

	// for _, groupId := range groupIDs {
	// log.Printf("Group ID %d with confirmed faces\n", groupId)
	// }
	return groupIDs, nil
}

func (fd *faceDetector) FaceGroupFaceIDs(db *gorm.DB, confirmed bool, groupId int) ([]int, error) {

	var faceIDs []int
	if err := db.
		Model(models.ImageFace{}).
		// Debug().
		Where("face_group_id = ?", groupId).
		Where("confirmed = ?", confirmed).
		Pluck("ID", &faceIDs).
		Error; err != nil {

		return make([]int, 0), err
	}

	var faceGroup models.FaceGroup
	if err := db.
		// Debug().
		Model(models.FaceGroup{}).
		Where("id = ?", groupId).
		First(&faceGroup).
		Error; err != nil {

		return make([]int, 0), err
	}
	// log.Println("Face ", imgFace)
	// log.Println("FaceGroup* ", *imgFace.FaceGroup)
	// log.Println("FaceGroup ", imgFace.FaceGroup)
	log.Printf("FaceGroup '%s'\n", *faceGroup.Label)
	return faceIDs, nil
}

func (fd *faceDetector) CheckFaceGroup(db *gorm.DB, groupID int32) {
	grps, _ := fd.ConfirmedFaceGroups(db)
	log.Printf("Groups with confirmed faces %d\n", len(grps))
	for _, grpID := range grps {
		log.Printf("Group ID %d\n", grpID)
		confFaceIDs, _ := fd.FaceGroupFaceIDs(db, true, grpID)
		log.Printf("Group %d with %d confirmed faces\n", grpID, len(confFaceIDs))
		// for _, faceID := range confFaceIDs {
		// log.Printf("GrpID %d faceID %d \n", grpID, faceID)
		// }
		faceIDs, _ := fd.FaceGroupFaceIDs(db, false, grpID)
		log.Printf("Group %d with %d unconfirmed faces\n", grpID, len(faceIDs))
		// for _, faceID := range faceIDs {
		// log.Printf("GrpID %d faceID %d \n", grpID, faceID)
		// }
		fd.nextFaceForGroup(db, grpID, faceIDs)
	}
}

func (fd *faceDetector) nextFaceForGroup(db *gorm.DB, groupID int, faceIDs []int) {

	log.Printf("nextFaceForGroup %d %d\n", groupID, len(faceIDs))

	saveFaceDescriptors := fd.faceDescriptors
	saveFaceGroupIDs := fd.faceGroupIDs
	saveImageFaceIDs := fd.imageFaceIDs

	for i, faceID := range faceIDs {
		idx := slices.IndexFunc(fd.imageFaceIDs, func(c int) bool { return c == faceID })
		// log.Printf("Index %d\n", idx)
		if idx < 0 {
			log.Printf("Error slicing Facegroups with idx %d(%d) ID %d", i, len(faceIDs), faceID)
			return
		}
		fd.faceDescriptors = append(fd.faceDescriptors, fd.faceDescriptors[idx])
		fd.faceDescriptors = slices.Delete(fd.faceDescriptors, idx, idx+1)
		fd.faceGroupIDs = append(fd.faceGroupIDs, 0)
		fd.faceGroupIDs = slices.Delete(fd.faceGroupIDs, idx, idx+1)
		fd.imageFaceIDs = append(fd.imageFaceIDs, faceID)
		fd.imageFaceIDs = slices.Delete(fd.imageFaceIDs, idx, idx+1)
	}
	fd.findFaceGroupCandidates(groupID, len(faceIDs))
	fd.faceDescriptors = saveFaceDescriptors
	fd.faceGroupIDs = saveFaceGroupIDs
	fd.imageFaceIDs = saveImageFaceIDs
	fd.rec.SetSamples(fd.faceDescriptors, fd.faceGroupIDs)
}

type distElem struct {
	sId  int
	tId  int
	dist float64
}

func deleteUnused(grp []int, newId int, dists []distElem) (r []distElem) {
	for _, grpId := range grp {
		dists = slices.DeleteFunc(dists, func(c distElem) bool {
			return (c.sId == grpId) && (c.tId == newId) || (c.tId == grpId) && (c.sId == newId)
		})
	}
	return dists
}

func mergeGroupsIfPossible(groups [][]int, dists []distElem, tId int, sGrpIndex int) (r bool, rgroups [][]int, rdists []distElem) {
	for j, grpT := range groups {
		if 0 <= slices.IndexFunc(grpT, func(c int) bool { return c == tId }) {
			log.Printf("   found two related groups that must be merged %d %d\n", sGrpIndex, j)
			// merge groups because they are close to each other
			for _, faceId := range grpT {
				dists = deleteUnused(groups[sGrpIndex], faceId, dists)
			}
			//groups[sGrpIndex] = slices.Concat(groups[sGrpIndex], grpT)
			groups[sGrpIndex] = slices.Concat(groups[sGrpIndex], groups[j])
			log.Println("   7a group", sGrpIndex, groups[sGrpIndex])
			groupn := slices.Delete(groups, j, j+1)
			log.Println("   7b groupCnt", len(groupn))
			return true, groupn, dists
		}
	}
	log.Printf("   No group merge to group %d\n", sGrpIndex)
	return false, groups, dists
}

func (fd *faceDetector) SplitFaceGroup(db *gorm.DB, groupID int32) {

	log.Printf("SplitFaceGroup groupID %d\n", groupID)

	fd.ReloadFacesFromDatabase(db)

	var faceGrp models.FaceGroup
	if err := db.Model(models.FaceGroup{}).Where("id = ?", groupID).First(&faceGrp).Error; err != nil {
		log.Println(" SplitFaceGroup err", err)
		return
	}

	if faceGrp.Label != nil {
		log.Println(" SearchFacesForFaceGroup for", *faceGrp.Label)
		fd.searchFacesForFaceGroup(db, groupID)
	} else {
		log.Println(" SplitFaceGroup Label is nii")
		fd.moveToClosestFaceGroup(db, groupID, false)
	}
}

func (fd *faceDetector) OrderFacesOfFaceGroup(db *gorm.DB, groupID int32) {

	log.Printf("OrderFacesOfFaceGroup groupID %d\n", groupID)

	fd.ReloadFacesFromDatabase(db)

	var faceGrp models.FaceGroup
	if err := db.Model(models.FaceGroup{}).Where("id = ?", groupID).First(&faceGrp).Error; err != nil {
		log.Println(" SplitFaceGroup err", err)
		return
	}

	if faceGrp.Label != nil {
		log.Println(" SplitFaceGroup ", *faceGrp.Label)
		fd.reOrderFaceGroupUnconfirmed(db, groupID)
	} else {
		log.Println(" SplitFaceGroup Label is nii")
		fd.moveToClosestFaceGroup(db, groupID, true)
	}
}

func (fd *faceDetector) findFaceGroupIndices(groupID int32) ( []int, []int) {
	grpIndexes := make([]int, 0)
	noGrpIdxes := make([]int, 0)
	for i, id := range fd.faceGroupIDs {
		if id == groupID {
			grpIndexes = append( grpIndexes, i)
			// log.Println(" findFaceGroupIndices", "index", i, "grp id", id, "faceGroupID", fd.faceGroupIDs[i], "imgFaceID", fd.imageFaceIDs[i], "grpIndexes", grpIndexes)
		} else {
			noGrpIdxes = append( noGrpIdxes, i)
		}
	}
	log.Println(" findFaceGroupIndices", grpIndexes, len(noGrpIdxes))
	return grpIndexes, noGrpIdxes
}

func (fd *faceDetector) findConfirmedFaceIndices( db *gorm.DB, groupIdIndexes []int, confirmed bool) ( []int) {
	confIndexes := make([]int, 0)
	for _, idx := range groupIdIndexes {
		var conf []bool
		//if err := db.Model(models.ImageFace{ID: fd.imageFaceIDs[idx]}).Select("confirmed").Take(&conf).Error; err != nil {
		//if err := db.Model(models.ImageFace{ID: fd.imageFaceIDs[idx]}).Pluck("confirmed",&conf).Error; err != nil {
		if err := db.Model(models.ImageFace{}).Where("id = ?", fd.imageFaceIDs[idx]).Pluck("confirmed",&conf).Error; err != nil {
		//var imgface models.ImageFace
		//if err := db.Model(imageFace, fd.imageFaceIDs[idx]).Take(imageFace).Error; err != nil {
			log.Println(" findConfirmedFaceIndices err", err)
			return nil
		}
		if confirmed {
			// log.Println(" findConfirmedFaceIndices", idx, confirmed, conf[0], "faceGroupID", fd.faceGroupIDs[idx], "imgFaceID", fd.imageFaceIDs[idx])
		}
		if conf[0] == confirmed {
			confIndexes = append( confIndexes, idx)
		}
	}
	if confirmed {
		log.Println(" findConfirmedFaceIndices", confIndexes, len(confIndexes))
	} else {
		log.Println(" findConfirmedFaceIndices", len(confIndexes))
	}
	return confIndexes
}

func loadFaceGroupWithDescriptor(db *gorm.DB, groupSel bool, groupID int32, confirmed bool) ( []*models.ImageFace, []face.Descriptor) {
	var grpOtherImageFaces []*models.ImageFace
	grpCond := "face_group_id = ?"
	if !groupSel {
		grpCond = "face_group_id != ?"
	}
	if err := db.Where("confirmed = ?", confirmed).Where(grpCond,groupID).Find(&grpOtherImageFaces).Error; err != nil {
		return nil, nil
	}
	log.Println("Loaded Faces", len(grpOtherImageFaces))
	otherDesc := make([]face.Descriptor, len(grpOtherImageFaces))
	for i, imgFace := range grpOtherImageFaces {
		otherDesc[i] = face.Descriptor(imgFace.Descriptor)
	}
	return grpOtherImageFaces, otherDesc
}

func calcImgMinDistGrp( imgDesc face.Descriptor, groupDesc []face.Descriptor) (float64, int) {
	var minDist float64
	var minDistIdx int
	for cIdx, confImageFace := range groupDesc {
		curDist := euclidean_distance( imgDesc, confImageFace)
		if (cIdx == 0) || (curDist < minDist) {
			minDist = curDist
			minDistIdx = cIdx
		}
		// log.Printf("cIdx %d minDist %d %f\n", cIdx, minDistIdx, minDist)
	}
	return minDist, minDistIdx
}

func (fd *faceDetector) calcMinDistIndex( imgIndex int, otherIndexes []int) (float64, int) {
	var minDist float64
	var minDistIdx int
	for cIdx, compIndex := range otherIndexes {
		curDist := euclidean_distance( fd.faceDescriptors[imgIndex], fd.faceDescriptors[compIndex])
		if (cIdx == 0) || (curDist < minDist) {
			minDist = curDist
			minDistIdx = compIndex
		}
		// log.Printf("cIdx %d minDist %d %f\n", cIdx, minDistIdx, minDist)
	}
	return minDist, minDistIdx
}

func (fd *faceDetector) searchFacesForFaceGroup(db *gorm.DB, groupID int32) {

	log.Printf("searchFacesForFaceGroup groupID %d\n", groupID)

	grpIndexes, noGroupIndexes := fd.findFaceGroupIndices( groupID)
	log.Println("in group:", len(grpIndexes), "out group:", len(noGroupIndexes))
	grpConfIndexes := fd.findConfirmedFaceIndices( db, grpIndexes, true)
	log.Println("Confirmed", len(grpConfIndexes))
	noGrpUnConfIndexes := fd.findConfirmedFaceIndices( db, noGroupIndexes, false)
	log.Println("Other Group", len(noGrpUnConfIndexes))

	grpConfirmedImageFaces, confDesc := loadFaceGroupWithDescriptor(db, true, groupID, true)
	// log.Println("Confirmed old", len(grpConfIndexes) == len(grpConfirmedImageFaces), len(grpConfirmedImageFaces))
	if len(grpConfIndexes) != len(grpConfirmedImageFaces) {
		log.Println("Confirmed old DIFFERS", len(grpConfIndexes) == len(grpConfirmedImageFaces), len(grpConfirmedImageFaces))
	}

	grpOtherImageFaces, otherDesc := loadFaceGroupWithDescriptor(db, false, groupID, false)
	// log.Println("Other Group old", len(noGrpUnConfIndexes) == len(grpOtherImageFaces), len(grpOtherImageFaces))
	if len(noGrpUnConfIndexes) != len(grpOtherImageFaces) {
		log.Println("Other Group old DIFFERS", len(noGrpUnConfIndexes) == len(grpOtherImageFaces), len(grpOtherImageFaces))
	}

	for i, idx := range grpConfIndexes {
		if int32(fd.imageFaceIDs[idx]) != int32(grpConfirmedImageFaces[i].ID) {
			log.Println("Check", int32(fd.faceGroupIDs[idx]) == int32(grpConfirmedImageFaces[i].ID), fd.faceGroupIDs[idx], grpConfirmedImageFaces[i].ID)
		}
	}

	for i, idx := range noGrpUnConfIndexes {
		// log.Println("Check", fd.faceGroupIDs[idx] == int32(grpOtherImageFaces[i].ID), fd.faceGroupIDs[idx], grpOtherImageFaces[i].ID)
		if int32(fd.faceGroupIDs[idx]) != int32(grpOtherImageFaces[i].ID) {
			log.Println("Check", int32(fd.faceGroupIDs[idx]) == int32(grpOtherImageFaces[i].ID), fd.faceGroupIDs[idx], grpOtherImageFaces[i].ID)
		}
	}

	index := make([]int, len(grpOtherImageFaces))
	odists := make([]float64, len(grpOtherImageFaces))
	for uIdx, otherImageFace := range otherDesc {
		index[uIdx] = uIdx
		odists[uIdx], _ = calcImgMinDistGrp( otherImageFace, confDesc)
	}
	// index := make([]int, len(noGrpUnConfIndexes))
	dists := make([]float64, len(noGrpUnConfIndexes))
	for uIdx, oIdx := range noGrpUnConfIndexes {
		// index[uIdx] = oIdx
		// log.Println("Check", index[uIdx] == oIdx, index[uIdx], oIdx)
		// if index[uIdx] != oIdx {
			// log.Println("Check", index[uIdx] == oIdx, index[uIdx], oIdx)
		// }
		dists[uIdx], _ = fd.calcMinDistIndex( oIdx, grpConfIndexes)
		// log.Println("Check", index[uIdx] == oIdx, index[uIdx], oIdx)
		if odists[uIdx] != dists[uIdx] {
			log.Println("Check", odists[uIdx] == oIdx, index[uIdx], oIdx)
		}
	}

	// log.Println("dists", dists)
	sort.Slice( index, func( i, j int) bool {
		return dists[index[i]] < dists[index[j]]
	})
	// log.Println("index", index)

	for uIdx := range min( 10, len(grpOtherImageFaces)) {
		var faceGrp models.FaceGroup
		if err := db.Model(models.FaceGroup{}).Where("id = ?", grpOtherImageFaces[index[uIdx]].FaceGroupID).First(&faceGrp).Error; err != nil {
			log.Println(" searchFacesForFaceGroup err", err)
			return
		}
		var media models.Media
		if err := db.Model(models.Media{}).Where("id = ?", grpOtherImageFaces[index[uIdx]].MediaID).First(&media).Error; err != nil {
			log.Println(" searchFacesForFaceGroup err", err)
			return
		}
		label := "Label is nil"
		if faceGrp.Label != nil {
			label = *faceGrp.Label
		}
		//_, oldConfDesc := loadFaceGroupWithDescriptor(db, true, int32(faceGrp.ID), true)
		_, oldConfDesc := loadFaceGroupWithDescriptor(db, true, int32(grpOtherImageFaces[index[uIdx]].FaceGroupID), true)
		oldConfDist, _ := calcImgMinDistGrp( otherDesc[index[uIdx]], oldConfDesc)
		log.Println("Img", grpOtherImageFaces[index[uIdx]].ID, "Media", media.ID, "Album", media.AlbumID, "Group", grpOtherImageFaces[index[uIdx]].FaceGroupID, label, "distance", dists[index[uIdx]], "oldConfDist", oldConfDist)
		if (dists[index[uIdx]] < 0.250) && (len(oldConfDesc) == 0 || (dists[index[uIdx]] < (0.75 * oldConfDist))) {
			if err := setImageFaceGroup(db, grpOtherImageFaces[index[uIdx]].ID, int(groupID)); err != nil {
				log.Println(" searchFacesForFaceGroup err", err)
				return
			}
			log.Println(" Moved Img to current GROUP")
		}
	}
}

func (fd *faceDetector) moveToClosestFaceGroup(db *gorm.DB, groupID int32, doIt bool) {

	log.Printf("moveToClosestFaceGroup groupID %d\n", groupID)

	var grpImageFaces []*models.ImageFace
	if err := db.Where("face_group_id = ?",groupID).Find(&grpImageFaces).Error; err != nil {
		return
	}
	log.Println("Group", len(grpImageFaces))
	grpDesc := make([]face.Descriptor, len(grpImageFaces))
	for i, imgFace := range grpImageFaces {
		grpDesc[i] = face.Descriptor(imgFace.Descriptor)
	}

	//grpConfirmedImageFaces, confDesc := loadFaceGroupWithDescriptor(db, true, groupID, false)
	//log.Println("Confirmed", len(grpConfirmedImageFaces))

	var otherImageFaces []*models.ImageFace
	if err := db.Where("face_group_id != ?",groupID).Find(&otherImageFaces).Error; err != nil {
		return
	}
	log.Println("Other", len(otherImageFaces))
	otherDesc := make([]face.Descriptor, len(otherImageFaces))
	for i, imgFace := range otherImageFaces {
		otherDesc[i] = face.Descriptor(imgFace.Descriptor)
	}
	//grpConfirmedImageFaces, confDesc := loadFaceGroupWithDescriptor(db, false, groupID, true)
	//log.Println("Confirmed", len(grpConfirmedImageFaces))


	index := make([]int, len(grpDesc))
	dists := make([]float64, len(grpDesc))
	distIdx := make([]int, len(grpDesc))
	for gIdx, grpImageFace := range grpDesc {
		var minDist float64
		var minDistIdx int
		for oIdx, otherImageFace := range otherDesc {
			curDist := euclidean_distance( grpImageFace, otherImageFace)
			if (oIdx == 0) || (curDist < minDist) {
				minDist = curDist
				minDistIdx = oIdx
			}
			// log.Printf("gIdx %d oIdx %d minDist %f\n", gIdx, oIdx, minDist)
		}
		log.Printf("gIdx %3d ID %6d minDistIdx %6d minDist %f\n", gIdx, grpImageFaces[gIdx].ID, minDistIdx, minDist)
		index[gIdx] = gIdx
		dists[gIdx] = minDist
		distIdx[gIdx] = minDistIdx

		minDist, minDistIdx = calcImgMinDistGrp( grpImageFace, otherDesc)
		if minDist != dists[gIdx] || minDistIdx != distIdx[gIdx] {
			log.Printf("gIdx %3d ID %6d minDist %f %d / %f %d\n", gIdx, grpImageFaces[gIdx].ID, minDist, minDistIdx, dists[gIdx], distIdx[gIdx])
		}
	}
	log.Println("dists", dists)
	sort.Slice( index, func( i, j int) bool {
		return dists[index[i]] < dists[index[j]]
	})
	log.Println("index", index)

	for gIdx := range len(grpImageFaces) {
		otherImg := otherImageFaces[distIdx[index[gIdx]]]
		if otherImg.FaceGroup == nil {
			var faceGroup models.FaceGroup
			if err := db.Model(&otherImg).Association("FaceGroup").Find(&faceGroup); err != nil {
				return
			}
			otherImageFaces[distIdx[index[gIdx]]].FaceGroup = &faceGroup
			otherImg = otherImageFaces[distIdx[index[gIdx]]]
		}
		if otherImg.FaceGroup.Label == nil {
			log.Println("FaceImageId", otherImg.ID, "FaceGroupId", otherImg.FaceGroupID, "dist", dists[index[gIdx]])
		} else {
			log.Println("FaceImageId", otherImg.ID, "FaceGroupId", otherImg.FaceGroupID, "dist", dists[index[gIdx]], *otherImg.FaceGroup.Label)
		}
		if doIt && dists[index[gIdx]] < 0.200 {
			if err := setImageFaceGroup(db, grpImageFaces[gIdx].ID, otherImg.FaceGroupID); err != nil {
				log.Println(" moveToClosestFaceGroup err", err)
				return
			}
		}
	}
}

func (fd *faceDetector) reOrderFaceGroupUnconfirmed(db *gorm.DB, groupID int32) {

	log.Printf("reOrderFaceGroupUnconfirmed groupID %d\n", groupID)

	_, confDesc := loadFaceGroupWithDescriptor(db, true, groupID, true)
	log.Println("Confirmed", len(confDesc))

	grpUnConfirmedImageFaces, unConfDesc := loadFaceGroupWithDescriptor(db, true, groupID, false)
	log.Println("Other Group", len(unConfDesc))

	index := make([]int, len(unConfDesc))
	dists := make([]float64, len(unConfDesc))
	for uIdx, unConfImageFace := range unConfDesc {
		index[uIdx] = uIdx
		dists[uIdx], _ = calcImgMinDistGrp( unConfImageFace, confDesc)
		log.Printf("uIdx %3d ID %6d minDist %f\n", uIdx, grpUnConfirmedImageFaces[uIdx].ID, dists[uIdx])
	}
	// log.Println("dists", dists)
	sort.Slice( index, func( i, j int) bool {
		return dists[index[i]] < dists[index[j]]
	})
	// log.Println("index", index)

	for uIdx := range len(unConfDesc) {
		var media models.Media
		if err := db.Model(models.Media{}).Where("id = ?", grpUnConfirmedImageFaces[index[uIdx]].MediaID).First(&media).Error; err != nil {
			log.Println(" reOrderFaceImages err", err)
			return
		}
		log.Println( "distance for FaceImageId", grpUnConfirmedImageFaces[index[uIdx]].ID, dists[index[uIdx]], "Media", media.ID, "Album", media.AlbumID)
		if err := setImageFaceSubGroup(db, grpUnConfirmedImageFaces[index[uIdx]].ID, uIdx); err != nil {
			log.Println("Saved Subgroup", uIdx, "error saving Subgroup for FaceImageId", grpUnConfirmedImageFaces[index[uIdx]].ID, err)
		}
	}
}

func (fd *faceDetector) SplitFaceGroupOld(db *gorm.DB, groupID int32) {

	log.Printf("splitFaceGroup groupID %d\n", groupID)

	fd.ReloadFacesFromDatabase(db)

	// collect descriptors and faceIDs of groupID and calculate all distances in this group
	descr := make([]face.Descriptor, 0)
	imIds := make([]int, 0)
	dists := make([]distElem, 0)
	for i, faceGroupID := range fd.faceGroupIDs {
		if faceGroupID == groupID {
			newId := fd.imageFaceIDs[i]
			for j, dstId := range imIds {
				var elem distElem
				elem.sId = newId
				elem.tId = dstId
				elem.dist = euclidean_distance(descr[j], fd.faceDescriptors[i])
				log.Printf(" loop i:%d j:%d newId:%d dstId:%d dist:%f\n", i, j, newId, dstId, elem.dist)
				dists = append(dists, elem)
			}
			descr = append(descr, fd.faceDescriptors[i])
			imIds = append(imIds, newId)
		}
	}

	// log.Printf("splitFaceGroup len(descr) %d len(dists) %d\n", len(descr), len(dists))

	// order distanes ascending
	sort.Slice(dists, func(i, j int) bool {
		return dists[i].dist < dists[j].dist
	})

	var groups [][]int
	for len(dists) > 0 {
		log.Printf("0 len(dists) %d len(groups) %d\n", len(dists), len(groups))
		nextDist := dists[0]
		dists = dists[1:]
		log.Printf(" 1 next distance between %d and %d is %f\n", nextDist.sId, nextDist.tId, nextDist.dist)
		foundA := false
		foundB := false
		for i, grp := range groups {
			if 0 <= slices.IndexFunc(grp, func(c int) bool { return c == nextDist.sId }) {
				log.Printf("  2 found face %d in group %d\n", nextDist.sId, i)
				var found bool
				found, groups, dists = mergeGroupsIfPossible(groups, dists, nextDist.tId, i)
				if found {
					foundA = true
					break
				}
				log.Printf("  3 found new face %d for group %d\n", nextDist.tId, i)
				dists = deleteUnused(grp, nextDist.tId, dists)
				groups[i] = append(grp, nextDist.tId)
				log.Println("  3a group", i, groups[i])
				foundA = true
				break
			}
			if 0 <= slices.IndexFunc(grp, func(c int) bool { return c == nextDist.tId }) {
				log.Printf("  4 found new face %d for group %d\n", nextDist.tId, i)
				var found bool
				found, groups, dists = mergeGroupsIfPossible(groups, dists, nextDist.sId, i)
				if found {
					foundB = true
					break
				}
				log.Printf("  5 found new face %d for group %d\n", nextDist.sId, i)
				dists = deleteUnused(grp, nextDist.sId, dists)
				groups[i] = append(grp, nextDist.sId)
				log.Println("  5a group", i, groups[i])
				foundB = true
				break
			}
		}
		if foundA || foundB {
			continue
		}
		log.Printf("  6 have to create new group for sId %d and tId %d\n", nextDist.sId, nextDist.tId)
		groups = append(groups, []int{nextDist.sId, nextDist.tId})
		log.Println("  6a group", len(groups)-1, groups[len(groups)-1])
		elemCnt := 0
		for _, grp := range groups {
			elemCnt += len(grp)
		}
		log.Println("  6b elements grouped", elemCnt, "of", len(descr))
		if elemCnt == len(descr) {
			log.Printf("ALL faces grouped\n")
			break
		}
	}
	log.Println("FINISHED groupCnt:", len(groups), "distsCnt", len(dists))
	// Now sort subgroups
	var groupsSorted [][]int
	for len(dists) > 0 {
		log.Printf("10 len(dists) %d len(groups) %d\n", len(dists), len(groups))
		nextDist := dists[0]
		dists = dists[1:]
		log.Printf(" 11 next distance between %d and %d is %f\n", nextDist.sId, nextDist.tId, nextDist.dist)
		foundS := false
		foundD := false
		for i, grp := range groups {
			if !foundS && 0 <= slices.IndexFunc(grp, func(c int) bool { return c == nextDist.sId }) {
				log.Printf("  12 found face %d in group %d\n", nextDist.sId, i)
				groupsSorted = append(groupsSorted, groups[i])
				groups = slices.Delete(groups, i, i+1)
				foundS = true
			}
			if !foundD && 0 <= slices.IndexFunc(grp, func(c int) bool { return c == nextDist.tId }) {
				log.Printf("  14 found new face %d for group %d\n", nextDist.sId, i)
				groupsSorted = append(groupsSorted, groups[i])
				groups = slices.Delete(groups, i, i+1)
				foundD = true
			}
		}
	}
	// order remaing group by lng descending
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})
	// Add them in this order
	for _, grp := range groups {
		groupsSorted = append(groupsSorted, grp)
	}
	log.Println("FINISHED sorting groups sortedCnt:", len(groupsSorted), "groupCnt:", len(groups))
	//
	j := -len(groups)
	for i, grp := range groupsSorted {
		for _, faceImageId := range grp {
			if err := setImageFaceSubGroup(db, faceImageId, j); err != nil {
				log.Println(" Saving Subgroup", j, "failed for FaceImageId", faceImageId, "err", err)
			} else {
				log.Println("Saved Subgroup", j, "succeeded for FaceImageId", faceImageId)
			}
		}
		log.Println(" Group", j, i, "len", len(grp), grp)
		j += 1
	}
}

func euclidean_distance(src face.Descriptor, dst face.Descriptor) (r float64) {

	sum_of_squares := 0.0
	for i, coord := range src {
		//sum_of_squares += math.Pow(float64(coord-dst[i]), 2)
		diff := float64(coord - dst[i])
		sum_of_squares += diff * diff
	}
	//return 1 / (1 + math.Sqrt(sum_of_squares))
	return math.Sqrt(sum_of_squares)
	//return 0.5 * sum_of_squares
}
