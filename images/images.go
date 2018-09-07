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
		workWithImg(w, r, findAlbum)
		return
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
	workWithImg(w, r, defaultAlbumId)
}

func workWithImg(w http.ResponseWriter, r *http.Request, albumId AlbumId) {
	var newImg Img
	file, handler, err := r.FormFile("img")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	if handler.Header["Content-Type"][0] != "image/jpeg" && handler.Header["Content-Type"][0] != "image/png" {
		http.Error(w, "unvalidate format", http.StatusBadRequest)
		return
	}
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + handler.Filename
	fCreate, err := os.Create("./src/img/" + fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newImg = Img{Name: handler.Filename, Url: fileName, Time: time.Now().Unix(), AlbumId: albumId.ID}
	defer fCreate.Close()
	io.Copy(fCreate, file)
	infoPicher, err := getInfoImg(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if infoPicher.Width < 400 {
		newImg.UrlPreview = fileName
		json.NewEncoder(w).Encode(&newImg)
		return
	}
	var preview string
	if infoPicher.Format == "jpeg" {
		previewName, err := createPreviewJpeg(fileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		preview = previewName
	}
	if infoPicher.Format == "png" {
		previewName, err := createPreviewPng(fileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		preview = previewName
	}
	newImg.UrlPreview = preview
	json.NewEncoder(w).Encode(&newImg)
}

func createPreviewJpeg(fileName string) (string, error) {
	openImg, err := os.Open("./src/img/" + fileName)
	if err != nil {
		return "", err
	}
	defer openImg.Close()
	img, err := jpeg.Decode(openImg)
	if err != nil {
		return "", err
	}
	m := resize.Resize(WIDTH_PREVIEW, 0, img, resize.Lanczos3)
	preview, err := os.Create("./src/img/" + "preview_" + fileName)
	defer preview.Close()
	jpeg.Encode(preview, m, nil)
	return "preview_" + fileName, nil
}

func createPreviewPng(fileName string) (string, error) {
	openImg, err := os.Open("./src/img/" + fileName)
	if err != nil {
		return "", err
	}
	defer openImg.Close()
	img, err := png.Decode(openImg)
	if err != nil {
		return "", err
	}
	m := resize.Resize(WIDTH_PREVIEW, 0, img, resize.Lanczos3)
	preview, err := os.Create("./src/img/" + "preview_" + fileName)
	defer preview.Close()
	png.Encode(preview, m)
	return "preview_" + fileName, nil
}

func getInfoImg(fileName string) (ImgInfo, error) {
	openImg, err := os.Open("./src/img/" + fileName)
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
