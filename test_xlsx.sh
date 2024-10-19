curl -H "Content-Type: application/json" -X POST --data "{\"b64xlsx\":\"$(base64 -w 0 test.xlsx)\",\"indexCol\":0,\"header\":0}" localhost:8080/util/convert/xlsx

curl -H "Content-Type: application/json" -X POST --data "{\"b64xlsx\":\"$(base64 -w 0 test.xlsx)\",\"indexCol\":-1,\"header\":-1}" localhost:8080/util/convert/xlsx

curl -H "Content-Type: application/json" -X POST --data "{\"b64xlsx\":\"$(base64 -w 0 test.xlsx)\",\"indexCol\":0,\"header\":-1}" localhost:8080/util/convert/xlsx

curl -H "Content-Type: application/json" -X POST --data "{\"b64xlsx\":\"$(base64 -w 0 test.xlsx)\",\"indexCol\":-1,\"header\":0}" localhost:8080/util/convert/xlsx

curl -H "Content-Type: application/json" -X POST --data "{\"b64xlsx\":\"$(base64 -w 0 test.xlsx)\",\"indexCol\":1,\"header\":1}" localhost:8080/util/convert/xlsx
