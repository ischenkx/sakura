package notify

import "time"

type TimeStampOption struct {
	UnixTime int64
}

func TimeStamp(t time.Time) TimeStampOption {
	return UnixTimestamp(t.UnixNano())
}

func UnixTimestamp(t int64) TimeStampOption {
	return TimeStampOption{UnixTime: t}
}

type MetaInfoOption struct {
	Data interface{}
}

func MetaInfo(data interface{}) MetaInfoOption {
	return MetaInfoOption{data}
}
