package softlinePayment

type Config struct {
	IdleConnTimeoutSec int
	RequestTimeoutSec  int
	SiteID             int
	Token              string
	URI                string
}
