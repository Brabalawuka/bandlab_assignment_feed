package config

import (
	"bandlab_feed_server/util"
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/spf13/viper"
)

type Config struct {
	MongoURL      string
	MongoDatabase string

	R2AccessKey                 string
	R2SecretKey                 string
	R2AccountId                 string
	R2BucketName                string
	R2ImagePresignExpirationSec int
	R2PublicBucketURL             string
	OriginalImagePath           string
	ProcessedImagePath          string
	PostRecentCommentsCount     int
}

// Global AppConfig
var AppConfig *Config

func Init() {
	LoadConfig()
}

// LoadConfig loads the configuration from the environment variables.
func LoadConfig() {

	viper.SetConfigFile("config/config.yml")               // 指定配置文件
	err := viper.ReadInConfig()                            // 读取配置信息
	viper.AutomaticEnv()                                   // 读取匹配的环境变量
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 读取环境变量的分隔符
	if err != nil {                                        // 读取配置信息失败
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	AppConfig = &Config{
		MongoURL:                    viper.GetString("mongo.url"),
		MongoDatabase:               viper.GetString("mongo.initdb_database"),
		R2AccessKey:                 viper.GetString("r2.access_key"),
		R2SecretKey:                 viper.GetString("r2.secret_key"),
		R2AccountId:                 viper.GetString("r2.account_id"),
		R2BucketName:                viper.GetString("r2.bucket_name"),
		R2ImagePresignExpirationSec: viper.GetInt("r2.presign_expiration_sec"),
		R2PublicBucketURL:           viper.GetString("r2.public_bucket_url"),
		OriginalImagePath:         viper.GetString("image.original_image_path"),
		ProcessedImagePath:        viper.GetString("image.processed_image_path"),
		PostRecentCommentsCount:   viper.GetInt("post.recent_comments_count"),
	}

	hlog.Info("AppConfig: %s", util.ToJsonString(AppConfig))
	hlog.Info("Config loaded successfully")
}
