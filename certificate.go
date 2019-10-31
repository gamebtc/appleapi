package appleapi

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

type ListCertificatesQuery struct {
	Certificates    string `json:"certificates,omitempty"` //Possible values: certificateContent, certificateType, csrContent, displayName, expirationDate, name, platform, serialNumber
	Id              string `json:"id,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
	Limit           int    `json:"limit,omitempty"`           //Maximum: 200
	Sort            string `json:"sort,omitempty"`            //Possible values: certificateType, -certificateType, displayName, -displayName, id, -id, serialNumber, -serialNumber
	CertificateType string `json:"certificateType,omitempty"` //Possible values: IOS_DEVELOPMENT, IOS_DISTRIBUTION, MAC_APP_DISTRIBUTION, MAC_INSTALLER_DISTRIBUTION, MAC_APP_DEVELOPMENT, DEVELOPER_ID_KEXT, DEVELOPER_ID_APPLICATION
	DisplayName     string `json:"displayName,omitempty"`
}

func (q *ListCertificatesQuery) QueryString() string {
	sb := strings.Builder{}
	if "" != q.Certificates {
		sb.WriteString("&fields[certificates]=")
		sb.WriteString(q.Certificates)
	}
	if "" != q.Id {
		sb.WriteString("&filter[id]=")
		sb.WriteString(q.Id)
	}
	if "" != q.SerialNumber {
		sb.WriteString("&filter[serialNumber]=")
		sb.WriteString(q.SerialNumber)
	}
	if q.Limit > 0 {
		sb.WriteString("&limit=")
		sb.WriteString(strconv.Itoa(q.Limit))
	}
	if "" != q.Sort {
		sb.WriteString("&sort=")
		sb.WriteString(q.Sort)
	}
	if "" != q.CertificateType {
		sb.WriteString("&filter[certificateType]=")
		sb.WriteString(q.CertificateType)
	}
	if "" != q.DisplayName {
		sb.WriteString("&filter[displayName]=")
		sb.WriteString(q.DisplayName)
	}
	if sb.Len() > 1 {
		return sb.String()[1:sb.Len()]
	}
	return sb.String()
}

type CertificateCreateRequest struct {
	Data struct {
		Attributes struct {
			/// Possible Values:
			/// IOS_DEVELOPMENT,
			/// IOS_DISTRIBUTION,
			/// MAC_APP_DISTRIBUTION,
			/// MAC_INSTALLER_DISTRIBUTION,
			/// MAC_APP_DEVELOPMENT,
			/// DEVELOPER_ID_KEXT,
			/// DEVELOPER_ID_APPLICATION
			CertificateType string `json:"certificateType,omitempty"`
			CsrContent      string `json:"csrContent,omitempty"`
		} `json:"attributes,omitempty"`
		Type string `json:"type,omitempty"` // certificates
	} `json:"data,omitempty"`
}

type Certificate struct {
	Attributes struct {
		CertificateContent string `json:"certificateContent,omitempty"`
		DisplayName        string `json:"displayName,omitempty"`
		ExpirationDate     string `json:"expirationDate,omitempty"`
		Name               string `json:"name,omitempty"`
		Platform           string `json:"platform,omitempty"`
		SerialNumber       string `json:"serialNumber,omitempty"`
		CertificateType    string `json:"certificateType,omitempty"` //Possible Values IOS_DEVELOPMENT IOS_DISTRIBUTION MAC_APP_DISTRIBUTION MAC_INSTALLER_DISTRIBUTION MAC_APP_DEVELOPMENT DEVELOPER_ID_KEXT DEVELOPER_ID_APPLICATION
	} `json:"attributes,omitempty"`
	Id    string        `json:"id,omitempty"`
	Type  string        `json:"type,omitempty"` //  "certificates"
	Links ResourceLinks `json:"links,omitempty"`
}

func (c *Certificate) SaveCertificateContent(file string) error {
	cerBin, err := base64.StdEncoding.DecodeString(c.Attributes.CertificateContent)
	if err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	f.Write(cerBin)
	f.Close()
	return nil
}

type CertificateResponse struct {
	Data  Certificate   `json:"data,omitempty"`
	Links DocumentLinks `json:"links,omitempty"`
}

type CertificatesResponse struct {
	Data  []Certificate      `json:"data,omitempty"`
	Links PagedDocumentLinks `json:"links,omitempty"`
	Meta  PagingInformation  `json:"meta,omitempty"`
}

type Certificates struct {
	*Token
}

// https://developer.apple.com/documentation/appstoreconnectapi/list_and_download_certificates
func (c *Certificates) Query(query *ListCertificatesQuery) ([]byte, error) {
	url := "https://api.appstoreconnect.apple.com/v1/certificates"
	if query != nil {
		url += "?" + query.QueryString()
	}
	return c.WebGet(url)
}

// 增加证书
// https://developer.apple.com/documentation/appstoreconnectapi/create_a_certificate
func (c *Certificates) CertificateCreate(csrContent, certificateType string) ([]byte, error) {
	req := new(CertificateCreateRequest)
	req.Data.Type = "certificates"
	req.Data.Attributes.CsrContent = csrContent
	req.Data.Attributes.CertificateType = certificateType

	url := "https://api.appstoreconnect.apple.com/v1/certificates"
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return c.WebPost(url, reqJson)
}
