package xchain

const (
	defaultTestingAccount = "1111111111111111"
	defaultAccountACL     = `
        {
            "pm": {
                "rule": 1,
                "acceptValue": 1.0
            },
            "aksWeight": {
                "` + "xchain" + `": 1.0
            }
        }
        `
)
