package api

import (
	"gitlab.void-ptr.org/go/schism/pkg/business"
)

// FeatureSet of the api
type FeatureSet struct {
	Devices business.DeviceSupport `json:"devices"`
}

// Features supported
var Features = FeatureSet{
	Devices: business.DeviceSupport{
		Enabled: true,
	},
}
