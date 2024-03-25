package models

const (
	CSEBranch      = "CSE"
	ISEBranch      = "ISE"
	CivilBranch    = "CIVIL"
	MechBranch     = "MECH"
	EceBranch      = "ECE"
	EeeBranch      = "EEE"
	AcceptedStatus = "ACCEPTED"
	PendingStatus  = "PENDING"
	RejectedStatus = "REJECTED"

	MinimumAge         = 22
	MinimumNameLength  = 3
	MinimumPhoneLength = 10
	MaximumPhoneLength = 12
	ReadoutTimer       = 1

	APIKey      = "X-API-KEY"
	ContentType = "application/json"

	// Categories
	Mass      = "MASS"
	DreamIT   = "DREAM IT"
	OpenDream = "OPEN DREAM"
	Core      = "CORE"

	// Table names
	CompanyTable = "company"
	StudentTable = "student"
)

// nolint : gochecknoglobals // required
var (
	ValidStudentStatus = []string{
		AcceptedStatus, PendingStatus, RejectedStatus,
	}
	ValidBranch = []string{
		CSEBranch, ISEBranch, CivilBranch, MechBranch, EeeBranch, EceBranch,
	}
	ValidCategory = []string{
		Mass, DreamIT, OpenDream, Core,
	}
)

func Contains(input []string, str string) bool {
	for _, val := range input {
		if val == str {
			return true
		}
	}

	return false
}
