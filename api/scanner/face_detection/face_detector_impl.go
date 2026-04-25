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
	"gorm.io/gorm"
	"gopkg.in/gographics/imagick.v3/imagick"
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

	return
}

func setImageFaceSubGroup(db *gorm.DB, imageFaceID int, subGroup int) (err error) {

	if err := db.Where("id = ?", imageFaceID).Model(&models.ImageFace{}).UpdateColumn("subgroup", subGroup).Error; err != nil {
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

	log.Println("  ", len(media.Faces)," Faces exist already for: ", media.ID)
	return nil
}

// DetectFaces finds the faces in the given image and saves them to the database
func (fd *faceDetector) DetectFaces(db *gorm.DB, media *models.Media) error {
	if err := db.Model(media).Preload("MediaURL").First(&media).Error; err != nil {
		return err
	}

	log.Println("  ", len(media.Faces)," Faces exist already for: ", media.ID)

	var faces []face.Face
	var thumbnailPath string
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
			if err = mw.AutoOrientImage(); err != nil {
				log.Printf("Err: AutoOrientImage %s\n", err)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			if err = mw.SetImageFormat("RGB"); err != nil {
				log.Printf("Err: SetImageFormat RGB %s\n", err)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			depth := mw.GetImageDepth()
			color := mw.GetImageColorspace()
			if color == 3 {
				log.Printf("   GRAY ColorSpace %d %s\n", color, cachedPath)
				if err = mw.AddImage(mw.Clone()); err != nil {
					log.Printf("Err: AddImage 1 %s\n", err)
					return errors.Wrap(err, "error read faces")
					continue
				}
				if err = mw.AddImage(mw.Clone()); err != nil {
					log.Printf("Err: AddImage 2 %s\n", err)
					return errors.Wrap(err, "error read faces")
					continue
				}
				if err = mw.AddImage(mw.Clone()); err != nil {
					log.Printf("Err: AddImage 3 %s\n", err)
					return errors.Wrap(err, "error read faces")
					continue
				}
				mw = mw.CombineImages(23);
				color = mw.GetImageColorspace()
			}
			if color != 23 {
				log.Printf("   ColorSpace %d  Depth %d\n", color, depth)
				// if err = mw.SetImageColorspace(COLORSPACE_SRGB); err != nil {
				if err = mw.SetImageColorspace(23); err != nil {
					log.Printf("Err: SetImageColorspace %s\n", err)
					return errors.Wrap(err, "error read faces")
					continue
				}
				// depthn := mw.GetImageDepth()
				// colorn := mw.GetImageColorspace()
				// log.Printf("   New ColorSpace %d  Depth %d\n", colorn, depthn)
			}
			if err = mw.SetImageFormat("JPEG"); err != nil {
				log.Printf("Err: SetImageFormat %s\n", err)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			var blob []byte
			if blob, err = mw.GetImageBlob(); err != nil {
				log.Printf("Err: GetImageBlob %s\n", err)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			fd.mutex.Lock()
			if faces, err = fd.rec.Recognize(blob); err != nil {
				fd.mutex.Unlock()
				log.Printf("Err: Recognize %s\n", err)
				// return errors.Wrap(err, "error read faces")
				continue
			}
			fd.mutex.Unlock()

			thumbnailPath = cachedPath
			break
		}

		if url.Purpose == models.PhotoThumbnail {
			log.Println("  Thumbnail URL found %s", cachedPath)

			fd.mutex.Lock()
			log.Printf("RecognizeFile %s\n", cachedPath)
			faces, err = fd.rec.RecognizeFile(cachedPath)
			fd.mutex.Unlock()

			if err != nil {
				return errors.Wrap(err, "error read faces")
			}
			thumbnailPath = cachedPath
			continue
		}

		// if url.Purpose == models.PhotoHighRes {
			// thumbnailURL = &url
			// thumbnailURL.Media = media
			// log.Println("  Thumbnail URL found")
			// break
		// }
	}

	log.Println("  ", len(faces)," Faces found for: ", thumbnailPath)
	for _, face := range faces {
		// log.Println("  Faces found for: ", face)
		fd.classifyFace(db, &face, media, thumbnailPath)
	}

	return nil
}

func (fd *faceDetector) classifyFace(db *gorm.DB, face *face.Face, media *models.Media, imagePath string) error {
	fd.mutex.Lock()
	defer fd.mutex.Unlock()

	match := int32(fd.rec.ClassifyThreshold(face.Descriptor, fd.classifyThreshold))

	dimension, err := media_encoding.GetPhotoDimensions(imagePath)
	if err != nil {
		return err
	}

	imageFace := models.ImageFace{
		MediaID:    media.ID,
		Descriptor: models.FaceDescriptor(face.Descriptor),
		Rectangle: models.FaceRectangle{
			// Converts a pixel absolute rectangle to a relative FaceRectangle.
			MinX: float64(face.Rectangle.Min.X) / float64(dimension.Width),
			MaxX: float64(face.Rectangle.Max.X) / float64(dimension.Width),
			MinY: float64(face.Rectangle.Min.Y) / float64(dimension.Height),
			MaxY: float64(face.Rectangle.Max.Y) / float64(dimension.Height),
		},
	}
	// log.Printf("    Face region: %f-%f:%f-%f", imageFace.Rectangle.MinX, imageFace.Rectangle.MaxX, imageFace.Rectangle.MinY, imageFace.Rectangle.MaxX)

	var faceGroup models.FaceGroup

	// If no match add it new to samples
	if match < 0 {
		// log.Println("     No match, assigning new face")

		faceGroup = models.FaceGroup{
			ImageFaces: []models.ImageFace{imageFace},
		}

		if err := db.Create(&faceGroup).Error; err != nil {
			return err
		}

	} else {
		// log.Println("     Found match")

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

		for _, id := range imageFaceIDs {
			if imageFaceID == id {
				fd.faceGroupIDs[i] = destFaceGroupID
				break
			}
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
	}

	return updatedImageFaces, nil
}

func (fd *faceDetector) findFaceGroupCandidates( faceGroup int, checkLng int) {

	fd.mutex.Lock()
	defer fd.mutex.Unlock()

	k := len(fd.faceDescriptors) - checkLng
	log.Printf("findFaceGroupCandidates Id %d Lng %d start %d chkLng %d\n", faceGroup, len(fd.faceDescriptors), k, checkLng)

	descriptor := fd.faceDescriptors[k]
	faceGroupID := fd.faceGroupIDs[k]
	imageFaceID := fd.imageFaceIDs[k]
	fd.faceDescriptors = slices.Delete(fd.faceDescriptors, k, k)
	fd.faceGroupIDs = slices.Delete(fd.faceGroupIDs, k, k)
	fd.imageFaceIDs = slices.Delete(fd.imageFaceIDs, k, k)
	for j := range checkLng {
		i := j + k
		// log.Printf("findFaceGroupCandidates index %d\n", i)
		fd.rec.SetSamples(fd.faceDescriptors, fd.faceGroupIDs)
		match := fd.rec.ClassifyThreshold( descriptor, fd.classifyThreshold)
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

func (fd *faceDetector) ConfirmedFaceGroups(db *gorm.DB) ( []int, error ) {
	var groupIDs []int
	if err := db.
		Model(models.ImageFace{}).
		// Debug().
		Where("confirmed = TRUE").
		Group("face_group_id").
		Order("COUNT(face_group_id)").
		Pluck("face_group_id",&groupIDs).
		Error; err != nil {

		return make([]int,0), err
	}

	// for _, groupId := range groupIDs {
		// log.Printf("Group ID %d with confirmed faces\n", groupId)
	// }
	return groupIDs, nil
}

func (fd *faceDetector) FaceGroupFaceIDs(db *gorm.DB, confirmed bool, groupId int) ( []int, error ) {

	var faceIDs []int
	if err := db.
		Model(models.ImageFace{}).
		// Debug().
		Where("face_group_id = ?", groupId).
		Where("confirmed = ?", confirmed).
		Pluck("ID",&faceIDs).
		Error; err != nil {

		return make([]int,0), err
	}

	var faceGroup models.FaceGroup
	if err := db.
		// Debug().
		Model(models.FaceGroup{}).
		Where("id = ?", groupId).
		First(&faceGroup).
		Error; err != nil {

		return make([]int,0), err
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
		confFaceIDs, _ := fd.FaceGroupFaceIDs( db, true, grpID)
		log.Printf("Group %d with %d confirmed faces\n", grpID, len(confFaceIDs))
		// for _, faceID := range confFaceIDs {
			// log.Printf("GrpID %d faceID %d \n", grpID, faceID)
		// }
		faceIDs, _ := fd.FaceGroupFaceIDs( db, false, grpID)
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
		fd.faceDescriptors = slices.Delete(fd.faceDescriptors, idx, idx)
		fd.faceGroupIDs = append(fd.faceGroupIDs, 0)
		fd.faceGroupIDs = slices.Delete(fd.faceGroupIDs, idx, idx)
		fd.imageFaceIDs = append(fd.imageFaceIDs, faceID)
		fd.imageFaceIDs = slices.Delete(fd.imageFaceIDs, idx, idx)
	}
	fd.findFaceGroupCandidates( groupID, len(faceIDs))
	fd.faceDescriptors = saveFaceDescriptors
	fd.faceGroupIDs = saveFaceGroupIDs
	fd.imageFaceIDs = saveImageFaceIDs
}

type distElem struct {
	sId	int
	tId	int
	dist	float64
}

func deleteUnused(grp []int, newId int, dists []distElem) (r []distElem) {
	for _, grpId := range grp {
		dists = slices.DeleteFunc(dists, func(c distElem) bool {
			return (c.sId == grpId) && (c.tId == newId) || (c.tId == grpId) && (c.sId == newId)
		})
	}
	return dists
}

func mergeGroups(groups [][]int, dists []distElem, tId int, i int) (r bool, rgroups [][]int, rdists []distElem) {
	for j, grpS := range groups {
		if 0 <= slices.IndexFunc(grpS, func(c int) bool { return c == tId }) {
			log.Printf("   found two related groups that must be merged %d %d\n", i, j)
			// merge groups because they are close to each other
			for _, faceId := range grpS {
				dists = deleteUnused(groups[i], faceId, dists)
			}
			groups[i] = slices.Concat(groups[i], grpS)
			log.Println("   7a group", i, groups[i])
			groups = slices.Delete(groups, j, j)
			log.Println("   7b groupCnt", len(groups))
			return true, groups, dists
		}
	}
	log.Printf("   No group merge to group %d\n", i)
	return false, groups, dists
}

func (fd *faceDetector) SplitFaceGroup(db *gorm.DB, groupID int32) {

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
		for i, grp := range groups {
			if 0 <= slices.IndexFunc(grp, func(c int) bool { return c == nextDist.sId }) {
				log.Printf("  2 found face %d in group %d\n", nextDist.sId, i)
				var found bool
				found, groups, dists = mergeGroups(groups, dists, nextDist.tId, i)
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
				log.Printf("  4 found new face %d for group %d\n", nextDist.sId, i)
				dists = deleteUnused(grp, nextDist.sId, dists)
				groups[i] = append(grp, nextDist.sId)
				log.Println("  4a group", i, groups[i])
				foundA = true
				break
			}
		}
		if foundA {
			continue
		}
		log.Printf("  5 have to create new group for sId %d and tId %d\n", nextDist.sId, nextDist.tId)
		groups = append(groups, []int {nextDist.sId, nextDist.tId})
		log.Println("  5a group", len(groups)-1, groups[len(groups)-1])
		elemCnt := 0
		for _, grp := range groups {
			elemCnt += len(grp)
		}
		log.Println("  5b elements grouped", elemCnt, "of", len(descr))
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
				groups = slices.Delete(groups, i, i)
				foundS = true
			}
			if !foundD && 0 <= slices.IndexFunc(grp, func(c int) bool { return c == nextDist.tId }) {
				log.Printf("  14 found new face %d for group %d\n", nextDist.sId, i)
				groupsSorted = append(groupsSorted, groups[i])
				groups = slices.Delete(groups, i, i)
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
	for i, grp := range groupsSorted {
		for _, faceImageId := range grp {
			if err := setImageFaceSubGroup(db, faceImageId, i); err != nil {
				log.Println(" Saving Subgroup", i, "failed for FaceImageId", faceImageId, "err", err)
			} else {
				log.Println("Saved Subgroup", i, "succeeded for FaceImageId", faceImageId)
			}
		}
		log.Println(" Group", i, "len", len(grp), grp)
	}
}

func euclidean_distance(src face.Descriptor, dst face.Descriptor) (r float64) {

	sum_of_squares := 0.0
	for i, coord := range src {
		//sum_of_squares += math.Pow(float64(coord-dst[i]), 2)
		diff := float64(coord-dst[i])
		sum_of_squares += diff * diff
	}
	//return 1 / (1 + math.Sqrt(sum_of_squares))
	return math.Sqrt(sum_of_squares)
	//return 0.5 * sum_of_squares
}
