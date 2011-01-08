include $(GOROOT)/src/Make.inc

TARG=bot
GOFILES=server.go ajax.go ircstuff.go data.go socket.go crypto.go

include $(GOROOT)/src/Make.cmd
