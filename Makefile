include $(GOROOT)/src/Make.inc

TARG=bot
GOFILES=server.go webstuff.go ircstuff.go data.go socket.go crypto.go

include $(GOROOT)/src/Make.cmd
