package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)

func mustNewAWSSession(accessKeyID, secretAccessKey, regionID string) *session.Session {
	sess, err := newAWSSession(accessKeyID, secretAccessKey, regionID)
	if err != nil {
		log.WithError(err).
			Fatal("Failed to create AWS session")
	}

	return sess
}

func newAWSSession(accessKeyID, secretAccessKey, regionID string) (*session.Session, error) {
	awsConf := aws.Config{
		Credentials: credentials.NewStaticCredentials(
			accessKeyID, secretAccessKey, "",
		),
		Region: aws.String(regionID),
	}

	sess, err := session.NewSession(&awsConf)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
