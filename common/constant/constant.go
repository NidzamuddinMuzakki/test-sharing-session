package constants

const (
	XRequestIDHeader = "X-REQUEST-ID"
	XUserIDHeader    = "X-User-Id"
	XUsernameHeader  = "X-Username"
	System           = "SYSTEM"
	XHeaderTest      = "X-header-test"
	KTP              = "ktp"
	KK               = "kk"
)

const (
	SQLErrNoRows        = "stmt[0]: sql: no rows in result set"
	SQLERRDuplicateRows = "duplicate key value violates unique constraint"
	SQLDuplicateError   = "Duplicate entry"
)

const (
	KtpKawin      = "kawin"
	KtpBelumKawin = "belum kawin"
	KtpCeraiHidup = "cerai hidup"
	KtpCeraiMati  = "cerai mati"

	EmployeeKawin      = "Sudah Menikah"
	EmployeeBelumKawin = "Belum Menikah"
	EmployeeCerai      = "Cerai"
)

const (
	SudahMenikahCode = 7
	BelumMenikahCode = 8
	CeraiCode        = 9

	SudahMenikah = "Sudah Menikah"
	BelumMenikah = "Belum Menikah"
	Cerai        = "Cerai"
)
