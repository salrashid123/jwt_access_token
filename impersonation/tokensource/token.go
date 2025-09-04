package impersonatedtoken

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	credentialspb "cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	jwt "github.com/golang-jwt/jwt/v5"

	"cloud.google.com/go/compute/metadata"
	"golang.org/x/oauth2"
)

const (
	GCP_CLOUD_PLATFORM_SCOPE = "https://www.googleapis.com/auth/cloud-platform"
)

type ClaimWithSubject struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

type ImpersonatedTokenConfig struct {
	Email string // ServiceAccount Email

	KeyId    string   // The service accounts key_id value
	Scopes   []string // list of scopes to use
	Duration time.Duration
}

type impersonatedTokenSource struct {
	refreshMutex *sync.Mutex
	oauth2.TokenSource
	email    string
	keyId    string
	scopes   []string
	myToken  *oauth2.Token
	duration time.Duration
}

func ImpersonatedTokenSource(tokenConfig *ImpersonatedTokenConfig) (oauth2.TokenSource, error) {

	if !metadata.OnGCE() {
		return nil, fmt.Errorf("salrashid123/x/oauth2/google: metadata server not available")
	}
	if tokenConfig.Email == "" {
		ctx := context.Background()
		email, err := metadata.EmailWithContext(ctx, "default")
		if err != nil {
			return nil, fmt.Errorf("salrashid123/x/oauth2/google: ImpersonatedTokenConfig.Email and cannot be nil")
		}
		tokenConfig.Email = email
	}

	if len(tokenConfig.Scopes) == 0 {
		tokenConfig.Scopes = []string{GCP_CLOUD_PLATFORM_SCOPE}
	}

	return &impersonatedTokenSource{
		refreshMutex: &sync.Mutex{},
		email:        tokenConfig.Email,

		keyId:    tokenConfig.KeyId,
		scopes:   tokenConfig.Scopes,
		duration: tokenConfig.Duration,
	}, nil

}

func (ts *impersonatedTokenSource) Token() (*oauth2.Token, error) {
	ts.refreshMutex.Lock()
	defer ts.refreshMutex.Unlock()
	if ts.myToken.Valid() {
		return ts.myToken, nil
	}

	ctx := context.Background()

	type idTokenJWT struct {
		jwt.RegisteredClaims
		Scopes []string `json:"scopes"`
	}

	iat := time.Now()
	exp := iat.Add(ts.duration) // we just need a small amount of time to get a token

	claims := &ClaimWithSubject{
		Scope: strings.Join(ts.scopes, " "),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(iat),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    ts.email,
			Subject:   ts.email,
		},
	}

	b, err := json.Marshal(claims)
	if err != nil {
		return &oauth2.Token{}, err
	}

	c, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		return &oauth2.Token{}, err
	}
	defer c.Close()

	idreq := &credentialspb.SignJwtRequest{
		Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", ts.email),
		Payload: string(b),
	}
	idresp, err := c.SignJwt(ctx, idreq)
	if err != nil {
		return &oauth2.Token{}, err
	}

	ts.myToken = &oauth2.Token{AccessToken: idresp.SignedJwt, TokenType: "Bearer", Expiry: exp}

	return ts.myToken, nil
}
