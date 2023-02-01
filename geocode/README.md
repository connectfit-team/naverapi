# naverapi/geocode

# Installation

`go get github.com/connectfit-team/naverapi/geocode`

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

	// find road address
	roadAddress := "경기도 성남시 분당구 불정로 6 그린팩토리"
	res, err := client.Query(ctx, roadAddress)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	
	// find jibun address
	jibunAddress := "경기도 성남시 분당구 정자동 178-1 그린팩토리"
	_, err = client.Query(ctx, jibunAddress)
	if err != nil {
		panic(err)
	}
	
	// find address as english
	res, err = client.
		Language(geocode.Eng).
		Query(ctx, jibunAddress)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
```

# References

* https://api.ncloud-docs.com/docs/ai-naver-mapsgeocoding-geocode
