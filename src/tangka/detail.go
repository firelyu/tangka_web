package tangka

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	uploadFileSuffix    = ".tka"
	uploadFilePermisson = 0600
)

const (
	UploadFolder        = "upload"
	idLength            = 8
)

type Tangka struct {
	Id     string // The [0-9a-z]*, the length is 8. It is unique.
	Name   string
	Author string
	//CreateDate time.Time
	//General string
	//Detail string
	//Price int
	//ImageUrl string
}

var idSeed = rand.NewSource(time.Now().UnixNano())

const idLetterBytes = "0123456789abcdefghijklmnopqrstuvwxyz"
const (
	idLetterIdxBits = 6                    // The length of letterBytes is 32, 6 bits to represent a letter index
	idLetterIdxMask = 077                  // The index range is 0-35, 033 is octal
	idLetterIdxMax  = 63 / idLetterIdxBits // # of letter indices fitting in 63 bits
)

func randStringBytesMask(n int) string {
	b := make([]byte, n)
	// A idSeed.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, idSeed.Int63(), idLetterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = idSeed.Int63(), idLetterIdxMax
		}
		if idx := int(cache & idLetterIdxMask); idx < len(idLetterBytes) {
			b[i] = idLetterBytes[idx]
			i--
		}
		cache >>= idLetterIdxBits
		remain--
	}

	return string(b)
}

func generaTangkaId() (string, error) {
	// TODO: check the id exist.
	return randStringBytesMask(idLength), nil
}

func NewTangka(name string, author string) *Tangka {
	id, err := generaTangkaId()
	if err != nil {
		return nil
	}

	return &Tangka{
		Id:     id,
		Name:   name,
		Author: author,
	}
}

func (t *Tangka) Save() error {
	if t.Id == "" {
		return errors.New("The thangka Id is NULL.")
	}
	fileName := t.Id + uploadFileSuffix
	permission := os.FileMode(uploadFilePermisson)

	return ioutil.WriteFile(UploadFolder + "/" + fileName, []byte(t.Name), permission)
}

func (t *Tangka) Delete() error {
	fileName := t.Id + uploadFileSuffix
	return os.Remove(UploadFolder + "/" + fileName)
}

func GetTangkaById(id string) (t *Tangka, err error) {
	fileName := UploadFolder + "/" + id + uploadFileSuffix
	body, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return &Tangka{Id: id, Name: string(body)}, nil
}

type TangkaList struct {
	List []*Tangka
}

func ListAllTangka() (tList *TangkaList, err error) {
	fileList, err := ioutil.ReadDir(UploadFolder)
	if err != nil {
		return nil, err
	}

	var tangkaList []*Tangka

	for _, file := range fileList {
		id := strings.TrimSuffix(file.Name(), uploadFileSuffix)
		name, err := ioutil.ReadFile(UploadFolder + "/" + file.Name())
		if err != nil {
			return nil, err
		}
		tangkaList = append(tangkaList, &Tangka{Id: id, Name: string(name)})
	}

	return &TangkaList{List: tangkaList}, nil
}
