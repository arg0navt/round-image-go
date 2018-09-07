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

const WIDTH_PREVIEW = 450
const PATH_TO_IMG = "./src/img/"

type RequestCreateAlbum struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

type Album struct {
	Name         string        `json:"name" bson:"name"`
	TimeOfCreate int64         `json:"timeOfCreate" bson:"timeOfCreate"`
	Description  string        `json:"description" bson:"description"`
	UserID       bson.ObjectId `json:"userId" bson:"userId"`
}

type Img struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	Name       string        `json:"name" bson:"name"`
	Time       int64         `json:"time" bson:"time"`
	AlbumId    bson.ObjectId `json:"albumId" bson:"albumId"`
	Url        string        `json:"url" bson:"url"`
	UrlPreview string        `json:"urlPreview" bson:"urlPreview"`
}

type ImgInfo struct {
	Width  int
	Height int
	Format string
}

type AlbumId struct {
	ID bson.ObjectId `json:"id" bson:"_id"`
}

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
	db.GetCollection("albums").Insert(&newAlbum)
	json.NewEncoder(w).Encode(&newAlbum)
}

func LoadImage(w http.ResponseWriter, r *http.Request) {
	id, err := db.ValidateToken(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	r.ParseMultipartForm(1024)
	if r.FormValue("albumId") != "" {
		findAlbum, err := foundAlbum(id, r.FormValue("albumId"))
		if err != nil {
			http.Error(w, "album is not found", http.StatusNoContent)
			return
		}
		workWithPicher(w, r, findAlbum)
	}
	defaultAlbumId, errDefault := foundDefaultAlbum(id)
	if errDefault != nil {
		createAlbumId, errD := createDefaultAlbum(id)
		if errD != nil {
			http.Error(w, errD.Error(), http.StatusBadRequest)
			return
		}
		defaultAlbumId = createAlbumId
	}
	workWithPicher(w, r, defaultAlbumId)
}

func workWithPicher(w http.ResponseWriter, r *http.Request, album AlbumId) {
	img, err := setImg(r, album)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = addPicherToBD(&img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(&img)
}

func addPicherToBD(img *Img) error {
	err := db.GetCollection("images").Insert(&img)
	if err != nil {
		return err
	}
	return nil
}

func setImg(r *http.Request, albumId AlbumId) (Img, error) {
	var newImg Img
	fileName, err := pushImg(r)
	if err != nil {
		return newImg, err
	}
	newImg = Img{Name: fileName, Url: fileName, Time: time.Now().Unix(), AlbumId: albumId.ID}
	infoPicher, err := getInfoImg(PATH_TO_IMG + fileName)
	if err != nil {
		return newImg, err
	}
	if infoPicher.Width < 400 {
		newImg.Id = bson.NewObjectId()
		newImg.UrlPreview = fileName
		return newImg, nil
	}
	preview, err := createPreview(fileName, infoPicher.Format)
	if err != nil {
		return newImg, err
	}
	newImg.Id = bson.NewObjectId()
	newImg.UrlPreview = preview
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
	fCreate, err := os.Create(PATH_TO_IMG + fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	defer fCreate.Close()
	io.Copy(fCreate, file)
	return fileName, nil
}

func createPreview(fileName string, format string) (string, error) {
	openImg, err := os.Open(PATH_TO_IMG + fileName)
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
	m := resize.Resize(WIDTH_PREVIEW, 0, img, resize.Lanczos3)
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

func foundAlbum(id string, albumId string) (AlbumId, error) {
	var getAlbum AlbumId
	err := db.GetCollection("albums").Find(bson.M{"_id": bson.ObjectId(albumId), "userId": bson.ObjectId(id)}).One(&getAlbum)
	if err != nil {
		return getAlbum, err
	}
	return getAlbum, nil
}

func foundDefaultAlbum(id string) (AlbumId, error) {
	var getAlbum AlbumId
	err := db.GetCollection("albums").Find(bson.M{"default": true, "userId": bson.ObjectId(id)}).One(&getAlbum)
	if err != nil {
		return getAlbum, err
	}
	return getAlbum, nil
}

func createDefaultAlbum(id string) (AlbumId, error) {
	i := bson.NewObjectId()
	err := db.GetCollection("albums").Insert(bson.M{"_id": i, "userId": bson.ObjectId(id), "name": "Default album", "description": "", "timeToCreate": time.Now().Unix(), "default": true})
	if err != nil {
		return AlbumId{}, errors.New("Failed created default album")
	}
	return AlbumId{ID: i}, nil
}
