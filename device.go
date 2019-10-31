package appleapi

import (
	"encoding/json"
	"strconv"
	"strings"
)

type DeviceCreateRequest struct {
	Data struct {
		Attributes struct {
			Name     string `json:"name,omitempty"`
			Platform string `json:"platform,omitempty"`
			Udid     string `json:"udid,omitempty"`
		} `json:"attributes,omitempty"`
		Type string `json:"type,omitempty"` // "devices"
	} `json:"data,omitempty"`
}

type DeviceUpdateRequest struct {
	Data struct {
		Attributes struct {
			Name   string `json:"name,omitempty"`
			Status string `json:"status,omitempty"` // Possible values: ENABLED, DISABLED
		} `json:"attributes,omitempty"`
		Id   string `json:"id,omitempty"`
		Type string `json:"type,omitempty"` // "devices"
	} `json:"data,omitempty"`
}

type Device struct {
	Attributes struct {
		DeviceClass string `json:"deviceClass,omitempty"` //Possible values: APPLE_WATCH, IPAD, IPHONE, IPOD, APPLE_TV, MAC
		Model       string `json:"model,omitempty"`
		Name        string `json:"name,omitempty"`
		Platform    string `json:"platform,omitempty"` // IOS,MAC_OS
		Status      string `json:"status,omitempty"`   // Possible values: ENABLED, DISABLED
		Udid        string `json:"udid,omitempty"`
		AddedDate   string `json:"addedDate,omitempty"`
	} `json:"attributes,omitempty"`
	Id    string        `json:"id,omitempty"`
	Type  string        `json:"type,omitempty"` // "devices"
	Links ResourceLinks `json:"links,omitempty"`
}

type DeviceResponse struct {
	Data  Device        `json:"data,omitempty"`
	Links DocumentLinks `json:"links,omitempty"`
}

type DevicesResponse struct {
	Data  []Device          `json:"data,omitempty"`
	Links DocumentLinks     `json:"links,omitempty"`
	Mate  PagingInformation `json:"mate,omitempty"`
}

// https://developer.apple.com/documentation/appstoreconnectapi/list_devices
type ListDevicesQuery struct {
	Devices  string `json:"devices,omitempty"` //Possible values: addedDate, deviceClass, model, name, platform, status, udid
	Id       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Platform string `json:"platform,omitempty"`
	Status   string `json:"status,omitempty"`
	Udid     string `json:"udid,omitempty"`
	Limit    int    `json:"limit,omitempty"` //Maximum: 200
	Sort     string `json:"sort,omitempty"`  //Possible values: id, -id, name, -name, platform, -platform, status, -status, udid, -udid
}

func (q *ListDevicesQuery) QueryString() string {
	sb := strings.Builder{}
	if "" != q.Devices {
		sb.WriteString("&fields[devices]=")
		sb.WriteString(q.Devices)
	}
	if "" != q.Id {
		sb.WriteString("&filter[id]=")
		sb.WriteString(q.Id)
	}
	if "" != q.Name {
		sb.WriteString("&filter[name]=")
		sb.WriteString(q.Name)
	}
	if "" != q.Platform {
		sb.WriteString("&filter[platform]=")
		sb.WriteString(q.Platform)
	}
	if "" != q.Status {
		sb.WriteString("&filter[status]=")
		sb.WriteString(q.Status)
	}
	if "" != q.Udid {
		sb.WriteString("&filter[udid]=")
		sb.WriteString(q.Udid)
	}
	if q.Limit > 0 {
		sb.WriteString("&limit=")
		sb.WriteString(strconv.Itoa(q.Limit))
	}
	if "" != q.Sort {
		sb.WriteString("&sort=")
		sb.WriteString(q.Sort)
	}
	if sb.Len() > 1 {
		return sb.String()[1:sb.Len()]
	}
	return sb.String()
}

type Devices struct {
	*Token
}

// https://developer.apple.com/documentation/appstoreconnectapi/list_devices
func (c *Devices) Query(query *ListDevicesQuery) ([]byte, error) {
	url := "https://api.appstoreconnect.apple.com/v1/devices"
	if query != nil {
		url += "?" + query.QueryString()
	}
	return c.WebGet(url)
}

// 增加设备
// https://developer.apple.com/documentation/appstoreconnectapi/register_a_new_device
func (c *Devices) DeviceCreate(udid, name string) ([]byte, error) {
	req := new(DeviceCreateRequest)
	req.Data.Type = "devices"
	req.Data.Attributes.Platform = PlatformIos
	req.Data.Attributes.Name = name
	req.Data.Attributes.Udid = udid

	url := "https://api.appstoreconnect.apple.com/v1/devices"
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return c.WebPost(url, reqJson)
}

// 更新设备
// https://developer.apple.com/documentation/appstoreconnectapi/modify_a_registered_device
func (c *Devices) DeviceUpdate(id, name, status string) ([]byte, error) {
	req := new(DeviceUpdateRequest)
	req.Data.Type = "devices"
	req.Data.Id = id
	req.Data.Attributes.Name = name
	req.Data.Attributes.Status = status

	url := "https://api.appstoreconnect.apple.com/v1/devices/" + id
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return c.WebPost(url, reqJson)
}
