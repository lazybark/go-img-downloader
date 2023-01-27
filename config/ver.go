package config

import "github.com/lazybark/go-helpers/semver"

var (
	//Ver is the current app version according to Semver
	Ver = semver.Ver{
		Major:       1,
		Minor:       5,
		Patch:       2,
		Stable:      true,
		ReleaseNote: "2023.01.27",
	}
)
