package config

import (
	"fmt"
	mexadomain "mexa/internal/domains/mexa"
)

type Config struct {
	Token    *string              `json:"token,omitempty"`
	DbPath   *string              `json:"db_path,omitempty"`
	CasesDir *string              `json:"cases_dir,omitempty"`
	AdminIds []mexadomain.UserId  `json:"admin_ids,omitempty"`
	Batch    *mexadomain.Batch    `json:"batch,omitempty"`
	Exercise *mexadomain.Exercise `json:"exercise,omitempty"`
}

func (c *Config) Validate() error {
	if c.Token == nil {
		return fmt.Errorf("token is required")
	} else if c.DbPath == nil {
		return fmt.Errorf("db_path is required")
	} else if c.CasesDir == nil {
		return fmt.Errorf("cases_dir is required")
	} else if c.Batch == nil {
		return fmt.Errorf("batch is required")
	} else if c.Exercise == nil {
		return fmt.Errorf("exercise is required")
	} else if len(c.AdminIds) == 0 {
		return fmt.Errorf("admin_ids is required")
	}

	return nil
}
