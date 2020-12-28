package cache_service

import (
	"strconv"
	"strings"

	"github.com/EDDYCJY/go-gin-example/pkg/e"
)

type Tag struct {
	ID    int
	Name  string
	State int

	PageNum  int
	PageSize int
}

func (t *Tag) GetTagsKey() string {
	keys := []string{
		e.CACHE_TAG,
		"LIST",
	}

	if t.Name != "" {
		keys = append(keys, "name_" + t.Name)
	}
	if t.State >= 0 {
		keys = append(keys, "state_" + strconv.Itoa(t.State))
	}
	if t.PageNum > 0 {
		keys = append(keys, "pagenum_" + strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		keys = append(keys, "pagesize_" + strconv.Itoa(t.PageSize))
	}

	return strings.Join(keys, "_")
}
