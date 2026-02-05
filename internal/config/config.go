package config

import (
	"fmt"
	mexadomain "mexa/internal/domains/mexa"
	"slices"
)

type Config struct {
	Token         *string              `json:"token,omitempty"`
	DbPath        *string              `json:"db_path,omitempty"`
	CasesDir      *string              `json:"cases_dir,omitempty"`
	AdminIds      []mexadomain.UserId  `json:"admin_ids,omitempty"`
	Batch         *mexadomain.Batch    `json:"batch,omitempty"`
	Exercise      *mexadomain.Exercise `json:"exercise,omitempty"`
	SectionSizes  [][]int              `json:"section_sizes,omitempty"`
	InvalidFourDs []int                `json:"invalid_4Ds,omitempty"`
}

type FourDValidator func(int) bool

func (c *Config) FourDValidator() FourDValidator {
	return func(fd int) bool {

		nPlt := len(c.SectionSizes)
		list := make([]int, 0, 4)
		{
			x := fd
			for x > 0 {
				list = append(list, x%10)
				x /= 10
			}
			slices.Reverse(list)
		}

		// Check that it is a four digit number
		if len(list) != 4 {
			return false
		}

		{
			pl := list[0]
			if pl < 1 || pl > nPlt {
				return false
			}

			sc := list[1]
			if sc < 1 || sc > len(c.SectionSizes[pl-1]) {
				return false
			}

			n := (list[2] * 10) + list[3]
			if n < 1 || n > c.SectionSizes[pl-1][sc-1] {
				return false
			}
		}

		for _, v := range c.InvalidFourDs {
			if fd == v {
				return false
			}
		}
		return true
	}
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
