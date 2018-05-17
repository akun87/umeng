package umeng

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"log"
	"github.com/heycayc/chkcki/httplib"
)

const (

	youmeng_host = "http://msg.umeng.com/api/send"

	listcast_device_tokens_max = 500
	customizedcast_alias_max   = 500

	// "type":"xx",        // 必填，消息发送类型,其值可以为:
	TYPE_UNICAST        = "unicast"        //单播
	TYPE_LISTCAST       = "listcast"       //列播，要求不超过500个device_token
	TYPE_FILECAST       = "filecast"       //文件播，多个device_token可通过文件形式批量发送
	TYPE_BROADCAST      = "broadcast"      //广播
	TYPE_GROUPCAST      = "groupcast"      //组播，按照filter筛选用户群, 请参照filter参数
	TYPE_CUSTOMIZEDCAST = "customizedcast" //通过alias进行推送，包括以下两种case:

	PAYLOAD_DISPLAY_TYPE_NOTIFICATION = "notification"
	PAYLOAD_DISPLAY_TYPE_MESSAGE      = "message"

	//"after_open": "xx",    // 可选，默认为"go_app"，值可以为:
	PAYLOAD_AFTER_OPEN_GO_APP      = "go_app"      // 打开应用
	PAYLOAD_AFTER_OPEN_GO_URL      = "go_url"      //跳转到URL
	PAYLOAD_AFTER_OPEN_GO_ACTIVITY = "go_activity" // 打开特定的activity
	PAYLOAD_AFTER_OPEN_GO_CUSTOM   = "go_custom"   // 用户自定义内容。
)

type (
	//https://developer.umeng.com/docs/66632/detail/68343#h3-u529Fu80FDu8BF4u660E
	Reqd struct {
		Appkey       string      `json:"appkey,omitempty"`
		Timestamp    int64      `json:"timestamp,omitempty"`
		Type         string      `json:"type,omitempty"`
		DeviceTokens string      `json:"device_tokens,omitempty"`
		AliasType    string      `json:"alias_type,omitempty"`
		Alias        string      `json:"alias,omitempty"`
		FileID       string      `json:"file_id,omitempty"`
		Filter       interface{} `json:"filter,omitempty"`
		Payload      struct {
			DisplayType string `json:"display_type,omitempty"`
			Body        struct {
				Ticker      string `json:"ticker,omitempty"`
				Title       string `json:"title,omitempty"`
				Text        string `json:"text,omitempty"`
				Icon        string `json:"icon,omitempty"`
				LargeIcon   string `json:"largeIcon,omitempty"`
				Img         string `json:"img,omitempty"`
				Sound       string `json:"sound,omitempty"`
				BuilderID   int    `json:"builder_id,omitempty"`
				PlayVibrate bool   `json:"play_vibrate,omitempty"`
				PlayLights  bool   `json:"play_lights,omitempty"`
				PlaySound   bool   `json:"play_sound,omitempty"`
				AfterOpen   string `json:"after_open,omitempty"`
				URL         string `json:"url,omitempty"`
				Activity    string `json:"activity,omitempty"`
				Custom      interface{} `json:"custom,omitempty"`
			} `json:"body,omitempty"`
			Extra interface{} `json:"extra,omitempty"`
		} `json:"payload,omitempty"`
		Policy struct {
			StartTime  string `json:"start_time,omitempty"`
			ExpireTime string `json:"expire_time,omitempty"`
			MaxSendNum int    `json:"max_send_num,omitempty"`
			OutBizNo   string `json:"out_biz_no,omitempty"`
		} `json:"policy,omitempty"`
		ProductionMode bool   `json:"production_mode,omitempty"`
		Description    string `json:"description,omitempty"`
		Mipush         bool   `json:"mipush,omitempty"`
		MiActivity     string `json:"mi_activity,omitempty"`
	}

	Resp struct {
		Ret  string      `json:"ret"`
		Data map[string]interface{} `json:"data"`
	}

	Context struct {
		AppKey    string
		AppSecret string
	}
)


func (this *Reqd) AddDeviceToken(token string) *Reqd {
	if this.Type == TYPE_UNICAST {
		this.DeviceTokens = token
	} else if this.Type == TYPE_LISTCAST {
		arr := strings.Split(this.DeviceTokens, ",")
		arr = append(arr, token)
		arr = arr[:listcast_device_tokens_max]
		this.DeviceTokens = strings.Join(arr, ",")
	}
	return this
}
func (this *Reqd) AddDeviceTokens(tokens []string) *Reqd {
	if len(tokens) == 0 {
		return this
	}
	if this.Type == TYPE_UNICAST {
		this.DeviceTokens = tokens[0]
	} else if this.Type == TYPE_LISTCAST {
		arr := strings.Split(this.DeviceTokens, ",")
		arr = append(arr, tokens...)
		arr = arr[:listcast_device_tokens_max]
		this.DeviceTokens = strings.Join(arr, ",")
	}
	return this
}


func (this *Reqd) SetAliasType(typ string) *Reqd {
	if this.Type != TYPE_CUSTOMIZEDCAST {
		return this
	}
	this.AliasType = typ
	return this
}

