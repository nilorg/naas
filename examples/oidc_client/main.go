package main

import (
	"context"
	"log"

	"github.com/coreos/go-oidc/v3/oidc"
)

func main() {
	ctx := context.Background()
	var (
		provider *oidc.Provider
		idToken  *oidc.IDToken
		err      error
	)
	provider, err = oidc.NewProvider(ctx, "http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: "1000",
	})
	rawIDToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiMTAwMCJdLCJleHAiOjE2MTg0ODI2OTIsImlhdCI6MTYxODQ3OTA5MiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwIiwibmJmIjoxNjE4NDc5MDkyLCJzY29wZSI6InByb2ZpbGUgZW1haWwgcGhvbmUiLCJzdWIiOiIxIn0.lxN7DXLiXEZxfb1OJmOYxiI3-4BaUQfdD1UjctqgrX25YeoHDGa52eBRlCx-ZTrbM5FkUqXen25A20aImyx2GDdyUKtMLZaaJg-8unhOIpHafjw2T__CLSpcxCg0BHf0Ncf8jw3dwf6nDJJAAeEIWK1K8pWOO_Zg5Fz9pvTVDojoCIgy8Kz1hmlFcmAiRvxE--4nfi_E1aRNliYJ4ZxCKHCtJIzo5_9OeXsQQVNjOLR36qClJhaugHWAX3R2_1__EZRMQ59JlvG2Ox89IHXrARlzv43oRQnyI0FwYboeqNHkGhJrB1ebPNanLgettT4D1Y31n3puBePLjyJDhW01wA"
	idToken, err = verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Fatalf("Failed to get verifier.Verify: %s", err)
		return
	}
	log.Printf("%+v\n", idToken)
}
