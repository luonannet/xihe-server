package config

import (
	"errors"
	"regexp"
	"strings"

	"github.com/huaweicloud/golangsdk"
	"github.com/opensourceways/community-robot-lib/mq"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

var reIpPort = regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}:[1-9][0-9]*$`)

func LoadConfig(path string) (*Config, error) {
	v := new(Config)

	if err := utils.LoadFromYaml(path, v); err != nil {
		return nil, err
	}

	v.setDefault()

	if err := v.validate(); err != nil {
		return nil, err
	}

	return v, nil
}

type Config struct {
	DefaultPassword string `json:"default_password" required:"true"`
	EncryptionKey   string `json:"encryption_key" required:"true"`
	MaxRetry        int    `json:"max_retry"`
	ActivityKeepNum int    `json:"activity_keep_num"`

	Authing  AuthingService `json:"authing_service" required:"true"`
	Resource Resource       `json:"resource" required:"true"`
	Mongodb  Mongodb        `json:"mongodb" required:"true"`
	Gitlab   Gitlab         `json:"gitlab" required:"true"`
	API      API            `json:"api" required:"true"`
	User     User           `json:"user"`
	MQ       MQ             `json:"mq" required:"true"`
}

func (cfg *Config) GetMQConfig() mq.MQConfig {
	return mq.MQConfig{
		Addresses: cfg.MQ.parseAddress(),
	}
}

func (cfg *Config) setDefault() {
	cfg.Resource.setdefault()
	cfg.User.setDefault()

	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 10
	}

	if cfg.ActivityKeepNum <= 0 {
		cfg.ActivityKeepNum = 25
	}
}

func (cfg *Config) validate() error {
	if _, err := golangsdk.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	if err := cfg.Resource.validate(); err != nil {
		return err
	}

	if _, err := domain.NewPassword(cfg.DefaultPassword); err != nil {
		return err
	}

	return cfg.MQ.validate()
}

type Mongodb struct {
	MongodbConn        string `json:"mongodb_conn" required:"true"`
	DBName             string `json:"mongodb_db" required:"true"`
	ProjectCollection  string `json:"project_collection" required:"true"`
	ModelCollection    string `json:"model_collection" required:"true"`
	DatasetCollection  string `json:"dataset_collection" required:"true"`
	UserCollection     string `json:"user_collection" required:"true"`
	LoginCollection    string `json:"login_collection" required:"true"`
	LikeCollection     string `json:"like_collection" required:"true"`
	ActivityCollection string `json:"activity_collection" required:"true"`
}

type AuthingService struct {
	APPId       string `json:"app_id" required:"true"`
	Secret      string `json:"secret" required:"true"`
	RedirectURI string `json:"redirect_uri" required:"true"`
}

type Gitlab struct {
	Endpoint  string `json:"endpoint" required:"true"`
	RootToken string `json:"root_token" required:"true"`
}

type API struct {
	APITokenExpiry int64  `json:"api_token_expiry" required:"true"`
	APITokenKey    string `json:"api_token_key" required:"true"`
}

type Resource struct {
	MaxNameLength int `json:"max_name_length"`
	MinNameLength int `json:"min_name_length"`
	MaxDescLength int `json:"max_desc_length"`

	Covers           []string `json:"covers" required:"true"`
	Protocols        []string `json:"protocols" required:"true"`
	ProjectType      []string `json:"project_type" required:"true"`
	TrainingPlatform []string `json:"training_platform" required:"true"`
}

func (r *Resource) setdefault() {
	if r.MaxNameLength == 0 {
		r.MaxNameLength = 50
	}

	if r.MinNameLength == 0 {
		r.MinNameLength = 5
	}

	if r.MaxDescLength == 0 {
		r.MaxDescLength = 100
	}
}

func (r *Resource) validate() error {
	if r.MaxNameLength < (r.MinNameLength + 10) {
		return errors.New("invalid name length")
	}

	return nil
}

type User struct {
	MaxNicknameLength int `json:"max_nickname_length"`
	MaxBioLength      int `json:"max_bio_length"`
}

func (u *User) setDefault() {
	if u.MaxNicknameLength == 0 {
		u.MaxNicknameLength = 20
	}

	if u.MaxBioLength == 0 {
		u.MaxBioLength = 200
	}
}

type MQ struct {
	Address        string `json:"address" required:"true"`
	TopicLike      string `json:"topic_like" required:"true"`
	TopicFollowing string `json:"topic_following" required:"true"`
}

func (cfg MQ) validate() error {
	if r := cfg.parseAddress(); len(r) == 0 {
		return errors.New("invalid mq address")
	}

	return nil
}

func (cfg MQ) parseAddress() []string {
	v := strings.Split(cfg.Address, ",")
	r := make([]string, 0, len(v))
	for i := range v {
		if reIpPort.MatchString(v[i]) {
			r = append(r, v[i])
		}
	}

	return r
}
