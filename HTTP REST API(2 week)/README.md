Написал пробную версию. Буду потихонько ее дорабатывать

все задачи:
```cmd
curl http://localhost:8080/tasks/
```

добавить задачу:
```cmd
curl -X POST http://localhost:8081/tasks/ \
  -H "Content-Type: application/json" \
  -d '{"text": "qwe", "done": false}'
```

метка что задача выполнена:
```cmd
curl -X PUT http://localhost:8081/tasks/1/done
```