# Сервис подсчета арифметических выражений 
Выполнен в рамках курса «Программирование на Go | 24» 

# Описание 
Веб-сервис, который принимает от пользователя POST-запрос с арифметическим выражением и возвращает значение арифметического выражения или ошибку, если выражение некорректно 

# Запуск
(Требования: go, git)</br>
Выполните команды:</br>
```
git clone git@github.com:Starostina-elena/yalms_go_task2.git
```
```
cd yalms_go_task2
```
```
go mod init github.com/Starostina-elena/yalms_go_task2
```
```
go run cmd/main.go
```

Запущенный сервер будет доступен по адресу 
`http://localhost:8080/api/v1/calculate`

# Примеры использования 
```
curl --location "http://localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"2+2\"}"
{"result":"4.000000"}
```
```
curl --location "http://localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"2+\"}"
{"error":"Expression is not valid"}
```

# Запуск тестов

Выполните команды:</br>
```
cd internal/application
```
```
go test
```