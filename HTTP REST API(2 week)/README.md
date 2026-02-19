Написал пробную версию. Буду потихонько ее дорабатывать

добавил удаление

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
curl -X PUT http://localhost:8080/tasks/1/done
```

удалить задачу:
```cmd
curl -X DELETE http://localhost:8080/tasks/2 -v
```