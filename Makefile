vet:
	jsonnet config.jsonnet | cue vet schema.cue -

send:
	jsonnet config.jsonnet | go run .

live-symlink:
	ln -sf "$(dir $(abspath $(firstword $(MAKEFILE_LIST))))live" \
		"$(shell echo ~)/Music/Ableton/User Library/Remote Scripts/ec4"

live-log:
	tail -f "$$(ls -d ~/Library/Preferences/Ableton/* | tail -1)/Log.txt"
