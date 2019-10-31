package appleapi

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

type ProfileCreateRequest struct {
	Data struct {
		Attributes struct {
			Name string `json:"name,omitempty"`
			//Possible values:
			// IOS_APP_DEVELOPMENT,
			// IOS_APP_STORE,
			// IOS_APP_ADHOC,
			// IOS_APP_INHOUSE,
			// MAC_APP_DEVELOPMENT,
			// MAC_APP_STORE,
			// MAC_APP_DIRECT,
			// TVOS_APP_DEVELOPMENT,
			// TVOS_APP_STORE,
			// TVOS_APP_ADHOC,
			// TVOS_APP_INHOUSE
			ProfileType string `json:"profileType,omitempty"`
		} `json:"attributes,omitempty"`
		Relationships struct {
			BundleId struct {
				Data TypeId `json:"data,omitempty"` //bundleIds
			} `json:"bundleId,omitempty"`

			Certificates struct {
				Data []TypeId `json:"data,omitempty"` //certificates
			} `json:"certificates,omitempty"`

			Devices struct {
				Data []TypeId `json:"data,omitempty"` //devices
			} `json:"devices,omitempty"`
		} `json:"relationships,omitempty"`

		Type string `json:"type,omitempty"` //"profiles"

	} `json:"data,omitempty"`
}

type Profile struct {
	Attributes struct {
		Name           string `json:"name,omitempty"`
		Platform       string `json:"platform,omitempty"`
		ProfileContent string `json:"profileContent,omitempty"`
		Uuid           string `json:"uuid,omitempty"`
		CreatedDate    string `json:"createdDate,omitempty"`
		ProfileState   string `json:"profileState,omitempty"` //Possible values: ACTIVE, INVALID
		ProfileType    string `json:"profileType,omitempty"`  //Possible values: IOS_APP_DEVELOPMENT, IOS_APP_STORE, IOS_APP_ADHOC, IOS_APP_INHOUSE, MAC_APP_DEVELOPMENT, MAC_APP_STORE, MAC_APP_DIRECT, TVOS_APP_DEVELOPMENT, TVOS_APP_STORE, TVOS_APP_ADHOC, TVOS_APP_INHOUSE
		ExpirationDate string `json:"expirationDate,omitempty"`
	} `json:"attributes,omitempty"`
	Id string `json:"id,omitempty"`

	Relationships struct {
		Certificates struct {
			Data  []TypeId `json:"data,omitempty"` //certificates
			Links struct {
				Related string `json:"related,omitempty"`
				Self    string `json:"self,omitempty"`
			} `json:"links,omitempty"`
			Meta PagingInformation `json:"meta,omitempty"`
		} `json:"certificates,omitempty"`

		Devices struct {
			Data  []TypeId `json:"data,omitempty"` //devices
			Links struct {
				Related string `json:"related,omitempty"`
				Self    string `json:"self,omitempty"`
			} `json:"links,omitempty"`
			Meta PagingInformation `json:"meta,omitempty"`
		} `json:"devices,omitempty"`

		BundleId struct {
			Data  TypeId `json:"data,omitempty"` //bundleIds
			Links struct {
				Related string `json:"related,omitempty"`
				Self    string `json:"self,omitempty"`
			} `json:"links,omitempty"`
		} `json:"bundleId,omitempty"`
	} `json:"relationships,omitempty"`

	Type  string        `json:"type,omitempty"` // profiles
	Links ResourceLinks `json:"links,omitempty"`
}

