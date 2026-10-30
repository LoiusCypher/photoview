package actions

import (
	"strings"

	"github.com/loiuscypher/photoview/api/database/drivers"
	"github.com/loiuscypher/photoview/api/graphql/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Search(db *gorm.DB, query string, userID int, limitMedia *int, limitAlbums *int) (*models.SearchResult, error) {
	limitMediaInternal := 10
	limitAlbumsInternal := 10
	limitFacesInternal := 10

	if limitMedia != nil {
		limitMediaInternal = *limitMedia
	}

	if limitAlbums != nil {
		limitAlbumsInternal = *limitAlbums
	}

	wildQuery := "%" + strings.ToLower(query) + "%"

	var media []*models.Media

	userSubquery := db.Table("user_albums").Where("user_id = ?", userID)
	if drivers.POSTGRES.MatchDatabase(db) {
		userSubquery = userSubquery.Where("album_id = \"Album\".id")
	} else {
		userSubquery = userSubquery.Where("album_id = Album.id")
	}

	err := db.Joins("Album").
		Where("EXISTS (?)", userSubquery).
		Where("LOWER(media.title) LIKE ? OR LOWER(media.path) LIKE ?", wildQuery, wildQuery).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{
				SQL:    "(CASE WHEN LOWER(media.title) LIKE ? THEN 2 WHEN LOWER(media.path) LIKE ? THEN 1 END) DESC",
				Vars:               []interface{}{wildQuery, wildQuery},
				WithoutParentheses: true},
		}).
		Limit(limitMediaInternal).Find(&media).Error

	if err != nil {
		return nil, errors.Wrapf(err, "searching media")
	}

	var albums []*models.Album

	err = db.
		Where("EXISTS (?)", db.Table("user_albums").Where("user_id = ?", userID).Where("album_id = albums.id")).
		Where("albums.title LIKE ? OR albums.path LIKE ?", wildQuery, wildQuery).
		Clauses(clause.OrderBy{
			Expression: clause.Expr{
				SQL:                "(CASE WHEN albums.title LIKE ? THEN 2 WHEN albums.path LIKE ? THEN 1 END) DESC",
				Vars:               []interface{}{wildQuery, wildQuery},
				WithoutParentheses: true},
		}).
		Limit(limitAlbumsInternal).
		Find(&albums).Error

	if err != nil {
		return nil, errors.Wrapf(err, "searching albums")
	}

//////
	var faceGroups []*models.FaceGroup
 
	err = db.
                Joins("JOIN image_faces ON image_faces.face_group_id = face_groups.id").
                Joins("JOIN media ON media.id = image_faces.face_group_id").
                Joins("JOIN user_albums ON user_albums.album_id = media.album_id").
                Group("face_groups.id").
                Where("face_groups.label IS NOT NULL").
		Where("user_albums.user_id = ?", userID).
		Where("face_groups.label LIKE ?", wildQuery).
                Order("face_groups.label ASC").
		Limit(limitFacesInternal).
		Find(&faceGroups).Error

	if err != nil {
		return nil, errors.Wrapf(err, "searching faceGroups")
	}

	result := models.SearchResult{
		Query:  query,
		Media:  media,
		Albums: albums,
		FaceGroups: faceGroups,
	}

	return &result, nil
}
