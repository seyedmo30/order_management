package pkg

// service
var (
	StatusOrderManagementHigh   = "High"
	StatusOrderManagementNormal = "Normal"

	StatusOrderManagementPending   = "Pending"
	StatusOrderManagementProcessed = "Processed"
	StatusOrderManagementFailed    = "Failed"
)

// repository
var (
	NotFoundRepositoryMessage            = "Error Repository: Item Not Found in the Repository"
	ExpectedExactlyOneRepositoryMessage  = "Error Repository: Expected exactly one item but found multiple records."
	InternalServerErrorRepositoryMessage = "Error Repository: An unexpected issue has occurred. Please try again later."
	DuplicateEntryRepositoryMessage      = "Error Repository: Duplicate entry. A record with this value already exists in the system."
	RequiredFieldRepositoryMessage       = "Error Repository: required cannot be null. This field is required and must contain a valid value for the transaction to proceed."
)
