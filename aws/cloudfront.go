package aws

import (
	"crypto"
	"net/url"
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
		privateKey crypto.Signer
		signer     *sign.URLSigner
	}

	SignedURL struct {
		URL        string
		Expiration time.Time
	}
)

// NewCloudFrontStore creates a CloudFrontStore for generating signed URLs.
func NewCloudFrontStore(domain string, keyPairID string, privateKey crypto.Signer) *CloudFrontStore {

	return &CloudFrontStore{
		domain:     strings.TrimSuffix(domain, "/"),
		privateKey: privateKey,
		signer:     sign.NewURLSigner(keyPairID, privateKey),
	}
}

func (c *CloudFrontStore) OneHourPresignedURL(key string) (*SignedURL, error) {
	return c.GeneratePresignedURL(key, time.Now().Add(time.Hour))
}

func (c *CloudFrontStore) GeneratePresignedURL(key string, expiration time.Time) (*SignedURL, error) {

	// Url should form cleanly without duplicate slashes
	// from domain and s3 key
	url, err := url.JoinPath(c.domain, key)
	if err != nil {
		return nil, err
	}

	signedURL, err := c.signer.Sign(url, expiration)
	if err != nil {
		return nil, err
	}

	return &SignedURL{
		URL:        signedURL,
		Expiration: expiration,
	}, nil
}
