package sms

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"strings"
)

// AliService use aliyun send sms message
var AliService AliYunSMSService

// AliYunSMSService aliyun sms service
type AliYunSMSService struct {
	client             *dysmsapi.Client
	RegionID           string
	AccessKeyID        string
	AccessSecret       string
	verifyCodeSignName string
	verifyCodeTemplate string
}

func (sms *AliYunSMSService) forTest() (inTest bool) {
	if sms.AccessKeyID == "dummy" {
		inTest = true
	}
	if strings.HasPrefix(sms.AccessKeyID, "dummy") {
		inTest = true
	}
	return
}

// ConfigAuth alloc new client for sms api
func (sms *AliYunSMSService) ConfigAuth(regionID, keyID, secret string) (err error) {
	if regionID == "" {
		regionID = "cn-hangzhou"
	}
	sms.RegionID = regionID
	sms.AccessKeyID = keyID
	sms.AccessSecret = secret
	sms.client, err = dysmsapi.NewClientWithAccessKey(sms.RegionID, sms.AccessKeyID, sms.AccessSecret)
	if err != nil {
		sms.client = nil
	}
	return
}

// SetupVerifyCode setup verify code template
func (sms *AliYunSMSService) SetupVerifyCode(signName, template string) {
	sms.verifyCodeSignName = signName
	sms.verifyCodeTemplate = template
}

// SendSms send template message
// 接收短信的手机号码。
// 格式：
// 国内短信：11位手机号码，例如15951955195。
// 国际/港澳台消息：国际区号+号码，例如85200000000。
// 支持对多个手机号码发送短信，手机号码之间以英文逗号（,）分隔。上限为1000个手机号码。批量调用相对于单条调用及时性稍有延迟。
func (sms *AliYunSMSService) SendSms(template string, param interface{}, toNum string) (id string, err error) {
	var templateParam string
	var vv []byte
	if sms.client == nil {
		err = fmt.Errorf("sms service does not configed")
		return
	}
	if sms.forTest() {
		id = "C3D3EAA5-9EDC-4546-8AD7-AABBCC700000"
		return
	}
	if vv, err = json.Marshal(param); err != nil {
		err = fmt.Errorf("invalid param:%w", err)
		return
	}
	templateParam = string(vv)
	var response *dysmsapi.SendSmsResponse
	var request *dysmsapi.SendSmsRequest
	request = dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = sms.verifyCodeSignName
	request.TemplateCode = template
	request.PhoneNumbers = toNum
	request.TemplateParam = templateParam // 短信模板变量对应的实际值，JSON格式
	response, err = sms.client.SendSms(request)
	if err != nil {
		return
	}
	// {"Message":"OK","RequestId":"C3D3EAA5-9EDC-4546-8AD7-28DA6F784562","BizId":"379318983464073993^0","Code":"OK"}
	// BizId 发送回执ID，可根据该ID在接口QuerySendDetails中查询具体的发送状态
	// RequestId 请求ID
	if response.Code == "OK" {
		id = response.RequestId
	} else {
		err = fmt.Errorf("send sms fail, code %s, message %s", response.Code, response.Message)
	}
	return
}

// SendVerifyCode send verify code to phone num
func (sms *AliYunSMSService) SendVerifyCode(code, toNum string) (id string, err error) {
	var param struct {
		Code string `json:"code"`
	}
	param.Code = code
	return sms.SendSms(sms.verifyCodeTemplate, param, toNum)
}
