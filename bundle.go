package appleapi

import (
	"encoding/json"
	"strconv"
	"strings"
)

type BundleId struct {
	Attributes struct {
		Identifier string `json:"identifier,omitempty"`
		Name       string `json:"name,omitempty"`
		Platform   string `json:"platform,omitempty"`
		SeedId     string `json:"seedId,omitempty"`
	} `json:"attributes,omitempty"`

	Id string `json:"id,omitempty"`

	Relationships struct {
		Profiles struct {
			Data  []TypeId `json:"data,omitempty"` //"profiles"
			Links struct {
				Related string `json:"related,omitempty"`
				Self    string `json:"self,omitempty"`
			} `json:"links,omitempty"`
			Meta PagingInformation `json:"meta,omitempty"`
		} `json:"profiles,omitempty"`
	} `json:"relationships,omitempty"`

	BundleIdCapabilities struct {
		Data  []TypeId `json:"data,omitempty"` //"bundleIdCapabilities"
		Links struct {
			Related string `json:"related,omitempty"`
			Self    string `json:"self,omitempty"`
		} `json:"links,omitempty"`
		Meta PagingInformation `json:"meta,omitempty"`
	} `json:"bundleIdCapabilities,omitempty"`

	Type  string        `json:"type,omitempty"` // "bundleIds"
	Links ResourceLinks `json:"links,omitempty"`
}

type ListBundlesQuery struct {
	BundleIds            string `json:"bundleIds,omitempty"` //Possible values: bundleIdCapabilities, identifier, name, platform, profiles, seedId
	Profiles             string `json:"profiles,omitempty"`  //Possible values: bundleId, certificates, createdDate, devices, expirationDate, name, platform, profileContent, profileState, profileType, uuid
	Id                   string `json:"id,omitempty"`
	Identifier           string `json:"identifier,omitempty"`
	Name                 string `json:"name,omitempty"`
	Platform             string `json:"platform,omitempty"` //Possible values: IOS, MAC_OS
	SeedId               string `json:"seedId,omitempty"`
	Include              string `json:"include,omitempty"`              //Possible values: bundleIdCapabilities, profiles
	Limit                int    `json:"limit,omitempty"`                //Maximum: 200
	LimitProfiles        int    `json:"limit_profiles,omitempty"`       //Maximum: 50
	Sort                 string `json:"sort,omitempty"`                 //Possible values: id, -id, name, -name, platform, -platform, seedId, -seedId
	BundleIdCapabilities string `json:"bundleIdCapabilities,omitempty"` //Possible values: bundleId, capabilityType, settings
}

func (q *ListBundlesQuery) QueryString() string {
	sb := strings.Builder{}
	if "" != q.BundleIds {
		sb.WriteString("&fields[bundleIds]=")
		sb.WriteString(q.BundleIds)
	}
	if "" != q.Profiles {
		sb.WriteString("&fields[profiles]=")
		sb.WriteString(q.Profiles)
	}
	if "" != q.Id {
		sb.WriteString("&filter[id]=")
		sb.WriteString(q.Id)
	}
	if "" != q.Identifier {
		sb.WriteString("&filter[identifier]=")
		sb.WriteString(q.Identifier)
	}
	if "" != q.Name {
		sb.WriteString("&filter[name]=")
		sb.WriteString(q.Name)
	}
	if "" != q.Platform {
		sb.WriteString("&filter[platform]=")
		sb.WriteString(q.Platform)
	}
	if "" != q.SeedId {
		sb.WriteString("&filter[seedId]=")
		sb.WriteString(q.SeedId)
	}
	if "" != q.Include {
		sb.WriteString("&include=")
		sb.WriteString(q.Include)
	}
	if q.Limit > 0 {
		sb.WriteString("&limit=")
		sb.WriteString(strconv.Itoa(q.Limit))
	}
	if q.LimitProfiles > 0 {
		sb.WriteString("&limit[profiles]=")
		sb.WriteString(strconv.Itoa(q.LimitProfiles))
	}
	if "" != q.Sort {
		sb.WriteString("&sort=")
		sb.WriteString(q.Sort)
	}
	if "" != q.BundleIdCapabilities {
		sb.WriteString("&fields[bundleIdCapabilities]=")
		sb.WriteString(q.BundleIdCapabilities)
	}
	if sb.Len() > 1 {
		return sb.String()[1:sb.Len()]
	}
	return sb.String()
}

type BundleIdCreateRequest struct {
	Data struct {
		Attributes struct {
			Identifier string `json:"identifier,omitempty"` // 标识符
			Name       string `json:"name,omitempty"`       // 名称
			Platform   string `json:"platform,omitempty"`   // 可能的值 IOS, MAC_OS
			SeedId     string `json:"seedId,omitempty"`
		} `json:"attributes,omitempty"`
		Type string `json:"type,omitempty"` //bundleIds
	} `json:"data,omitempty"`
}

type BundleIdUpdateRequest struct {
	Data struct {
		Attributes struct {
			Name string `json:"name,omitempty"` // 名称
		} `json:"attributes,omitempty"`
		Id   string `json:"id,omitempty"`
		Type string `json:"type,omitempty"` //bundleIds
	} `json:"data,omitempty"`
}

type BundleIdResponse struct {
	Data     BundleId      `json:"data,omitempty"`
	Links    DocumentLinks `json:"links,omitempty"`
	Included interface{}   `json:"included,omitempty"` //Possible types: Profile, BundleIdCapability
}

type BundleIdsResponse struct {
	Data     []BundleId         `json:"data,omitempty"`
	Links    PagedDocumentLinks `json:"links,omitempty"`
	Meta     PagingInformation  `json:"meta,omitempty"`
	Included interface{}        `json:"included,omitempty"` //Possible types: Profile, BundleIdCapability
}

type Bundles struct {
	*Token
}

// https://developer.apple.com/documentation/appstoreconnectapi/list_bundle_ids
func (b *Bundles) Query(query *ListBundlesQuery) ([]byte, error) {
	url := "https://api.appstoreconnect.apple.com/v1/bundleIds"
	if query != nil {
		url += "?" + query.QueryString()
	}
	return b.WebGet(url)
}

// 创建BundleId
// https://developer.apple.com/documentation/appstoreconnectapi/register_a_new_bundle_id
func (b *Bundles) BundleIdCreate(identifier, name string) ([]byte, error) {
	req := new(BundleIdCreateRequest)
	req.Data.Type = "bundleIds"
	req.Data.Attributes.Platform = PlatformIos
	req.Data.Attributes.Identifier = identifier
	req.Data.Attributes.Name = name

	url := "https://api.appstoreconnect.apple.com/v1/bundleIds"
	query, _ := json.Marshal(req)
	return b.WebPost(url, query)
}
