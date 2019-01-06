package gf_core

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/globalsign/mgo"
)
//-------------------------------------------------
type Runtime_sys struct {
	Service_name_str string
	Log_fun          func(string,string)
	Mongodb_coll     *mgo.Collection
}