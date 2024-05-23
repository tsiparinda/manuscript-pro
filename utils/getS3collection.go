package utils

import (
	"context"
	"database/sql"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	brconfig "github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/logging"
)

func getCredentialsFromCognito(ctx context.Context, identityPoolID string, brconfig brconfig.Configuration, logger logging.Logger) (aws.CredentialsProvider, error) {
	region := brconfig.GetStringDefault("storage:aws:region", "eu-west-1")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	cognitoSvc := cognitoidentity.NewFromConfig(cfg)
	identityResp, err := cognitoSvc.GetId(ctx, &cognitoidentity.GetIdInput{
		IdentityPoolId: aws.String(identityPoolID),
	})
	if err != nil {
		return nil, err
	}

	credsResp, err := cognitoSvc.GetCredentialsForIdentity(ctx, &cognitoidentity.GetCredentialsForIdentityInput{
		IdentityId: identityResp.IdentityId,
	})
	if err != nil {
		return nil, err
	}

	creds := aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     *credsResp.Credentials.AccessKeyId,
			SecretAccessKey: *credsResp.Credentials.SecretKey,
			SessionToken:    *credsResp.Credentials.SessionToken,
		}, nil
	})

	return creds, nil
}

func listDZIFiles(ctx context.Context, bucketName, prefix string, brconfig brconfig.Configuration, logger logging.Logger) ([]string, error) {
	token := brconfig.GetStringDefault("storage:aws:token", "sample")
	region := brconfig.GetStringDefault("storage:aws:region", "eu-west-1")
	creds, err := getCredentialsFromCognito(ctx, token, brconfig, logger)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region), config.WithCredentialsProvider(creds))
	if err != nil {
		return nil, err
	}

	s3Svc := s3.NewFromConfig(cfg)
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}

	var dziFiles []string
	paginator := s3.NewListObjectsV2Paginator(s3Svc, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, content := range page.Contents {
			key := aws.ToString(content.Key)
			if strings.HasSuffix(key, ".dzi") {
				dziFiles = append(dziFiles, key)
			}
		}
	}

	return dziFiles, nil
}

func saveToPostgres(dziFiles []string, brconfig brconfig.Configuration, logger logging.Logger) error {
	//conn := h.Configuration.GetString("authorization:recaptchasecret")
	driver := brconfig.GetStringDefault("sql:driverName", "postgres")
	connectionUrl, found := brconfig.GetString("sql:connectionUrl")

	var connectionStr string
	var err error
	if !found {
		logger.Panic("openDB: Cannot read SQL connection string from config")
	} else {
		logger.Debug("openDB: found SQL connection string from config")
	}
	connectionStr, err = pq.ParseURL(connectionUrl)
	if err != nil {
		logger.Panic("openDB: Error converting SQL URL connection from config to connection string")
	}
	logger.Debugf("openDB: Connection string: ", connectionStr)

	// connStr := "host=localhost port=5432 user=postgres password=postgres dbname=brucheion sslmode=disable"
	db, err := sql.Open(driver, connectionStr)
	if err != nil {
		return err
	}
	defer db.Close()

	// In this modified version of the saveToPostgres function,
	// we added the ON CONFLICT (filepath) DO NOTHING clause to the SQL statement.
	// This means that if a record with the same filepath already exists in the imagecollectionlist table,
	// the database will not perform any action, effectively merging the dziFiles array with the existing records in the table.
	// If you want to update other columns in the existing record,
	// you can replace DO NOTHING with DO UPDATE SET column_name = value in the ON CONFLICT clause.
	// With this change, the saveToPostgres function will merge the DZI files with the imagecollectionlist table,
	// avoiding duplicate records.
	for _, dziFile := range dziFiles {

		_, err := db.Exec(`
			INSERT INTO imagecollectionlist (filepath)
			VALUES ($1)
			ON CONFLICT (filepath) DO NOTHING
		`, dziFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func SyncDZIcollection(brconfig brconfig.Configuration, logger logging.Logger) {

	bucketName := brconfig.GetStringDefault("storage:aws:bucketname", "brucheion")
	prefix := brconfig.GetStringDefault("storage:aws:prefix", "")
	//prefix := "nbh/J1img/positive/"

	dziFiles, err := listDZIFiles(context.Background(), bucketName, prefix, brconfig, logger)
	if err != nil {
		logger.Infof("Error listing DZI files: %v", err)
	}

	err = saveToPostgres(dziFiles, brconfig, logger)
	if err != nil {
		logger.Infof("Error saving DZI files to PostgreSQL: %v", err)
	}

	// err = saveToHstorePostgres(dziFiles)
	// if err != nil {
	// 	log.Fatalf("Error saving DZI files to PostgreSQL: %v", err)
	// }
	logger.Info("Successfully saved DZI files to PostgreSQL")
}
