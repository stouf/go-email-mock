# Go email mock

## Motivation

It is nowadays very common to use an external service that takes care of sending
emails on bahalf of your application.
As an effort of making testing as reliable as possible, this repository
provides a very simple HTTP server that lets you simulate sending emails to a
third party service. It also provides an endpoint that will let you access the
information it has collected.

## How to run?

```bash
make build
HTTP_PORT=1234 make start
# ...
make stop
```

## HTTP API

### Send an email

*Request*

```json
POST /send
Content-Type: application/json

{
  "sender": "me@test.com",
  "senderName": "Me Test",
  "recipients": ["foo@test.com", "bar@test.com"],
  "carbonCopy": ["foo@test.com", "bar@test.com"],
  "blindCarbonCopy": ["foo@test.com", "bar@test.com"],
  "subject": "User registration",
  "content": "...."
}
```

The mandotory fields of the payload above are:

- `sender`
- `recipients`
- `subject`
- `content`

All the other fields are optional.

*Response*

```
200 OK
```


### List emails

*Request*

```
POST /get
Content-Type: application/json
```

*Response*

```json
200 OK

{
  "emails": [
    {
      "sender": "me@test.com",
      "senderName": "Me Test",
      "recipients": ["foo@test.com", "bar@test.com"],
      "carbonCopy": ["foo@test.com", "bar@test.com"],
      "blindCarbonCopy": ["foo@test.com", "bar@test.com"],
      "subject": "User registration",
      "content": "...."
    }
  ]
}
```

The following properties of an email can be empty:

- `senderName` (empty string)
- `carbonCopy` (empty array)
- `blindCarbonCopy` (empty array)


### Flush all data

*Request*

```
POST /flush
```

*Response*

```
200 OK
```


### Tests

```bash
HTTP_PORT=1234 make test
```
