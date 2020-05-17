package main

import (
	"context"
	"log"

	"github.com/coreos/go-oidc"
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
		ClientID: "naas-oidc-test",
	})
	for i := 0; i < 10; i++ {
		idToken, err = verifier.Verify(ctx, "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsibmFhcy1vaWRjLXRlc3QiLCIxeHh4eDAwMCJdLCJleHAiOjE1ODk3NzA4NzEsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6ODA4MCIsIm5iZiI6MTQ1MTYwNjQwMCwic3ViIjoic3ViamVjdCJ9.sNW-1gf_f8G-mheD8nnfZwFwlRx-fM5wvar8iD6CEWuYFMZGa4hWWjewm8JrEtalW31EmjFVPVE7vFhKE7Z1bXrelSsuA9FApVkKClqTN6Gs9A5CNj75_6jSauJNIRqXmL39RGWyzgnfgBro7R2jMBYLkTceQvbOrREPky2uu2kBcHEFOhjuUhbKHPEnjgjyRn754axIBGcCloebC5XmscdI568C9HcL0T7V1LvyxeC2RKu8HmTtahqVDhvDOr1k0F2czlO7ivgN_yFSTuxROq8tk7MpgX1fNO3peg3DCkQVAyW5sfDoAbFCtekwTUAEwAGE3zckdK-_2iPLElJh0g")
		if err != nil {
			log.Fatalf("Failed to get verifier.Verify: %s", err)
			return
		}
		log.Printf("%+v\n", idToken)
	}
}
