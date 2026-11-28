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
		domain string
		signer *sign.URLSigner
	}

	// SignedURL represents a CloudFront signed URL with its expiration time.
	SignedURL struct {
		URL       string    `json:"url"`
		ExpiresAt time.Time `json:"expiresAt"`
	}
)

// NewCloudFrontStore creates a CloudFrontStore for generating signed URLs.
func NewCloudFrontStore(domain string, keyPairID string, privateKey crypto.Signer) *CloudFrontStore {

	return &CloudFrontStore{
		domain: strings.TrimSuffix(domain, "/"),
		signer: sign.NewURLSigner(keyPairID, privateKey),
	}
}

func (c *CloudFrontStore) OneHourPresignedURL(path string) (*SignedURL, error) {
	return c.PresignedURL(path, time.Now().Add(time.Hour))
}

// Generate a presigned URL for the given key with the specified expiration time.
func (c *CloudFrontStore) PresignedURL(path string, expiration time.Time) (*SignedURL, error) {

	// Url should form cleanly without duplicate slashes
	// from domain and s3 key
	url, err := url.JoinPath(c.domain, path)
	if err != nil {
		return nil, err
	}

	signedURL, err := c.signer.Sign(url, expiration)
	if err != nil {
		return nil, err
	}

	return &SignedURL{
		URL:       signedURL,
		ExpiresAt: expiration,
	}, nil
}
