package gcp

import (
	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"context"
	"log"
)

type AccountType string

const (
	User           AccountType = "USR"
	ServiceAccount AccountType = "SAC"
)

func generateCreds() {
	ctx := context.Background()

	c, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create IAM credentials client: %v", err)
	}
	defer c.Close()

	req := &credentialspb.GenerateAccessTokenRequest{
		// TODO: Fill request struct fields.
		// See https://pkg.go.dev/cloud.google.com/go/iam/credentials/apiv1/credentialspb#GenerateAccessTokenRequest.
	}

	resp, err := c.GenerateAccessToken(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}
