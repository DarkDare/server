package server

import (
	"time"
)

const (
	ServerVersion                 = "0.1"
	OneTB                         = 1099511627776
	TwoToTheSixtyThree            = 9223372036854775808
	SubmissionInitialAttempts     = 5
	SubmissionInitialBackoff      = 2 * time.Microsecond
	SubmissionMaxSubmitDelay      = 2 * time.Second
	VarIdleTimeoutMin             = 50 * time.Millisecond
	VarIdleTimeoutRange           = 250
	FrameLockMinExcessSize        = 100
	FrameLockMinRatio             = 2
	ConnectionRestartDelayRangeMS = 5000
	ConnectionRestartDelayMin     = 3 * time.Second
	MostRandomByteIndex           = 7 // will be the lsb of a big-endian client-n in the txnid.
)
