curl -H "Content-Type: application/json" -X POST --data "{\"b64xlsx\":\"$(base64 -w 0 test.xlsx)\"}" localhost:8080/utils/xlsx/sheets
curl -H "Content-Type: application/json" -X POST --data "{\"b64xlsx\":\"$(base64 -w 0 test.xlsx)\",\"indexes\":0,\"headers\":0}" localhost:8080/utils/xlsx/convert
curl -H "Content-Type: application/json" -X POST --data "{\"b64xlsx\":\"$(base64 -w 0 test.xlsx)\",\"indexes\":1,\"headers\":1,\"skipRows\":1}" localhost:8080/utils/xlsx/convert
#curl -H "Content-Type: application/json" -X POST --data "{\"b64xlsx\":\"$(base64 -w 0 test.xlsx)\",\"indexes\":0,\"headers\":-1}" localhost:8080/utils/xlsx/convert
#curl -H "Content-Type: application/json" -X POST --data "{\"b64xlsx\":\"$(base64 -w 0 test.xlsx)\",\"indexes\":-1,\"headers\":0}" localhost:8080/utils/xlsx/convert
#curl -H "Content-Type: application/json" -X POST --data "{\"b64xlsx\":\"$(base64 -w 0 test.xlsx)\",\"indexes\":1,\"headers\":1}" localhost:8080/utils/xlsx/convert
