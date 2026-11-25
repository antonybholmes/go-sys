package aws

import (
	"crypto"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
)

type (
	CloudFrontStore struct {
		// domain is the CloudFront distribution domain
		// without a trailing slash. Can also be a custom domain name
		// if that is associated with CloudFront via Route 53 etc.
		domain     string
		keyID      string
		privateKey crypto.Signer
		signer     *sign.URLSigner
	}

	SignedURL struct {
		URL        string
		Expiration time.Time
	}
)

// NewCloudFrontStore creates a CloudFrontStore for generating signed URLs.
func NewCloudFrontStore(domain string, keyID string, privateKey crypto.Signer) *CloudFrontStore {
	domain = strings.TrimSuffix(domain, "/")

	return &CloudFrontStore{
		domain:     domain,
		keyID:      keyID,
		privateKey: privateKey,
		signer:     sign.NewURLSigner(keyID, privateKey),
	}
}

func (c *CloudFrontStore) OneHourSignedURL(key string) (*SignedURL, error) {
	return c.GenerateSignedURL(key, time.Now().Add(time.Hour))
}

func (c *CloudFrontStore) GenerateSignedURL(key string, expiration time.Time) (*SignedURL, error) {
	// Standardize the key by removing any leading slash.
	key = strings.TrimPrefix(key, "/")

	// Url should form cleanly without duplicate slashes
	// from domain and s3 key
	url := c.domain + "/" + key

	signedURL, err := c.signer.Sign(url, expiration)
	if err != nil {
		return nil, err
	}

	return &SignedURL{
		URL:        signedURL,
		Expiration: expiration,
	}, nil
}
