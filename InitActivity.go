// +build android

package go_ricochet_mobile

import (
	"Java/android/databinding/DataBindingUtil"
	"Java/android/os"
	"Java/android/support/v7/app"
	gopkg "Java/RicochetMobile"
	rlayout "Java/RicochetMobile/R/layout"
	"Java/RicochetMobile/databinding"
	"Java/RicochetMobile/ActivityInitBinding"
)

type InitActivity struct {
	app.AppCompatActivity
	binding databinding.ActivityInitBinding
}