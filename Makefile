build:
	unzip resources/tiles.zip
	docker buildx build --platform linux/amd64 --push -t registry.digitalocean.com/atlas-game-lookout/lookoutbot .