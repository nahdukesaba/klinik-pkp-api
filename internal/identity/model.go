package identity

import (
	"strings"

	"github.com/google/uuid"
)

const KlinikApplicationCode = "klinik_pkp"

// KlinikIdentity is the canonical read model used by Klinik authentication and
// session endpoints. It intentionally contains no klinikpkp_user fields.
type KlinikIdentity struct {
	AppUserID        uuid.UUID  `gorm:"column:app_user_id"`
	AuthUserID       *uuid.UUID `gorm:"column:auth_user_id"`
	Email            string     `gorm:"column:email"`
	FullName         string     `gorm:"column:full_name"`
	Phone            *string    `gorm:"column:phone"`
	NIP              *string    `gorm:"column:nip"`
	Role             string     `gorm:"column:role"`
	MembershipStatus string     `gorm:"column:membership_status"`
	ApplicationCode  string     `gorm:"column:application_code"`
}

func (i *KlinikIdentity) Normalize() {
	if i == nil {
		return
	}
	i.Email = strings.TrimSpace(i.Email)
	i.FullName = strings.TrimSpace(i.FullName)
	i.Role = strings.ToLower(strings.TrimSpace(i.Role))
	i.MembershipStatus = strings.ToLower(strings.TrimSpace(i.MembershipStatus))
	i.ApplicationCode = strings.ToLower(strings.TrimSpace(i.ApplicationCode))
}

func (i *KlinikIdentity) IsApproved() bool {
	if i == nil {
		return false
	}
	i.Normalize()
	return i.AppUserID != uuid.Nil &&
		i.AuthUserID != nil && *i.AuthUserID != uuid.Nil &&
		i.ApplicationCode == KlinikApplicationCode &&
		i.MembershipStatus == "approved" &&
		(i.Role == "admin" || i.Role == "user")
}
