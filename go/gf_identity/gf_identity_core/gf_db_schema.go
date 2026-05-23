
package gf_identity_core

import (
	"context"
	gf_core "github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------

func DBsqlCreateTables(pCtx context.Context,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	sqlStr := `
	CREATE TABLE IF NOT EXISTS gf_users (
		v                 VARCHAR(255),
		id                TEXT,
		deleted           BOOLEAN DEFAULT FALSE,
		creation_time     TIMESTAMP DEFAULT NOW(),
		user_type         VARCHAR(255), -- "admin" or "standard"
		user_name         TEXT NOT NULL,
		screen_name       TEXT,
		description       TEXT,
		addresses_eth     TEXT[],
		email             TEXT,
		email_confirmed   BOOLEAN DEFAULT FALSE,
		profile_image_url TEXT,
		banner_image_url  TEXT,

		PRIMARY KEY(id)
	);

	CREATE TABLE IF NOT EXISTS gf_auth_session (
		v              VARCHAR(255),
		id             TEXT,
		deleted        BOOLEAN DEFAULT FALSE,
		creation_time  TIMESTAMP DEFAULT NOW(),
		user_id        TEXT,

		login_complete              BOOLEAN NOT NULL DEFAULT FALSE,
		login_success_redirect_url  TEXT,
		logout_success_redirect_url TEXT,

		profile JSON,

		auth_subsystem_type VARCHAR(20) DEFAULT 'auth0',
		auth_method         VARCHAR, -- "google-oauth2", "github-oauth2", "email", etc.
		user_id_idp         TEXT,    -- user ID as per the identity provider (idp)
		user_agent          VARCHAR, -- user agent string of the browser/client

		mcp BOOLEAN DEFAULT FALSE, -- if its a AI/MCP session

		PRIMARY KEY(id),
		FOREIGN KEY (user_id) REFERENCES gf_users(id)
	);

	CREATE TABLE IF NOT EXISTS gf_users_creds (
		v             VARCHAR(255),
		id            TEXT,
		deleted       BOOLEAN DEFAULT FALSE,
		creation_time TIMESTAMP DEFAULT NOW(),
		user_id       TEXT,
		user_name     TEXT,
		pass_salt     TEXT,
		pass_hash     TEXT,

		PRIMARY KEY(id),
		FOREIGN KEY (user_id) REFERENCES gf_users(id)
	);

	CREATE TABLE IF NOT EXISTS gf_users_email_confirm (
		user_name          TEXT,
		user_id            TEXT,
		confirm_code       TEXT,
		creation_time      TIMESTAMP DEFAULT NOW(),

		PRIMARY KEY(confirm_code),
		FOREIGN KEY (user_id) REFERENCES gf_users(id)
	);

	CREATE TABLE IF NOT EXISTS gf_auth_login_attempts (
		v                VARCHAR(255),
		id               TEXT,
		deleted          BOOLEAN DEFAULT FALSE,
		creation_time    TIMESTAMP DEFAULT NOW(),

		user_type        VARCHAR(255), -- "admin" or "standard"
		user_id          TEXT,
		user_name        VARCHAR(255),

		-- only used for auth0 authentication since initially only the session_id
		-- is know, and not the user_name/id. that info is only known afterwards
		-- once the user is redirected back to GF from auth0 dialogs
		auth_session_id TEXT,

		pass_confirmed   BOOLEAN,
		email_confirmed  BOOLEAN,
		mfa_confirmed    BOOLEAN,

		PRIMARY KEY(id),
		FOREIGN KEY (user_id) REFERENCES gf_users(id),
		FOREIGN KEY (auth_session_id) REFERENCES gf_auth_session(id)
	);

	CREATE TABLE IF NOT EXISTS gf_users_invite_list (
		user_email      TEXT UNIQUE NOT NULL,
		creation_time   TIMESTAMP DEFAULT NOW(),
		deleted         BOOLEAN DEFAULT FALSE,

		PRIMARY KEY(user_email)
	);
	`

	_, err := pRuntimeSys.SQLdb.ExecContext(pCtx, sqlStr)
	if err != nil {
		gfErr := gf_core.ErrorCreate("failed to create gf_identity related tables in the DB",
			"sql_table_creation",
			map[string]interface{}{},
			err, "gf_identity_core", pRuntimeSys)
		return gfErr
	}

	return nil
}