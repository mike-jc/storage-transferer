package testsInit

import (
	"github.com/astaxie/beego"
	"service-recordingStorage/system"
)

func init() {
	system.SetAppDirToCurrentDir(1)
	beego.TestBeegoInit(system.AppDir())
}
