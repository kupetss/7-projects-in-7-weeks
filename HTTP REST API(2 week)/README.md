## День 1:
Написал пробную версию. Буду потихонько ее дорабатывать

## День 2:
добавил удаление

## День 3:
сделал все по тз. Добавил логирование

# Инструкция
все задачи:
```cmd
curl http://localhost:8080/tasks/
```

добавить задачу:
```cmd
curl -X POST http://localhost:8080/tasks/ \
  -H "Content-Type: application/json" \
  -d '{"text": "qwe", "done": false}'
```

метка что задача выполнена:
```cmd
curl -X PATCH http://localhost:8080/task/1
```

удалить задачу:
```cmd
curl -X DELETE http://localhost:8080/tasks/2 -v
```