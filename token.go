package appleapi

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

var (
	ErrAuthKeyNotPem   = errors.New("token: AuthKey must be a valid .p8 PEM file")
	ErrAuthKeyNotECDSA = errors.New("token: AuthKey must be of type ecdsa.PrivateKey")
	ErrAuthKeyNil      = errors.New("token: AuthKey was nil")
)

const PlatformIos = "IOS"
const PlatformMac = "MAC_OS"

func ReadPrivate(bytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, ErrAuthKeyNotPem
	}

	//log.Printf("blockKey:%v", block.Bytes)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		key, err = x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}
	switch pk := key.(type) {
	case *ecdsa.PrivateKey:
		return pk, nil
	default:
		return nil, ErrAuthKeyNotECDSA
	}
}

// AuthKeyFromFile loads a .p8 certificate from a local file and returns a
// *ecdsa.PrivateKey.
func AuthKeyFromFile(fileName string) (*ecdsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return ReadPrivate(bytes)
}

type ErrorResponse struct {
	Error []struct {
		Code   string            `json:"code,omitempty"`
		Status string            `json:"status,omitempty"`
		Id     string            `json:"id,omitempty"`
		Title  string            `json:"title,omitempty"`
		Detail string            `json:"detail,omitempty"`
		Source map[string]string `json:"source,omitempty"`
	} `json:"errors,omitempty"`
}

type PagingInformation struct {
	Paging struct {
		Total int `json:"total,omitempty"` //The total number of resources matching your request.
		Limit int `json:"limit,omitempty"` //The maximum number of resources to return per page, from 0 to 200.
	} `json:"paging,omitempty"`
}

type PagedDocumentLinks struct {
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Self  string `json:"self,omitempty"`
}

type TypeId struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type DocumentLinks struct {
	Self string `json:"self,omitempty"`
}

type ResourceLinks struct {
	Self string `json:"self,omitempty"`
}

type Token struct {
	Secret  string
	Kid     string
	Iss     string
	bearer  string
	created time.Time
	key     *jwt.ECDSASHA
}

type ApiPayload struct {
	Aud string `json:"aud,omitempty"`
	Exp int64  `json:"exp,omitempty"`
	Iss string `json:"iss,omitempty"`
}

func (t *Token) getAuthorization() (string, error) {
	if t.key == nil {
		pk, err := ReadPrivate([]byte(t.Secret))
		if err != nil {
			return "", err
		}
		t.key = jwt.NewES256(jwt.ECDSAPrivateKey(pk))
	}

	now := time.Now()
	since := now.Sub(t.created).Seconds()
	if since > 3000 || t.bearer == "" {
		t.created = now
		p1 := &ApiPayload{
			Aud: "appstoreconnect-v1",
			// 注意,apple建议20分钟，最长支持60分钟
			Exp: now.Add(time.Minute * 20).Unix(),
			Iss: t.Iss,
		}
		bearer, err := jwt.Sign(p1, t.key, jwt.KeyID(t.Kid))
		if err != nil {
			return "", err
		}
		t.bearer = string(bearer)
	}
	return t.bearer, nil
}

func (t *Token) Verify(bearer string) (jwt.Header, error) {
	pl := &ApiPayload{}
	return jwt.Verify([]byte(bearer), t.key, pl)
}

func (t *Token) WebGet(url string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	auth, _ := t.getAuthorization()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+auth)
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil
	}
	//if resp.StatusCode != http.StatusOK {
	//	return nil, fmt.Errorf("invalid respoinse code: %s", resp.Status)
	//}
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (t *Token) WebPost(url string, reqJson []byte) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqJson))
	auth, _ := t.getAuthorization()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+auth)
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil
	}
	//if resp.StatusCode != http.StatusOK {
	//	return nil, fmt.Errorf("invalid respoinse code: %s", resp.Status)
	//}
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}
