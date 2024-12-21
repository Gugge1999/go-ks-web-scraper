package constants

// TODO: Fixa så att man inte behöver sätta prod och dev separat
const (
	IntervalInMin    uint32 = 10
	IntervalInMsProd        = IntervalInMin * 60_000
	IntervalInMsDev         = IntervalInMin * 1_500
)
