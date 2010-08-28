
include $(GOROOT)/src/Make.inc

TARG=tyrant
CGOFILES=\
	 tyrant.go
#CGO_LDFLAGS=ttwrapper.o -ltokyotyrant
CGO_LDFLAGS=-ltokyotyrant


CLEANFILES+=connect

include $(GOROOT)/src/Make.pkg

%: ttwrapper.o install %.go
	$(GC) $*.go
	$(LD) -o $@ $*.$O

ttwrapper.o: ttwrapper.c
	gcc -fPIC -O2 -o ttwrapper.o -c ttwrapper.c





