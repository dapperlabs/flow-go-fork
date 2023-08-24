package synchronization

// spamProbabilityMultiplier is used to convert probability factor to an integer as well as a maximum value - 1
// random number that can be generated by the random number generator.
const spamProbabilityMultiplier = 1001

type SpamReportConfig struct {
	syncRequestProbability float32
}

func NewSpamReportConfig() *SpamReportConfig {
	return &SpamReportConfig{
		// create misbehavior report 1/100 message requests
		syncRequestProbability: 0.01,
	}
}
