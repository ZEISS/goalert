package acs

import (
	"database/sql"
	"net/http"

	"github.com/target/goalert/user/contactmethod"
)

const (
	msgParamID     = "msgID"
	msgParamSubID  = "msgSubjectID"
	msgParamBody   = "msgBody"
	msgParamBundle = "msgBundle"
)

// Config contains the details needed to interact with Azure Communication Services for SMS
type Config struct {
	// BaseURL can be used to override the Azure Communication Services API and Lookup URL bases.
	BaseURL string

	// Client is an optional net/http client to use, if nil the global default is used.
	Client *http.Client

	// CMStore is used for storing and fetching metadata (like carrier information).
	CMStore *contactmethod.Store

	// DB is used for storing DB connection data (needed for carrier metadata dbtx).
	DB *sql.DB
}
