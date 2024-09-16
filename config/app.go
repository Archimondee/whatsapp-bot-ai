package config

import "github.com/spf13/viper"

type Config struct {
	VerifyToken                    string `mapstructure:"VERIFY_TOKEN"`
	WhatsappToken                  string `mapstructure:"WHATSAPP_TOKEN"`
	HTTPServerAddress              string `mapstructure:"HTTP_SERVER_ADDRESS"`
	OpenAIApiKey                   string `mapstructure:"OPENAI_API_KEY"`
	OpenAIGptModelName             string `mapstructure:"OPENAI_GPT_MODEL_NAME"`
	OpenAIGptModelToken            string `mapstructure:"OPENAI_GPT_MODEL_TOKEN"`
	OpenAIGptModelTemperature      string `mapstructure:"OPENAI_GPT_MODEL_TEMPERATURE"`
	OpenAIGptModelTopP             string `mapstructure:"OPENAI_GPT_MODEL_TOP_P"`
	OpenAIGptModelPenaltyPresence  string `mapstructure:"OPENAI_GPT_MODEL_PENALTY_PRESENCE"`
	OpenAIGptModelPenaltyFrequency string `mapstructure:"OPENAI_GPT_MODEL_PENALTY_FREQUENCY"`

	LeonardoApiKey    string `mapstructure:"LEONARDO_API_KEY"`
	DreamshaperV6     string `mapstructure:"DREAMSHAPER_V6"`
	LeonardoCreative  string `mapstructure:"LEONARDO_CREATIVE"`
	LeonardoSelect    string `mapstructure:"LEONARDO_SELECT"`
	LeonardoSignature string `mapstructure:"LEONARDO_SIGNATURE"`

	DreamshaperTag string `mapstructure:"DREAMSHAPER_TAG"`
	CreativeTag    string `mapstructure:"CREATIVE_TAG"`
	SelectTag      string `mapstructure:"SELECT_TAG"`
	SignatureTag   string `mapstructure:"SIGNATURE_TAG"`

	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`

	WhatsappDeleteConversation string `mapstructure:"WHATSAPP_DELETE_CONVERSATION"`
	WhatsappNewConversation    string `mapstructure:"WHATSAPP_NEW_CONVERSATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
