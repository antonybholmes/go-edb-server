curl -v -X POST -H 'Content-Type: application/json' -d  '{"locations":[{"chr":"chr1","start":10,"end":10}], "level":"gene", "tss":[2000, 1000]}' "localhost:8080/annotation/grch38?n=5"

curl -v -X POST -H 'Content-Type: application/json' -d  '{"locations":[{"chr":"chr1","start":10,"end":10}], "level":"gene", "tss":[2000, 1000]}' "localhost:8080/annotation/grch38?n=5&output=text"

curl -v "localhost:8080/dna/grch38"

curl -v -X POST -H 'Content-Type: application/json' -d  '{"locations":[{"chr":"chr10","start":1043441,"end":1044114}], "level":"gene", "tss":[2000, 1000]}' "localhost:8080/v1/annotate/grch38?n=5&output=text"
curl -v -X POST -H 'Content-Type: application/json' -d  '{"locations":[{"chr":"chr10","start":1043441,"end":1044114},{"chr":"chr10","start":104349828,"end":104350217}], "level":"gene", "tss":[2000, 1000]}' "localhost:8080/v1/annotate/grch38?n=5&output=text"




curl -v -X POST -H 'Content-Type: application/json' -d  '{"locations":[{"chr":"chr10","start":100014303,"end":100014664}], "level":"gene", "tss":[2000, 1000]}' "https://api.rdf-lab.org/v1/annotate/grch38?n=5&output=text"



