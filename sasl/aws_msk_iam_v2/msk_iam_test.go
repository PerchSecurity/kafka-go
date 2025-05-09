package aws_msk_iam_v2

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/PerchSecurity/kafka-go/sasl"
	"github.com/stretchr/testify/assert"

	signer "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

const (
	accessKeyId     = "ACCESS_KEY"
	secretAccessKey = "SECRET_KEY"
)

// using a fixed time allows the signature to be verifiable in a test
var signTime = time.Date(2021, 10, 14, 13, 5, 0, 0, time.UTC)

func TestAwsMskIamMechanism(t *testing.T) {
	creds := credentials.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, "")
	ctxWithMetadata := func() context.Context {
		return sasl.WithMetadata(context.Background(), &sasl.Metadata{
			Host: "localhost",
			Port: 9092,
		})
	}

	tests := []struct {
		description string
		ctx         func() context.Context
		shouldFail  bool
	}{
		{
			description: "with metadata",
			ctx:         ctxWithMetadata,
		},
		{
			description: "without metadata",
			ctx: func() context.Context {
				return context.Background()
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			ctx := tt.ctx()

			mskMechanism := &Mechanism{
				Signer:      signer.NewSigner(),
				Credentials: creds,
				Region:      "us-east-1",
				SignTime:    signTime,
			}
			sess, auth, err := mskMechanism.Start(ctx)
			if tt.shouldFail { // if error is expected
				if err == nil { // but we don't find one
					t.Fatal("error expected")
				} else { // but we do find one
					return // return early since the remaining assertions are irrelevant
				}
			} else { // if error is not expected (typical)
				if err != nil { // but we do find one
					t.Fatal(err)
				}
			}

			if sess != mskMechanism {
				t.Error(
					"Unexpected session",
					"expected", mskMechanism,
					"got", sess,
				)
			}

			expectedMap := map[string]string{
				"version":             "2020_10_22",
				"action":              "kafka-cluster:Connect",
				"host":                "localhost",
				"user-agent":          signUserAgent,
				"x-amz-algorithm":     "AWS4-HMAC-SHA256",
				"x-amz-credential":    "ACCESS_KEY/20211014/us-east-1/kafka-cluster/aws4_request",
				"x-amz-date":          "20211014T130500Z",
				"x-amz-expires":       "300",
				"x-amz-signedheaders": "host",
				"x-amz-signature":     "6b8d25f9b45b9c7db9da855a49112d80379224153a27fd279c305a5b7940d1a7",
			}
			expectedAuth, err := json.Marshal(expectedMap)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(expectedAuth, auth) {
				t.Error("Unexpected authentication",
					"expected", expectedAuth,
					"got", auth,
				)
			}
		})
	}
}

func TestDefaultExpiry(t *testing.T) {
	expiry := time.Second * 5
	testCases := map[string]struct {
		Expiry time.Duration
	}{
		"with default":    {Expiry: expiry},
		"without default": {},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := defaultExpiry(testCase.Expiry)
			if testCase.Expiry == 0 {
				assert.Equal(t, time.Minute*5, actual)
			} else {
				assert.Equal(t, expiry, actual)
			}

		})
	}
}

func TestDefaultSignTime(t *testing.T) {
	testCases := map[string]struct {
		SignTime time.Time
	}{
		"with default":    {SignTime: signTime},
		"without default": {},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := defaultSignTime(testCase.SignTime)
			if testCase.SignTime.IsZero() {
				assert.True(t, actual.After(signTime))
			} else {
				assert.Equal(t, signTime, actual)
			}
		})
	}
}

func TestNewMechanism(t *testing.T) {
	region := "us-east-1"
	creds := credentials.StaticCredentialsProvider{}
	awsCfg := aws.Config{
		Region:      region,
		Credentials: creds,
	}
	m := NewMechanism(awsCfg)
	assert.Equal(t, m.Region, region)
	assert.Equal(t, m.Credentials, creds)
}
