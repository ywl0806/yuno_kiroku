swag:
	~/go/bin/swag i --output docs
kill:
	kill -9 $(shell lsof -t -i:1323)