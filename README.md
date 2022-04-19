# POC SQS em Golang
 - Produzindo e consumindo mensagens do SQS

## Configurações

 - Criar as variaveis de ambiente:
``` 
export AWS_ACCESS_KEY=""
export AWS_SECRET_ACCESS_KEY=""
export AWS_URL_QUEUE=""
export AWS_REGION=""
```

 - Run:

``` 
go run src/main.go 
```

## API

- Produção de mensagens:

``` 
Url: http://localhost:3000/producer
Method: POST
Content-Type: application/json

Body:{
  "id": "123",
  "message": "test 123"
}

Response:{
"statusCod": 201,
"message": "message send successfully, id: %!(EXTRA string=id message)"
}
```