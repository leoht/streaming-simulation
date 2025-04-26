package producer

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-msk-iam-sasl-signer-go/signer"
	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func CreateToken() kafka.OAuthBearerToken {
	token, tokenExpirationTime, err := signer.GenerateAuthToken(context.TODO(), os.Getenv("AWS_REGION"))
	if err != nil {
		panic(err)
	}
	seconds := tokenExpirationTime / 1000
	nanoseconds := (tokenExpirationTime % 1000) * 1000000
	bearerToken := kafka.OAuthBearerToken{
		TokenValue: token,
		Expiration: time.Unix(seconds, nanoseconds),
	}

	return bearerToken
}
