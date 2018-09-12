package images

import (
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"../db"
	"github.com/nfnt/resize"
	"gopkg.in/mgo.v2/bson"
)

// WidthPreview width of image preview
const WidthPreview = 450

// PathToImg path to derectory with pichers
const PathToImg = "./src/img"

// RequestCreateAlbum struct request for create new album
type RequestCreateAlbum struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

// Album user album
type Album struct {
	Name         string        `json:"name" bson:"name"`
	TimeOfCreate int64         `json:"timeOfCreate" bson:"timeOfCreate"`
	Description  string        `json:"description" bson:"description"`
	UserID       bson.ObjectId `json:"userId" bson:"userId"`
}

// Img new picher
type Img struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Name       string        `json:"name" bson:"name"`
	Time       int64         `json:"time" bson:"time"`
	AlbumID    bson.ObjectId `json:"albumId" bson:"albumId"`
	URL        string        `json:"url" bson:"url"`
	URLPreview string        `json:"urlPreview" bson:"urlPreview"`
}

// ImgInfo is being created after DecodeConfig
type ImgInfo struct {
	Width  int
	Height int
	Format string
}

// AlbumID this id of album
type AlbumID struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

// CreateAlbum create new album
func CreateAlbum(w http.ResponseWriter, r *http.Request) {
	id, err := db.ValidateToken(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var target RequestCreateAlbum
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newAlbum := Album{
		Name:         target.Name,
		TimeOfCreate: time.Now().Unix(),
		Description:  target.Description,
		UserID:       bson.ObjectId(id),
	}

	var s db.UseDb = &db.Session{}
	defer s.CloseSession()
	s.GetCollection("albums").Insert(&newAlbum)
	json.NewEncoder(w).Encode(&newAlbum)
}

// LoadImage load image and add to db
func LoadImage(w http.ResponseWriter, r *http.Request) {
	id, err := db.ValidateToken(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseMultipartForm(1024)
	var s db.UseDb = &db.Session{}
	err = s.CreateSession()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer s.CloseSession()
	if r.FormValue("albumId") != "" {
		findAlbum, err := foundAlbum(id, r.FormValue("albumId"), s)
		if err != nil {
			http.Error(w, "album is not found", http.StatusNoContent)
			return
		}
		workWithPicher(w, r, findAlbum, s)
	}
	defaultAlbumID, errDefault := foundDefaultAlbum(id, s)
	if errDefault != nil {
		createAlbumID, errD := createDefaultAlbum(id, s)
		if errD != nil {
			http.Error(w, errD.Error(), http.StatusBadRequest)
			return
		}
		defaultAlbumID = createAlbumID
	}
	workWithPicher(w, r, defaultAlbumID, s)
}

func workWithPicher(w http.ResponseWriter, r *http.Request, album AlbumID, s db.UseDb) {
	img, err := setImg(r, album)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = addPicherToBD(&img, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(&img)
}

func addPicherToBD(img *Img, s db.UseDb) error {
	err := s.GetCollection("images").Insert(&img)
	if err != nil {
		return err
	}
	return nil
}

func setImg(r *http.Request, id AlbumID) (Img, error) {
	var newImg Img
	fileName, err := pushImg(r)
	if err != nil {
		return newImg, err
	}
	newImg = Img{Name: fileName, URL: fileName, Time: time.Now().Unix(), AlbumID: id.ID}
	infoPicher, err := getInfoImg(PathToImg + fileName)
	if err != nil {
		return newImg, err
	}
	if infoPicher.Width < 400 {
		newImg.ID = bson.NewObjectId()
		newImg.URLPreview = fileName
		return newImg, nil
	}
	preview, err := createPreview(fileName, infoPicher.Format)
	if err != nil {
		return newImg, err
	}
	newImg.ID = bson.NewObjectId()
	newImg.URLPreview = preview
	return newImg, nil
}

func pushImg(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("img")
	if err != nil {
		return "", err
	}
	if handler.Header["Content-Type"][0] != "image/jpeg" && handler.Header["Content-Type"][0] != "image/png" {
		return "", errors.New("unvalidate format")
	}
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + handler.Filename
	fCreate, err := os.Create(PathToImg + fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	defer fCreate.Close()
	io.Copy(fCreate, file)
	return fileName, nil
}

func createPreview(fileName string, format string) (string, error) {
	openImg, err := os.Open(PathToImg + fileName)
	if err != nil {
		return "", err
	}
	defer openImg.Close()
	var img image.Image
	if format == "jpeg" {
		img, err = jpeg.Decode(openImg)
	} else {
		img, err = png.Decode(openImg)
	}
	if err != nil {
		return "", err
	}
	m := resize.Resize(WidthPreview, 0, img, resize.Lanczos3)
	preview, err := os.Create("./src/img/" + "preview_" + fileName)
	defer preview.Close()
	if format == "jpeg" {
		jpeg.Encode(preview, m, nil)
	} else {
		png.Encode(preview, m)
	}
	return "preview_" + fileName, nil
}

func getInfoImg(fileName string) (ImgInfo, error) {
	openImg, err := os.Open(fileName)
	if err != nil {
		return ImgInfo{}, err
	}
	defer openImg.Close()
	infoPicher, format, err := image.DecodeConfig(openImg)
	if err != nil {
		return ImgInfo{}, err
	}
	return ImgInfo{Width: infoPicher.Width, Height: infoPicher.Height, Format: format}, nil
}

func foundAlbum(id string, idA string, s db.UseDb) (AlbumID, error) {
	var getAlbum AlbumID
	err := s.GetCollection("albums").Find(bson.M{"_id": bson.ObjectId(idA), "userId": bson.ObjectId(id)}).One(&getAlbum)
	if err != nil {
		return getAlbum, err
	}
	return getAlbum, nil
}

func foundDefaultAlbum(id string, s db.UseDb) (AlbumID, error) {
	var getAlbum AlbumID
	err := s.GetCollection("albums").Find(bson.M{"default": true, "userId": bson.ObjectId(id)}).One(&getAlbum)
	if err != nil {
		return getAlbum, err
	}
	return getAlbum, nil
}

func createDefaultAlbum(id string, s db.UseDb) (AlbumID, error) {
	i := bson.NewObjectId()
	err := s.GetCollection("albums").Insert(bson.M{"_id": i, "userId": bson.ObjectId(id), "name": "Default album", "description": "", "timeToCreate": time.Now().Unix(), "default": true})
	if err != nil {
		return AlbumID{}, errors.New("Failed created default album")
	}
	return AlbumID{ID: i}, nil
}
