package pkg

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/spf13/viper"
)

var node *snowflake.Node

func Init() (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", viper.GetString("snowflake.start_time"))
	if err != nil {
		return err
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(viper.GetInt64("snowflake.machine_id"))
	return
}

func GenId() int64 {
	Init()
	return node.Generate().Int64()
}
