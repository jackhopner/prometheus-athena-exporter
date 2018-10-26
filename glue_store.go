package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	log "github.com/sirupsen/logrus"
)

type glueStore struct {
	clients []glueClient
}

func mustNewGlueStore(tcs []tenantConfig, ssmPrefix, clusterID, regionID string, awsSession *session.Session) glueStore {
	clients := []glueClient{}
	for _, tc := range tcs {
		tenant := tc.Tenant
		log.Debugf("Setting up new glue store for tenant %s", tenant)
		ssmClient := ssm.New(awsSession)

		accessKeyIDOutput, err := ssmClient.GetParameter(
			&ssm.GetParameterInput{
				Name:           aws.String(fmt.Sprintf("%s%s-id", ssmPrefix, tenant)),
				WithDecryption: aws.Bool(true),
			},
		)
		if err != nil {
			log.WithError(err).Fatalf("Failed to get access key id for %s from parameter store", tenant)
		}
		secretAccessKeyOutput, err := ssmClient.GetParameter(
			&ssm.GetParameterInput{
				Name:           aws.String(fmt.Sprintf("%s%s-secret", ssmPrefix, tenant)),
				WithDecryption: aws.Bool(true),
			},
		)
		if err != nil {
			log.WithError(err).Fatalf("Failed to get secret access key for %s from parameter store", tenant)
		}

		client := mustNewGlueClient(
			*accessKeyIDOutput.Parameter.Value,
			*secretAccessKeyOutput.Parameter.Value,
			regionID,
			fmt.Sprintf("s3://aws-athena-query-results-%s-%s", clusterID, regionID),
			tc.DBName,
			tc.Tenant,
		)

		clients = append(clients, client)
	}

	return glueStore{clients}
}
