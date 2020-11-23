package main

import (
	"context"
	"crawler_ig_content/instagram_scraper"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Option struct{
	Address string
	Password string
	Stage string
}


type Response struct{
	Status int
	Message string
	Data interface{}
}


var ctx = context.Background()

func main(){
	scheduller := flag.Bool("scheduller", false, "Set scheduller for running or manual (ex: ./apps -scheduller=true)")
	flag.Parse()
	if *scheduller == true{
		gocron.Every(1).Days().Do(task)
		<-gocron.Start()
	}else{
		task()
		fmt.Println("ULL")
	}
}


func task(){

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
		Address: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("PASSWORD"),
		Stage: os.Getenv("APP_LEVEL"),
	}

	if *redisConfig != ""{
		opt :=  strings.Split(*redisConfig, "|")
		option.Address = fmt.Sprintf("%s:%s", opt[0], opt[2])
		option.Password = opt[1]
		option.Stage = opt[3]
	}

	flag.Parse()

	client := connectRedis(option)

	configs, err := readConfig(*fileConfig)


	if err != nil{
		panic(err)
	}


	for v := range *configs{
		fmt.Println(v)
		var (
			response interface{}
		)
		module := (*configs)[v].(map[string]interface{})

		prefix := fmt.Sprintf("%s_%s", module["prefix"].(string), option.Stage)
		limit := int(module["limit"].(float64))

		keyName := fmt.Sprintf("%s_%s_%d",prefix,"get_instagram_limit", limit)

		data, err := instagram_scraper.FetchMediaImage(module["username"].(string), limit)
		if err != nil {
			panic(err)
		}


		if data != nil{
			client.Del(ctx, keyName)
			response = data
			raw, _ := json.Marshal(&Response{Status: http.StatusOK, Message: fmt.Sprintf("Fetch Instagram Media %s", module["username"]), Data: response})
			client.Set(ctx, keyName, raw, 0)

		}

	}

}

func readConfig(cfg string) (*map[string]interface{}, error){
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

func connectRedis(opt *Option) *redis.Client{
	rdb := redis.NewClient(&redis.Options{
		Addr: opt.Address,
		Password: opt.Password,
		DB: 0,
	})
	return rdb
}