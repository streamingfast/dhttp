package dhttp

import (
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()
var _ = logging.PackageLogger("dhttp", "github.com/streamingfast/dhttp", &zlog)
