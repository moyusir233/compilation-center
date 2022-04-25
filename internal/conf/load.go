package conf

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
)

func LoadConfig(path string, logger log.Logger) (*Bootstrap, error) {
	c := config.New(
		config.WithSource(
			file.NewSource(path),
		),
		config.WithLogger(logger),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		return nil, err
	}

	var bc Bootstrap
	if err := c.Scan(&bc); err != nil {
		return nil, err
	}

	return &bc, nil
}
