package main

import (
	"fmt"
	"os"

	pubsub "cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	projectID = "core-eso"
)

func main() {

	ctx := context.Background()

	keyBytes, err := os.ReadFile("../certs/jwt-access-svc-account.json")
	if err != nil {
		panic(err)
	}

	// aud := "https://pubsub.googleapis.com/google.pubsub.v1.Publisher"
	// tokenSource, err := google.JWTAccessTokenSourceFromJSON(keyBytes, aud)
	sscope := "https://www.googleapis.com/auth/cloud-platform"
	tokenSource, err := google.JWTAccessTokenSourceWithScope(keyBytes, sscope)
	if err != nil {
		panic(err)
	}

	t, _ := tokenSource.Token()
	fmt.Printf("%s\n", t)
	pubsubClient, err := pubsub.NewClient(ctx, projectID, option.WithTokenSource(tokenSource))
	if err != nil {
		fmt.Printf("pubsub.NewClient: %v", err)
		return
	}
	defer pubsubClient.Close()

	pit := pubsubClient.Topics(ctx)
	for {
		topic, err := pit.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("pubssub.Iterating error: %v", err)
			return
		}
		fmt.Printf("Topic Name: %s\n", topic.ID())
	}

}
