jsonnet:
	jsonnet run.jsonnet > run.json
	jsonnet load.jsonnet > load.json
etl: jsonnet
	go run ../... etl --runscript "${SCRIPT}" --date "${DATE}"