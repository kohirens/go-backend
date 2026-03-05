package backend

import (
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www/storage"
)

const (
	KeyGoogleProvider = "gp"
	KeySessionManager = "sm"

	// MetaRefresh HTML template to redirect the client.
	MetaRefresh = `<!DOCTYPE html>
<html>
	<head><meta http-equiv="refresh" content="0; url='%s'"></head>
	<body></body>
</html>`
	TmplSuffix = "tmpl"
)

const (
	KeyAccountManager = "am"
	PrefixAccounts    = "accounts"
	PrefixSecrets     = "secrets"
)

var (
	Log     = &logger.Standard{}
	TmplDir = "templates"
)

func NewAccountExec(store storage.Storage) *AccountExec {
	return &AccountExec{store: store}
}
