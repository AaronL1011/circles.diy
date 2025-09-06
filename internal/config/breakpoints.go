package config

// Breakpoints defines all responsive breakpoint values used in CSS
type Breakpoints struct {
	// Standard Bootstrap-like breakpoints (mobile-first)
	XS  string // 0px - extra small devices
	SM  string // 576px - small devices (landscape phones)
	MD  string // 768px - medium devices (tablets)
	LG  string // 992px - large devices (desktops)
	XL  string // 1200px - extra large devices (large desktops)
	XXL string // 1400px - extra extra large devices

	// Legacy breakpoints (for gradual migration)
	Mobile        string // 480px
	Tablet        string // 768px
	SmallDesktop  string // 900px
	Desktop       string // 1024px
}

// GetBreakpoints returns the configured breakpoint values
func GetBreakpoints() Breakpoints {
	return Breakpoints{
		// Standard breakpoints
		XS:  "0",
		SM:  "576px",
		MD:  "768px",
		LG:  "992px",
		XL:  "1200px",
		XXL: "1400px",

		// Legacy breakpoints
		Mobile:       "480px",
		Tablet:       "768px",
		SmallDesktop: "900px",
		Desktop:      "1024px",
	}
}

// GetVariableMap returns a map of CSS variable names to their values
// for use in CSS variable substitution
func GetVariableMap() map[string]string {
	bp := GetBreakpoints()
	
	return map[string]string{
		"--breakpoint-xs":             bp.XS,
		"--breakpoint-sm":             bp.SM,
		"--breakpoint-md":             bp.MD,
		"--breakpoint-lg":             bp.LG,
		"--breakpoint-xl":             bp.XL,
		"--breakpoint-xxl":            bp.XXL,
		"--breakpoint-mobile":         bp.Mobile,
		"--breakpoint-tablet":         bp.Tablet,
		"--breakpoint-small-desktop":  bp.SmallDesktop,
		"--breakpoint-desktop":        bp.Desktop,
	}
}