func (this *Reqd) AddAlias(alias string) *Reqd {
	if this.Type != TYPE_CUSTOMIZEDCAST {
		return this
	}

	arr := strings.Split(this.Alias, ",")
	if len(arr) > customizedcast_alias_max {
		return this
	}
	arr = append(arr, alias)
	this.Alias = strings.Join(arr, ",")

	return this
}

func (this *Reqd) SetFileId(fn string) *Reqd {
	//not support
	return this
}

func (this *Reqd) SetFilter(ft interface{}) *Reqd {
	this.Filter = ft
	return this
}

func (this *Reqd) SetPayloadDisplyType(typ string) *Reqd {
	this.Payload.DisplayType = typ
	return this
}

func (this *Reqd) SetPayloadTicker(val string) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Body.Ticker = val
	return this
}

func (this *Reqd) SetPayloadTitle(val string) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Body.Title = val
	return this
}

func (this *Reqd) SetPayloadText(val string) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Body.Text = val
	return this
}

func (this *Reqd) SetPayloadIcon(val string) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Body.Icon = val
	return this
}

func (this *Reqd) SetPayloadLargeIcon(val string) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Body.LargeIcon = val
	return this
}

func (this *Reqd) SetPayloadImg(val string) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Body.Img = val
	return this
}

func (this *Reqd) SetPayloadSound(val string) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Body.Sound = val
	return this
}

func (this *Reqd) SetPayloadBuilderID(val int) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Body.BuilderID = val
	return this
}

func (this *Reqd) SetPayloadPlayVibrate(val bool) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Body.PlayVibrate = val
	return this
}

func (this *Reqd) SetPayloadPlayLights(val bool) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Body.PlayLights = val
	return this
}

func (this *Reqd) SetPayloadPlaySound(val bool) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Body.PlaySound = val
	return this
}

func (this *Reqd) SetPayloadAfterOpen(typ string) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Body.AfterOpen = typ
	return this
}

func (this *Reqd) SetPayloadURL(url string) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	if this.Payload.Body.AfterOpen != PAYLOAD_AFTER_OPEN_GO_URL {
		return this
	}
	this.Payload.Body.URL = url
	return this
}

func (this *Reqd) SetPayloadActivity(act string) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	if this.Payload.Body.AfterOpen != PAYLOAD_AFTER_OPEN_GO_ACTIVITY {
		return this
	}
	this.Payload.Body.Activity = act
	return this
}

func (this *Reqd) SetPayloadCustom(data string) *Reqd {
	if this.Payload.DisplayType != PAYLOAD_DISPLAY_TYPE_MESSAGE && this.Payload.Body.AfterOpen != PAYLOAD_AFTER_OPEN_GO_ACTIVITY {
		return this
	}
	this.Payload.Body.Custom = data
	return this
}

func (this *Reqd) SetPayloadExtra(i interface{}) *Reqd {
	if this.Payload.DisplayType == PAYLOAD_DISPLAY_TYPE_MESSAGE {
		return this
	}
	this.Payload.Extra = i
	return this
}

//发送策略
func (this *Reqd) SetTime(start int64, expire int64)*Reqd {
	this.Policy.StartTime = time.Unix(start,0).Format("2006-01-02 15:04:05")
    this.Policy.ExpireTime= time.Unix(expire,0).Format("2006-01-02 15:04:05")
	return this
}

func (this *Reqd) SetSendMaxNum(num int) *Reqd {
	this.Policy.MaxSendNum = num
	return this
}

func (this *Reqd) SetUniqueId(id string) *Reqd {
	this.Policy.OutBizNo = id
	return this
}

func (this *Reqd) SetProductionMode(b bool) *Reqd {
	this.ProductionMode = b
	return this
}

func (this *Reqd) SetDescription(str string) *Reqd {
	this.Description = str
	return this
}

func (this *Reqd) SetMipush(b bool) *Reqd {
	this.Mipush = b
	return this
}

func (this *Reqd) SetMiActivity(str string) *Reqd {
	this.MiActivity = str
	return this
}

func NewConext(appKey, appSecret string) *Context {
	p := &Context{}
	p.AppKey = appKey
	p.AppSecret = appSecret
	return p
}

func (this *Context) url(reqd *Reqd) string {
	bodystr, _ := json.Marshal(reqd)

	sign := fmt.Sprintf("%s%s%s%s","POST",youmeng_host,bodystr,this.AppSecret)
	md5sum := md5.Sum([]byte(sign[:]))

	sign = hex.EncodeToString(md5sum[:])

	fmt.Println(sign)

	return fmt.Sprintf("%s?sign=%s",youmeng_host,sign)
}

func (this *Context) NewReqd(typ string) *Reqd {
	p := Reqd{}
	p.Appkey = this.AppKey
	p.Timestamp = time.Now().Unix()
	p.Type = typ
	return &p
}

func(this* Context) Send(reqd *Reqd) *Resp{
	url :=this.url(reqd)
	fmt.Println(url)
	req := httplib.Post(url).SetHeader("Content-Type","application/json")
	req.JSONBody(reqd)
	res := Resp{}
	err := req.ToJSON(&res)
	if err != nil {
		log.Fatalln(err)
	}
	return  &res
}
