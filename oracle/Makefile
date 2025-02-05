generate:
	python3 scripts/generate.py

deploy:
	GOARCH=amd64 GOOS=linux go build .
	ssh root@37.27.192.7 'killall gringotts -9'
	scp gringotts root@37.27.192.7:/root/gringotts
	scp -r config root@37.27.192.7:/root/gringotts
	ssh root@37.27.192.7 'cd gringotts; ./run' &
	rm gringotts