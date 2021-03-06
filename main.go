package main

import (
	"context"
	"crawler_ig_content/instagram_scraper"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
)

type Option struct {
	Address  string
	Password string
	Stage    string
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var ctx = context.Background()

func main() {
	scheduler := flag.Bool("scheduler", false, "Set scheduler for running or manual (ex: ./apps -scheduler=true)")
	time := flag.String("time", "", "Set time")
	interval := flag.Int("interval", 0, "Set time")
	flag.Parse()

	timeInterval := uint64(*interval)
	if *scheduler == true {
		switch *time {
		case "week":
			gocron.Every(timeInterval).Weeks().Do(task)
		case "day":
			gocron.Every(timeInterval).Days().Do(task)
		case "hour":
			gocron.Every(timeInterval).Hours().Do(task)
		case "minute":
			gocron.Every(timeInterval).Minutes().Do(task)
		case "second":
			gocron.Every(timeInterval).Seconds().Do(task)
		default:
			gocron.Every(1).Days().Do(task)
		}
		<-gocron.Start()
	} else {
		task()
	}
}

func task() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//File Config
	fileConfig := flag.String("config", "", "Your config file (default in .env file)")
	//if usage file config in parameter config
	if *fileConfig == "" {
		*fileConfig = os.Getenv("FILE_CONFIG")
	}

	//Redis Config
	redisConfig := flag.String("redis", "", "Your redis config (default in .env file). ex: 'address|password|port|prefix'")

	option := &Option{
		Address:  fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		Stage:    os.Getenv("APP_LEVEL"),
	}

	if *redisConfig != "" {
		opt := strings.Split(*redisConfig, "|")
		option.Address = fmt.Sprintf("%s:%s", opt[0], opt[2])
		option.Password = opt[1]
		option.Stage = opt[3]
	}

	flag.Parse()

	client := connectRedis(option)

	configs, err := readConfig(*fileConfig)

	if err != nil {
		panic(err)
	}

	for v := range *configs {

		var (
			response interface{}
		)
		module := (*configs)[v].(map[string]interface{})

		prefix := fmt.Sprintf("%s_%s", module["prefix"].(string), option.Stage)
		limit := int(module["limit"].(float64))

		keyName := fmt.Sprintf("%s_%s_%d", prefix, "get_instagram_limit", limit)

		data, statusCode, err := instagram_scraper.FetchMediaImage(module["username"].(string), limit)
		if err != nil {
			panic(err)
		}

		if *data != nil && statusCode != 429 {
			fmt.Println(keyName)
			client.Del(ctx, keyName)
			response = data
			raw, err := json.Marshal(&Response{Status: http.StatusOK, Message: fmt.Sprintf("Fetch Instagram Media %s", module["username"]), Data: response})
			if err != nil {
				panic(err)
			}

			fmt.Println(string(raw))
			err = client.Set(ctx, keyName, raw, 0).Err()
			if err != nil {
				panic(err)
			}

		} else {
			format := fmt.Sprintf("[%s] => STATUS CODE %d", v, statusCode)
			fmt.Println(format)
		}

	}

}

func readConfig(cfg string) (*map[string]interface{}, error) {
	var jsonMap map[string]interface{}

	file, err := os.Open(cfg)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	val, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(val, &jsonMap)

	if err != nil {
		return nil, err
	}

	return &jsonMap, nil

}

func connectRedis(opt *Option) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     opt.Address,
		Password: opt.Password,
		DB:       0,
	})
	return rdb
}
