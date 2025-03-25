package config

type Config struct {
	Port string
}

var AppConfig Config

func InitConfig() {
	AppConfig = Config{
		Port: ":8080", // 默认端口，可以从环境变量或配置文件加载
	}
}