func (c *Profile) SaveProfileContent(file string) error {
	cerBin, err := base64.StdEncoding.DecodeString(c.Attributes.ProfileContent)
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

type ProfileResponse struct {
	Data     []Profile     `json:"data,omitempty"`
	Links    DocumentLinks `json:"links,omitempty"`
	Included []interface{} `json:"included,omitempty"`
}

// https://developer.apple.com/documentation/appstoreconnectapi/profilesresponse
type ProfilesResponse struct {
	Data     []Profile          `json:"data,omitempty"`
	Links    PagedDocumentLinks `json:"links,omitempty"`
	Meta     PagingInformation  `json:"meta,omitempty"`
	Included []interface{}      `json:"included,omitempty"`
}

type ListProfilesQuery struct {
	Certificates      string `json:"certificates,omitempty"` //Possible values: certificateContent, certificateType, csrContent, displayName, expirationDate, name, platform, serialNumber
	Devices           string `json:"devices,omitempty"`      //Possible values: addedDate, deviceClass, model, name, platform, status, udid
	Profiles          string `json:"profiles,omitempty"`     //Possible values: bundleId, certificates, createdDate, devices, expirationDate, name, platform, profileContent, profileState, profileType, uuid
	Id                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	Include           string `json:"include,omitempty"`            //Possible values: bundleId, certificates, devices
	Limit             int    `json:"limit,omitempty"`              //Maximum: 200
	LimitCertificates int    `json:"limit_certificates,omitempty"` //Maximum: 50
	LimitDevices      int    `json:"limit_devices,omitempty"`      //Maximum: 50
	Sort              string `json:"sort,omitempty"`               //Possible values: id, -id, name, -name, profileState, -profileState, profileType, -profileType
	BundleIds         string `json:"bundleIds,omitempty"`          //Possible values: bundleIdCapabilities, identifier, name, platform, profiles, seedId
	ProfileState      string `json:"profileState,omitempty"`       //Possible values: ACTIVE, INVALID
	ProfileType       string `json:"profileType,omitempty"`        //Possible values: IOS_APP_DEVELOPMENT, IOS_APP_STORE, IOS_APP_ADHOC, IOS_APP_INHOUSE, MAC_APP_DEVELOPMENT, MAC_APP_STORE, MAC_APP_DIRECT, TVOS_APP_DEVELOPMENT, TVOS_APP_STORE, TVOS_APP_ADHOC, TVOS_APP_INHOUSE
}

func (q *ListProfilesQuery) QueryString() string {
	sb := strings.Builder{}

	if "" != q.Certificates {
		sb.WriteString("&fields[certificates]=")
		sb.WriteString(q.Certificates)
	}
	if "" != q.Devices {
		sb.WriteString("&fields[devices]=")
		sb.WriteString(q.Devices)
	}
	if "" != q.Profiles {
		sb.WriteString("&fields[profiles]=")
		sb.WriteString(q.Profiles)
	}
	if "" != q.Id {
		sb.WriteString("&filter[id]=")
		sb.WriteString(q.Id)
	}
	if "" != q.Name {
		sb.WriteString("&filter[name]=")
		sb.WriteString(q.Name)
	}
	if "" != q.Include {
		sb.WriteString("&include=")
		sb.WriteString(q.Include)
	}
	if q.Limit > 0 {
		sb.WriteString("&limit=")
		sb.WriteString(strconv.Itoa(q.Limit))
	}
	if q.LimitCertificates > 0 {
		sb.WriteString("&limit[certificates]=")
		sb.WriteString(strconv.Itoa(q.LimitCertificates))
	}
	if q.LimitDevices > 0 {
		sb.WriteString("&limit[devices]=")
		sb.WriteString(strconv.Itoa(q.LimitDevices))
	}
	if "" != q.Sort {
		sb.WriteString("&sort=")
		sb.WriteString(q.Sort)
	}
	if "" != q.BundleIds {
		sb.WriteString("&fields[bundleIds]=")
		sb.WriteString(q.BundleIds)
	}
	if "" != q.ProfileState {
		sb.WriteString("&filter[profileState]=")
		sb.WriteString(q.ProfileState)
	}
	if "" != q.ProfileType {
		sb.WriteString("&filter[profileType]=")
		sb.WriteString(q.ProfileType)
	}

	if sb.Len() > 1 {
		return sb.String()[1:sb.Len()]
	}
	return sb.String()
}

type Profiles struct {
	*Token
}

func (c *Profiles) Query(query *ListProfilesQuery) ([]byte, error) {
	url := "https://api.appstoreconnect.apple.com/v1/profiles"
	if query != nil {
		url += "?" + query.QueryString()
	}
	return c.WebGet(url)
}

func (c *Profiles) ReadProfile(id string) ([]byte, error) {
	url := "https://api.appstoreconnect.apple.com/v1/profiles/" + id
	return c.WebGet(url)
}

func (c *Profiles) ProfileCreate(name, bundleId string, certificates, devices []string) ([]byte, error) {
	req := new(ProfileCreateRequest)
	req.Data.Type = "profiles"
	req.Data.Attributes.Name = name
	req.Data.Relationships.BundleId.Data.Id = bundleId
	req.Data.Relationships.BundleId.Data.Type = "bundleIds"
	req.Data.Relationships.Certificates.Data = make([]TypeId, len(certificates))
	req.Data.Relationships.Devices.Data = make([]TypeId, len(devices))

	for i := 0; i < len(certificates); i++ {
		req.Data.Relationships.Certificates.Data[i].Type = "certificates"
		req.Data.Relationships.Certificates.Data[i].Id = certificates[i]
	}

	for i := 0; i < len(devices); i++ {
		req.Data.Relationships.Devices.Data[i].Type = "devices"
		req.Data.Relationships.Devices.Data[i].Id = devices[i]
	}

	url := "https://api.appstoreconnect.apple.com/v1/profiles"
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return c.WebPost(url, reqJson)
}
