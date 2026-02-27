package scanner_tasks

import (
	"io/fs"
	"log"

	"github.com/loiuscypher/photoview/api/scanner/scanner_task"
	ignore "github.com/sabhiram/go-gitignore"
)

type IgnorefileTask struct {
	scanner_task.ScannerTaskBase
}

type ignorefileTaskKey string

const albumIgnoreKey ignorefileTaskKey = "album_ignore_key"

func getAlbumIgnore(ctx scanner_task.TaskContext) *ignore.GitIgnore {
	return ctx.Value(albumIgnoreKey).(*ignore.GitIgnore)
}

func (t IgnorefileTask) BeforeScanAlbum(ctx scanner_task.TaskContext) (scanner_task.TaskContext, error) {
	a := ctx.GetAlbum().Path
	log.Printf("Path %s \n", a)
	c := ctx.GetCache()
	log.Printf("Cache %d \n", c)
	d := c.IsPathMedia("/photos/album5")
	log.Printf("Is media %d \n", d)
	e := ctx.Value("album_ignore_key")
	log.Printf("Value %s \n", e)
	b := *c.GetAlbumIgnore(a)
	log.Printf("Ignore %s \n", b)
	albumIgnore := ignore.CompileIgnoreLines(*ctx.GetCache().GetAlbumIgnore(ctx.GetAlbum().Path)...)
	return ctx.WithValue(albumIgnoreKey, albumIgnore), nil
}

func (t IgnorefileTask) MediaFound(ctx scanner_task.TaskContext, fileInfo fs.FileInfo, mediaPath string) (bool, error) {

	// Match file against ignore data
	if getAlbumIgnore(ctx).MatchesPath(fileInfo.Name()) {
		log.Printf("File %s ignored\n", fileInfo.Name())
		return true, nil
	}

	return false, nil
}
