package l10n

type Localzation struct {
	AppDescription  string
	ClientTitle     string
	WelcomeFirstRun string
	Command         map[string]string
}

func GetLocalizationStrings() Localzation {
	return EN_US_STRING_MAP
}
