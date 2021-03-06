package cache

import (
	"context"
	"time"

	"go-common/app/interface/main/mcn/conf"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	//VideoUpTypeCache cache for video types, map key is tid
	VideoUpTypeCache = make(map[int]archive.Type)
)

//RefreshUpTypeAsync refresh in a goroutine
func RefreshUpTypeAsync(tm time.Time) {
	go RefreshUpType(tm)
}

//RefreshUpType refresh
func RefreshUpType(tm time.Time) {
	var url = conf.Conf.Host.Videoup + "/videoup/types"
	var client = bm.NewClient(conf.Conf.HTTPClient)
	var result struct {
		Code int                  `json:"code"`
		Data map[int]archive.Type `json:"data"`
	}
	var err = client.Get(context.Background(), url, "", nil, &result)
	if err != nil {
		log.Error("refresh videoup types fail, err=%v", err)
		return
	}
	log.Info("refresh videoup types ok")
	VideoUpTypeCache = result.Data
}

//GetTidName get tid name
func GetTidName(tid int64) string {
	info, ok := VideoUpTypeCache[int(tid)]
	if !ok {
		return ""
	}
	return info.Name
}

// GetTidNames get tid name
func GetTidNames(tids []int64) (tpNames map[int64]string) {
	tpNames = make(map[int64]string, len(tids))
	for _, tid := range tids {
		info, ok := VideoUpTypeCache[int(tid)]
		if !ok {
			continue
		}
		tpNames[tid] = info.Name
	}
	return
}
