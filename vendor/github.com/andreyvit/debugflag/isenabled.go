package debugflag

func IsEnabled(flag string) bool {
	for ; flag != ""; flag = parentNamespace(flag) {
		ena, prsnt := enabled[flag]
		if prsnt {
			return ena
		}
	}
	return enabled["all"]
}
