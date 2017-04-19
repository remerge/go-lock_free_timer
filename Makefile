PROJECT := go-lock_free_timer
PACKAGE := github.com/remerge/$(PROJECT)

GOMETALINTER_OPTS := -D golint

include Makefile.common
