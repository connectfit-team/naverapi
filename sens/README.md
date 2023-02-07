# SENS(Simple & Easy Notification Service)

sens package provides a client to communicate with the Naver Cloud Platform's SENS API.

Also see the [Official API documentation](https://api.ncloud-docs.com/docs/ai-application-service-sens).

## ⚡️ Quickstart

```Go
package main

import (
	"context"

	"github.com/connectfit-team/naverapi/sens"
)

func main() {
	client, err := sens.NewClient("my-access-key", "my-secret-key", nil)
	if err != nil {
		panic(err)
	}

	req := sens.SendSMSRequest{
		Type:        sens.SMSTypeLMS,
		ContentType: sens.ContentTypeAD,
		CountryCode: sens.CountryCodeKorea,
		From:        "me",
		Subject:     "my-subject",
		Content:     "my-content",
		Messages: []sens.Message{
			{
				To:      "you",
				Subject: "I'm happy",
				Content: ":)",
			},
		},
	}
	resp, err := client.SendSMS(context.Background(), req)
	if err != nil {
		panic(err)
	}
}
```
