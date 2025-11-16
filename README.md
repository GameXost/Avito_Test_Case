# AVITO TEST_CASE


## ЗАПУСК

### Через докер композ

Из корневой директории:
```
docker-compose up --build
```
Работает на `http://localhost:8080`
Миграции сработают из `./db/db.sql`
ВСЕ нужные параметры заданы в docker-compose 

P.S.
- ошибки, не прописанные в документации, отправляю с статус кодом 500 или 400(для отсутствующих данных в теле запроса или некорректной обработке json)



### Нагрузочное тестирование использую wrg-go
1. создание команды
```Bash
curl -X POST http://localhost:8080/team/add -H "Content-Type: application/json" -d '{
  "team_name": "wow",
  "members": [
    {"user_id": "1", "username": "GameXost", "is_active": true},
    {"user_id": "2", "username": "PuPu", "is_active": true},
    {"user_id": "3", "username": "Karabaxzzz", "is_active": true}
  ]
}'
```
2. Чекнул создалась ли команда: `curl "http://localhost:8080/team/get?team_name=wow"`
```Bash
curl "http://localhost:8080/team/get?team_name=wow"
```
3. Запускаем тесты на GetTeam  и идем смотреть, как всему приходит ГГ:))))
```Bash
go-wrk -c 128 -d 30 "http://localhost:8080/team/get?team_name=wow"
```
4. Ловим вывод:
Как будто всё сломалось
```Bash
$ go-wrk -c 128 -d 30 "http://localhost:8080/team/get?team_name=wow"
Running 30s test @ http://localhost:8080/team/get?team_name=wow
  128 goroutine(s) running concurrently
26522 requests in 2.47086742s, 7.26MB read
Requests/sec:           10733.88
Transfer/sec:           2.94MB
Overall Requests/sec:   865.08
Overall Transfer/sec:   242.46KB
Fastest Request:        1.082ms
Avg Req Time:           11.924ms
Slowest Request:        201.167ms
Number of Errors:       3584
Error Counts:           net/http: timeout awaiting response headers=3584
10%:                    2.112ms
50%:                    2.715ms
75%:                    3.07ms
99%:                    3.204ms
99.9%:                  3.205ms
99.9999%:               3.205ms
99.99999%:              3.205ms
stddev:                 16.024ms

```
Ну кстати, сервис в итоге устоял, не крашнулся, ластовые запросы, видимо долго обрабатывался и wrk чет не понрав:
```Bash
app-1  | 2025/11/16 20:33:13 [2ee9634c35fd/NmRPoMDv82-030105] "GET http://localhost:8080/team/get?team_name=wow HTTP/1.1" from 172.20.0.1:56264 - 200 195B in 515.064µs               
app-1  | 2025/11/16 20:33:13 [2ee9634c35fd/NmRPoMDv82-030106] "GET http://localhost:8080/team/get?team_name=wow HTTP/1.1" from 172.20.0.1:56278 - 200 195B in 405.148µs               
```




### линтер, запушил в репо
```yml
linters:
  enable:
    - govet
    - gofmt
    - errcheck
    - staticcheck

run:
  timeout: 2m
```