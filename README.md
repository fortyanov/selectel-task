### Задача
Необходимо реализовать микросервис который, используя публичное API Vscale,
может создавать и удалять группу серверов.

### Техническое задание
API должен реализовывать два метода метода: POST для создания указанного
количества серверов (передается параметром), DELETE для удаления всех созданных
ранее.

### Условия
- Метод создания должен работать по принципу "все или ничего": если создание хотя
бы одного сервера завершилось ошибкой, необходимо удалить все уже созданные в
эту операцию.
- Исходящие запросы в публичное API необходимо распараллеливать, чтобы
минимизировать время ответа разрабатываемого микросервиса.
При этом стоит учесть что API может вернуть ошибку 429 Too many requests в случае,
если обращаться слишком часто, в таком случае не считать это ошибкой, а
попытаться повторить позже.
- Приложение должно быть написано на Go и подготовлено к сборке: присутствует
Dockerfile и/или Makefile.
- Приложения должны сопровождаться минимальной документацией по сборке, запуску
и использованию.
- Исходный код приложений необходимо предоставить в виде ссылки на github /
  bitbucket.
  
# Развертывание
1. Установить docker и docker-compose
2. В файле docker-compose.yml прописать XTOKEN
3. Перейти в консоли в корневую директорию проекта
4. Выполнить в консоли `docker-compose up -d`

Или

1. Перейти в консоли в корневую директорию проекта
2. Выполнить в консоли `docker build -t selectel-task .`
3. Выполнить в консоли `docker run --name selectel-task -p 8080:8080 -d selectel-task `

### Примеры запросов
Создание 5 скалет
```
curl 127.0.0.1:8080/scalets -X POST -d '{"count": 5}'

```
Удаление всех скалет пользователя
```
curl 127.0.0.1:8080/scalets -X DELETE
```
