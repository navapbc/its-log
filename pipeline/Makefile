etl:
	jsonnet pipeline.jsonnet > pipeline.json
	go run ../... etl --runscript pipeline.json --sqlite 2025-12-24.sqlite