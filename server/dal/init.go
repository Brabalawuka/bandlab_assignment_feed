package dal

import (
	"bandlab_feed_server/config"
	"bandlab_feed_server/dal/cloudflare"
)




func InitDal(config *config.Config) {
	cloudflare.Initialize(&cloudflare.Config{
		AccessKey:  config.R2AccessKey,
		SecretKey:  config.R2SecretKey,
		AccountID:  config.R2AccountID,
		BucketName: config.R2BucketName,
	})
}