package task

import (
	"github.com/robfig/cron/v3"
	"regexp"
	"time"
)

const DBPath = "data/change-detector.db"

var MulRn, _ = regexp.Compile("\n{2,}")
var MulSp, _ = regexp.Compile("\\s{2,}")

var TimeZone, _ = time.LoadLocation("Asia/Shanghai")
var cronObj = cron.New(cron.WithLocation(TimeZone))
var cronId cron.EntryID
