package webhookresponse

import "errors"

type WebhookResponseTemplate struct {
	URL       string
	Method    string `validate:"oneof=GET POST PUT PATCH DELETE HEAD OPTIONS"`
	BodyType  string `validate:"oneof=none form-data x-www-form-urlencoded raw json binary"`
	UrlParams []map[string]string
	Headers   []map[string]string
	Body      interface{} `validate:"required_unless=BodyType none"`
	Cookies   []map[string]string
}

type RawBody struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (t *WebhookResponseTemplate) ReflectBodyToRaw() *RawBody {
	rbd := &RawBody{}
	rb, _ := t.Body.(map[string]interface{})
	for k, v := range rb {
		switch k {
		case "type":
			rbd.Type, _ = v.(string)
		case "content":
			rbd.Content, _ = v.(string)
		}
	}
	return rbd
}

func resolveIntFieldsFromActionOptions(actionOptions map[string]interface{}, fieldName string) (int, error) {
	raw, hit := actionOptions[fieldName]
	if !hit {
		return 0, errors.New("missing " + fieldName + " field")

	}
	numberInFloat, numberAssertPass := raw.(float64)
	number := int(numberInFloat)
	if !numberAssertPass {
		return 0, errors.New(fieldName + " field which in action options assert failed")

	}
	return number, nil
}
