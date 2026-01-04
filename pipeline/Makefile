SCRIPT?=pipeline.json
DB?=2026-01-03.sqlite

jsonnet:
	jsonnet pipeline.jsonnet > pipeline.json
etl: jsonnet
	go run ../... etl --runscript "${SCRIPT}" --sqlite "${DB}"