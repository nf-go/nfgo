package nsecurity

import (
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"nfgo.ga/nfgo/nconf"
)

// NewEnforcer -
func NewEnforcer(securityConfig *nconf.SecurityConfig, db *gorm.DB) (casbin.IEnforcer, error) {
	adpt, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("fail to create casbin gorm-adapter: %w", err)
	}

	m := model.NewModel()
	if err := m.LoadModelFromText(securityConfig.Model); err != nil {
		return nil, fmt.Errorf("fail to load model from seucrity config: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(m, adpt)
	if err != nil {
		return nil, fmt.Errorf("fail to create enforcer: %w", err)
	}
	return enforcer, nil
}

// InitPolicy -
func InitPolicy(enforcer casbin.IEnforcer, securityConfig *nconf.SecurityConfig, rules [][]string) error {
	// load from db
	if err := enforcer.LoadPolicy(); err != nil {
		return err
	}
	// clear all
	enforcer.ClearPolicy()

	// add policies from config
	for _, anno := range securityConfig.Anons {
		if _, err := enforcer.AddNamedPolicy("p", "anonymous", anno, "*"); err != nil {
			return err
		}
	}
	for _, policy := range securityConfig.Policies {
		ps := strings.Split(policy, ",")
		if len(ps) > 2 {
			params := make([]interface{}, len(ps))
			for i := range ps {
				params[i] = strings.TrimSpace(ps[i])
			}
			enforcer.AddNamedPolicy(ps[0], params[1:]...)
		}
	}

	// save policies to db
	if err := enforcer.SavePolicy(); err != nil {
		return err
	}

	if _, err := enforcer.AddPolicies(rules); err != nil {
		return err
	}

	return nil
}
