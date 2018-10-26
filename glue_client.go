//go:generate mockgen -source=glue_client.go -package=main -destination=mock_glue_client_for_test.go
package main

import (
	"database/sql"
	"time"

	athena "github.com/segmentio/go-athena"
	log "github.com/sirupsen/logrus"
)

type glueClient interface {
	runQuery(query string) (map[string]interface{}, error)
	getDBName() string
	getTenant() string
	getDB() *sql.DB
}

type glueClientImpl struct {
	db     *sql.DB
	dbName string
	tenant string
}

func (gc *glueClientImpl) getDB() *sql.DB {
	return gc.db
}

func (gc *glueClientImpl) getDBName() string {
	return gc.dbName
}

func (gc *glueClientImpl) getTenant() string {
	return gc.tenant
}

func (gc *glueClientImpl) runQuery(query string) (map[string]interface{}, error) {
	log.WithField("query", query).Debug("Running query")

	rows, err := gc.db.Query(query)
	if err != nil {
		return nil, err
	}
	cols, _ := rows.Columns()

	m := map[string]interface{}{}
	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return m, err
		}

		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
	}

	return m, nil
}

func mustNewGlueClient(
	accessKeyID,
	secretAccessKey,
	regionID,
	outputLocation,
	database,
	tenant string,
) glueClient {
	awsSession := mustNewAWSSession(accessKeyID, secretAccessKey, regionID)

	cfg := athena.Config{
		Session:        awsSession,
		Database:       database,
		OutputLocation: outputLocation,
		PollFrequency:  time.Second * 5,
	}

	db, err := athena.Open(cfg)
	if err != nil {
		log.WithError(err).Fatal("Failed to open connection to athena")
	}

	return &glueClientImpl{
		db,
		database,
		tenant,
	}
}
