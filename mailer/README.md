# Cloud Outbound Mailer

mailer package provides a client to communicate with the Naver Cloud Platform's Cloud Outbound Mailer API.

Also see the [Official API documentation](https://api.ncloud-docs.com/docs/ai-application-service-cloudoutboundmailer).

## ⚡️ Quickstart

```Go
package main

import (
	"context"

	"github.com/connectfit-team/naverapi/mailer"
)

func main() {
	client, err := mailer.NewCloudOutboundMailerClient("my-access-key", "my-secret-key", nil)
	if err != nil {
		panic(err)
	}

	req := mailer.CreateMailRequest{
		SenderAddress: "my-address@naver.com",
		SenderName:    "Me",
		Title:         "Foo",
		Body:          "Bar",
		Recipients: []*mailer.Recipient{
			{
				Address: "your-address@naver.com",
				Name:    "You",
				Type:    mailer.RecipientTypeBlindCarbonCopy,
			},
		},
		AttachFileIDs: []string{"file-id"},
	}
	resp, err := client.CreateMail(context.Background(), req)
	if err != nil {
		panic(err)
	}
}

```
