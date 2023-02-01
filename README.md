# naverapi

# Installation

`go get github.com/connectfit-team/naverapi`

# Example

### Query Simple Address

```golang
package main

import (
	"context"
	"fmt"

	"github.com/connectfit-team/naverapi/geocode"
)

func main() {
	ctx := context.Background()
	client, err := geocode.NewClient("[CLIENT_ID]", "[CLIENT_SECRET]", nil)
	if err != nil {
		panic(err)
	}

	targetAddress := "경기도 성남시 분당구 불정로 6 그린팩토리"
	res, err := client.Query(ctx, targetAddress)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
```

# References

* https://api.ncloud-docs.com/docs/ai-naver-mapsgeocoding-geocode
