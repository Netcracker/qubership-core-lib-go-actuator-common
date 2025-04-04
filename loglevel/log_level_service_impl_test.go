package loglevel

import (
	"os"
	"strings"
	"testing"

	"github.com/netcracker/qubership-core-lib-go/v3/configloader"
	"github.com/netcracker/qubership-core-lib-go/v3/logging"
	"github.com/stretchr/testify/require"
)

func TestEnvPropertySource(t *testing.T) {
	var loggerName = "logpackage"

	os.Setenv("LOGGING_LEVEL_LOGPACKAGE", logging.LvlCrit.String())
	defer func() {
		os.Clearenv()
	}()

	configloader.Init(configloader.EnvPropertySource())

	_ = logging.GetLogger(loggerName)

	logLevelService, err := NewLogLevelService()
	require.Nil(t, err)

	resp, err := logLevelService.GetLogLevels()
	require.Nil(t, err)
	logLevel := (*resp)[loggerName]
	require.Equal(t, strings.ToUpper(logging.LvlCrit.String()), logLevel)
}
