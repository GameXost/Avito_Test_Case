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
- потестировал нагрузкой и простенький линтер поставил 


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


P.P.S
Еще захотел поиграть с тестированием
В докере ничего не падает, сервис продолжает отвечать
```Bash
GameXost@моржовый MINGW64 /d/apps/go_files/new_go_vue/Avito_Test_Case (main)                                                                                                          
$ go-wrk -c 128 -d 300 "http://localhost:8080/team/get?team_name=wow"
Running 300s test @ http://localhost:8080/team/get?team_name=wow
  128 goroutine(s) running concurrently
71802 requests in 15.164266082s, 19.65MB read
Requests/sec:           4734.95
Transfer/sec:           1.30MB
Overall Requests/sec:   239.07
Overall Transfer/sec:   67.00KB
Fastest Request:        1.058ms
Avg Req Time:           27.032ms
Slowest Request:        1.046911s
Number of Errors:       38121
Error Counts:           net/http: timeout awaiting response headers=36025,No connection could be made because the target machine actively refused it.=2096
10%:                    2.239ms
50%:                    3.228ms
75%:                    3.405ms
99%:                    3.686ms
99.9%:                  3.696ms
99.9999%:               3.696ms
99.99999%:              3.696ms
stddev:                 95.966ms
```

```Bash
rror Counts:           net/http: timeout awaiting response headers=36025,No connection could be made because the target machine actively refused it.=2096
10%:                    2.239ms
50%:                    3.228ms
75%:                    3.405ms
99%:                    3.686ms
99.9%:                  3.696ms
99.9999%:               3.696ms
99.99999%:              3.696ms
stddev:                 95.966ms

GameXost@моржовый MINGW64 /d/apps/go_files/new_go_vue/Avito_Test_Case (main)                                                                                                          
$ curl "http://localhost:8080/team/get?team_name=wow"
{"team_name":"wow","members":[{"user_id":"1","username":"GameXost","is_active":true},{"user_id":"2","username":"PuPu","is_active":true},{"user_id":"3","username":"Karabaxzzz","is_active":true}]}

```



P.P.P.S. 20.11 Решил еще чуть потестить , уже с пом. Vegeta. Итог: держит 2к\с запросов около 10 секунд с 100%, больше запросов или больше времени - ГГ, проседает
```Bash
$ echo "GET http://localhost:8080/team/get?team_name=wow" | vegeta attack -rate=2000 -duration=10s | vegeta report
Requests      [total, rate, throughput]  20000, 2000.15, 1999.85
Duration      [total, attack, wait]      10.0007527s, 9.9992362s, 1.5165ms
Latencies     [mean, 50, 95, 99, max]    1.958228ms, 1.072013ms, 4.026074ms, 19.088386ms, 122.4616ms
Bytes In      [total, mean]              4920000, 246.00
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:20000
Error Set:



$ echo "GET http://localhost:8080/team/get?team_name=wow" | vegeta attack -rate=2000 -duration=30s | vegeta report
Requests      [total, rate, throughput]  59999, 2000.00, 452.07
Duration      [total, attack, wait]      59.9998851s, 29.9995715s, 30.0003136s
Latencies     [mean, 50, 95, 99, max]    2.96181755s, 17.037288ms, 24.769279513s, 29.129629685s, 30.0012119s
Bytes In      [total, mean]              6672504, 111.21
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    45.21%
Status Codes  [code:count]               0:32875  200:27124
Error Set:
Get "http://localhost:8080/team/get?team_name=wow": dial tcp 0.0.0.0:0->[::1]:8080: connectex: No connection could be made because the target machine actively refused it.
Get "http://localhost:8080/team/get?team_name=wow": EOF
Get "http://localhost:8080/team/get?team_name=wow": context deadline exceeded (Client.Timeout exceeded while awaiting headers)

$ echo "GET http://localhost:8080/team/get?team_name=wow" | vegeta attack -rate=1000 -duration=30s | vegeta report
Requests      [total, rate, throughput]  30000, 1000.04, 1000.01
Duration      [total, attack, wait]      29.9997855s, 29.9987587s, 1.0268ms
Latencies     [mean, 50, 95, 99, max]    870.639µs, 1.024834ms, 1.553076ms, 2.570445ms, 23.4776ms
Bytes In      [total, mean]              7380000, 246.00
Bytes Out     [total, mean]              0, 0.00
Success       [ratio]                    100.00%
Status Codes  [code:count]               200:30000
Error Set:


```
Скорее всего нужно пересмотреть подключение к БД, пересмотреть работу с http и ограничениея с таймаутами, переработать максимальное кол-во пулов к БД и запускать на более мощной машине :)))
