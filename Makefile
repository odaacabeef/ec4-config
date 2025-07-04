vet:
	jsonnet config.jsonnet | cue vet schema.cue -

send:
	jsonnet config.jsonnet | go run .
