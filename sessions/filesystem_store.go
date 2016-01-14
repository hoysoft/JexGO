package sessions

import (
	"github.com/gorilla/sessions"
)

type FilesystemStore interface {
	Store
	Options(Options)
	MaxLength(int)
}

func NewFilesystemStore(path string, keyPairs ...[]byte) FilesystemStore {
	return &filesystemStore{sessions.NewFilesystemStore(path, keyPairs...)}
}

type filesystemStore struct {
	*sessions.FilesystemStore
}

func (c *filesystemStore) MaxLength(l int) {
	c.FilesystemStore.MaxLength(l)
}

func (c *filesystemStore) Options(options Options) {
	c.FilesystemStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
