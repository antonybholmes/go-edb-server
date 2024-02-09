curl -v -X POST -H 'Content-Type: application/json' -d  '{"locations":[{"chr":"chr1","start":10,"end":10}], "level":"gene", "tss":[2000, 1000]}' "localhost:8080/annotation/grch38?n=5"

curl -v -X POST -H 'Content-Type: application/json' -d  '{"locations":[{"chr":"chr1","start":10,"end":10}], "level":"gene", "tss":[2000, 1000]}' "localhost:8080/annotation/grch38?n=5&output=text"

curl -v "localhost:8080/dna/grch38"

curl -v "localhost:8080/dna/grch38?start=1&end=100&format=lower&mask=n"