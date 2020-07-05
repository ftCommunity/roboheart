package packages

import (
	"github.com/ftCommunity/roboheart/package/marshallers"
)

type Package struct {
	Id       string
	Variants map[string]Variant
}

type Variant struct {
	PackageId     string                         `json:"pkgid"`
	VariantId     string                         `json:"variantid"`
	AuthorId      string                         `json:"authorid"`
	Version       marshallers.Version            `json:"version"`  //version of this variant
	Platform      *marshallers.Regexp            `json:"platform"` //allowed platform as regex
	Device        *marshallers.Regexp            `json:"device"`   //allowed device as regex
	Firmware      *marshallers.Range             `json:"firmware"` //allowed firmware versions
	Multiversion  bool                           `json:"multiversion"`
	OtherVersions *marshallers.Range             `json:"otherversions"` //if set defnies allowed versions of same package. if unset all versions are allowed
	Dependencies  map[string]Dependency          `json:"dependencies"`  //maps named dependency folders to another package
	Implements    map[string]marshallers.Version `json:"implements"`
	Name          *string                        `json:"name"`
	Desc          *string                        `json:"desc"`
	Languages     map[string]Language            `json:"languages"`
	Icon          *string                        `json:"icon"`
	Category      *string                        `json:"category"`
	AuthorName    *string                        `json:"authorname"`
	Url           *string                        `json:"url"`
	Management    *string                        `json:"management"`
	Frontend      *string                        `json:"frontend"`
}

type Language struct {
	Name *string `json:"name"`
	Desc *string `json:"desc"`
}

type Dependency struct {
	OneOf []DependencyOption `json:"oneof"`
}

type DependencyOption struct {
	ID      string             `json:"id"`
	Version *marshallers.Range `json:"version"`
}
