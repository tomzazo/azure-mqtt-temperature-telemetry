.PHONY: setup build zip upload deploy

setup:
	sudo apt install zip direnv

build:
	env GOARCH=arm go build -o bin/mqtt_sender

zip:
	-rm mqtt_sender.zip
	zip -r -j mqtt_sender.zip bin/mqtt_sender
	zip -r mqtt_sender.zip scripts systemd sensors config.yml

upload:
	scp mqtt_sender.zip $(TARGET_USER)@$(TARGET_IP):$(TARGET_DIRECTORY)

deploy: build zip upload
	