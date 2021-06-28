# Btc-Requester

Примеры работы:

### Регистрация:

```bash
curl "http://localhost:9990/user/create?email=example0@mail.org&password=1qaz_"
```
```javascript
[200] {"status":"Ok"}
```

### Повторная регистрация:

```bash
curl "http://localhost:9990/user/create?email=example0@mail.org&password=1qaz_"
```
```javascript
{"status":"Email already used"}
```

### Логин:

```bash
curl "http://localhost:9990/user/login?email=example0@mail.org&password=1qaz_"
```
```javascript
[200] {"status":"Ok","token":"kuap0bzu0qd1"}
```

### Запрос цены:

```bash
curl "http://localhost:9990/btcRate?token=kuap0bzu0qd1"
```
```bash
curl "http://localhost:9990/btcRate" -H "X-API-Key: kuap0bzu0qd1"
```
```javascript
[200] {"status":"Ok","token":"kuap0bzu0qd1"}
```
### Неудачные запросы цены:
```bash
curl "http://localhost:9990/btcRate?token="
```
```javascript
[400] {"status":"Missing token"}
```
```bash
curl "http://localhost:9990/btcRate?token=WRONG_TOKEN"
```
```javascript
[403] {"status":"Invalid token"}
```
