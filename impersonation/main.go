package main

import (
	"fmt"
	"log"
	"time"

	"impersonatedtoken"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	projectID = "core-eso"
)

func main() {

	ctx := context.Background()

	tokenSource, err := impersonatedtoken.ImpersonatedTokenSource(&impersonatedtoken.ImpersonatedTokenConfig{
		Duration: time.Duration(10 * time.Second), // note, even if you give 10s, seems google will allow for +5mins leeway for an expired jwtaccesstoken to be used
	})
	if err != nil {
		log.Fatal(err)
	}

	t, err := tokenSource.Token()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", t.AccessToken)

	// try to use the token as a non-refreshed token and then when you loop the storage client below, it'll stop working 5mins later
	staticTokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: t.AccessToken,
		},
	)

	storageClient, err := storage.NewClient(ctx, option.WithTokenSource(staticTokenSource))
	if err != nil {
		log.Fatal(err)
	}

	// if you use the token source directly, it'll refresh from the metadata server when it expires
	// storageClient, err := storage.NewClient(ctx, option.WithTokenSource(tokenSource))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	for {
		sit := storageClient.Buckets(ctx, projectID)
		for {
			battrs, err := sit.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%s\n", battrs.Name)

		}
		time.Sleep(time.Duration(10 * time.Second))
	}
}
