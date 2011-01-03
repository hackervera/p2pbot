include $(GOROOT)/src/Make.inc

TARG=bot
GOFILES=server.go webstuff.go ircstuff.go data.go socket.go

include $(GOROOT)/src/Make.cmd
