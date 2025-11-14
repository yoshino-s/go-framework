package authz_casbin

import (
	"github.com/casbin/casbin/v2/log"
	"go.uber.org/zap"
)

var _ log.Logger = (*casbinZapLogger)(nil)

type casbinZapLogger struct {
	logger *zap.Logger
}

// EnableLog implements log.Logger.
func (c *casbinZapLogger) EnableLog(enabled bool) {
}

// IsEnabled implements log.Logger.
func (c *casbinZapLogger) IsEnabled() bool {
	return true
}

// LogEnforce implements log.Logger.
func (c *casbinZapLogger) LogEnforce(matcher string, request []interface{}, result bool, explains [][]string) {
	c.logger.Debug("Casbin Enforce",
		zap.String("matcher", matcher),
		zap.Any("request", request),
		zap.Bool("result", result),
		zap.Any("explains", explains),
	)
}

// LogError implements log.Logger.
func (c *casbinZapLogger) LogError(err error, msg ...string) {
	if len(msg) > 0 {
		c.logger.Error(msg[0], zap.Error(err))
	} else {
		c.logger.Error("Casbin Error", zap.Error(err))
	}
}

// LogModel implements log.Logger.
func (c *casbinZapLogger) LogModel(model [][]string) {
	c.logger.Debug("Casbin Model", zap.Any("model", model))
}

// LogPolicy implements log.Logger.
func (c *casbinZapLogger) LogPolicy(policy map[string][][]string) {
	c.logger.Debug("Casbin Policy", zap.Any("policy", policy))
}

// LogRole implements log.Logger.
func (c *casbinZapLogger) LogRole(roles []string) {
	c.logger.Debug("Casbin Role", zap.Any("roles", roles))
}

func newCasbinZapLogger(logger *zap.Logger) *casbinZapLogger {
	return &casbinZapLogger{
		logger: logger,
	}
}

var _ log.Logger = (*mockLogger)(nil)

type mockLogger struct {
	Roles []string
}

func (m *mockLogger) EnableLog(bool) {
	//
}

// IsEnabled implements log.Logger.
func (m *mockLogger) IsEnabled() bool {
	return true
}

// LogEnforce implements log.Logger.
func (m *mockLogger) LogEnforce(matcher string, request []interface{}, result bool, explains [][]string) {
	panic("unimplemented")
}

// LogError implements log.Logger.
func (m *mockLogger) LogError(err error, msg ...string) {
	panic("unimplemented")
}

// LogModel implements log.Logger.
func (m *mockLogger) LogModel(model [][]string) {
	panic("unimplemented")
}

// LogPolicy implements log.Logger.
func (m *mockLogger) LogPolicy(policy map[string][][]string) {
	panic("unimplemented")
}

// LogRole implements log.Logger.
func (m *mockLogger) LogRole(roles []string) {
	m.Roles = roles
}
