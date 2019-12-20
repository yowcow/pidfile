all:
	go run main.go

sig/%: PIDFILE := var/run/test.pid
sig/%:
	cat $(PIDFILE) | xargs -I{} kill -$* {}

.PHONY: all sig/%
