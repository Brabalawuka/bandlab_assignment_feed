package dal

import (
	"bandlab_feed_server/config"
	"bandlab_feed_server/dal/cloudflare"
	"bandlab_feed_server/dal/mongodb"
)




func InitDal(config *config.Config) {
	err := cloudflare.Initialize(&cloudflare.Config{
		AccessKey:  config.R2AccessKey,
		SecretKey:  config.R2SecretKey,
		AccountID:  config.R2AccountID,
		BucketName: config.R2BucketName,
	})
	if err != nil {
		panic(err)
	}
	err = mongodb.Initialize(&mongodb.Config{
		URL:      config.MongoURL,
		Database: config.MongoDatabase,
	})
	if err != nil {
		panic(err)
	}
}