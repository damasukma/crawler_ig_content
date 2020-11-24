# crawler_ig_content
For Service Crawler instagram set in redis


FILE .env
---------

* APP_LEVEL=development
* REDIS_HOST=127.0.0.1
* REDIS_PASSWORD=
* REDIS_PORT=6379
* FILE_CONFIG=(path file config)


FILE_CONFIG
-----------

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


RUNNING WITH SCHEDULER
-------

```shell
  
  go run main.go -scheduler=true -time=day -interval=1

```

