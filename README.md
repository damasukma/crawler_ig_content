# crawler_ig_content
For Service Crawler instagram set in redis


for file .env
* APP_LEVEL=development
* REDIS_HOST=127.0.0.1
* REDIS_PASSWORD=
* REDIS_PORT=6379
* FILE_CONFIG=(path file config)


FOR CONFIG_FILE
```json
  {
  "config1": {
    "username": "jhon",
    "limit": 5,
    "prefix": "jhon_" //prefix for redis
  },
  "config_2":{
    "username": "doe",
    "limit": 5,
    "prefix": "doe_" //prefix for redis
  }
}
  
```
