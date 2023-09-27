package sessions_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sessionReaperConfig struct{}

func (c sessionReaperConfig) SessionTimeout() models.Duration {
	return models.MustMakeDuration(42 * time.Second)
}

func (c sessionReaperConfig) SessionReaperExpiration() models.Duration {
	return models.MustMakeDuration(142 * time.Second)
}

func TestSessionReaper_ReapSessions(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	config := sessionReaperConfig{}
	lggr := logger.TestLogger(t)
	orm := sessions.NewORM(db, config.SessionTimeout().Duration(), lggr, pgtest.NewQConfig(true), audit.NoopLogger)

	rw := sessions.NewSessionReaperWorker(db.DB, config, lggr)
	r := utils.NewSleeperTask(rw)

	t.Cleanup(func() {
		assert.NoError(t, r.Stop())
	})

	tests := []struct {
		name     string
		lastUsed time.Time
		wantReap bool
	}{
		{"current", time.Now(), false},
		{"expired", time.Now().Add(-config.SessionTimeout().Duration()), false},
		{"almost stale", time.Now().Add(-config.SessionReaperExpiration().Duration()), false},
		{"stale", time.Now().Add(-config.SessionReaperExpiration().Duration()).
			Add(-config.SessionTimeout().Duration()), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user := cltest.MustRandomUser(t)
			require.NoError(t, orm.CreateUser(&user))

			session := sessions.NewSession()
			session.Email = user.Email

			_, err := db.Exec("INSERT INTO sessions (last_used, email, id, created_at) VALUES ($1, $2, $3, now())", test.lastUsed, user.Email, test.name)
			require.NoError(t, err)

			t.Cleanup(func() {
				_, err2 := db.Exec("DELETE FROM sessions where email = $1", user.Email)
				require.NoError(t, err2)
			})

			r.WakeUp()
			<-rw.RunSignal()
			sessions, err := orm.Sessions(0, 10)
			assert.NoError(t, err)

			if test.wantReap {
				assert.Len(t, sessions, 0)
			} else {
				assert.Len(t, sessions, 1)
			}
		})
	}
}
