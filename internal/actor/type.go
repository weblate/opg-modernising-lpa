package actor

type Type int

const (
	TypeNone Type = iota
	TypeDonor
	TypeAttorney
	TypeReplacementAttorney
	TypeCertificateProvider
	TypePersonToNotify
)

func (t Type) String() string {
	switch t {
	case TypeDonor:
		return "donor"
	case TypeAttorney:
		return "attorney"
	case TypeReplacementAttorney:
		return "replacementAttorney"
	case TypeCertificateProvider:
		return "certificateProvider"
	case TypePersonToNotify:
		return "personToNotify"
	default:
		return ""
	}
}

type Types struct {
	None                Type
	Donor               Type
	Attorney            Type
	ReplacementAttorney Type
	CertificateProvider Type
	PersonToNotify      Type
}

var ActorTypes = Types{
	None:                TypeNone,
	Donor:               TypeDonor,
	Attorney:            TypeAttorney,
	ReplacementAttorney: TypeReplacementAttorney,
	CertificateProvider: TypeCertificateProvider,
	PersonToNotify:      TypePersonToNotify,
}