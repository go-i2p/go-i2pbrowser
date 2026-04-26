package goi2pbrowser

import (
	"log"

	fcw "github.com/go-wbg/go-fpw"
)

// BrowseStrict launches a Firefox browser configured to use I2P and waits for it to exit.
// The profile is in "Strict" mode
func (i *I2PBrowser) BrowseStrict(url ...string) {
	profilePath, err := UnpackBase(i.ProfileDir)
	if err != nil {
		log.Println(err)
		return
	}
	FIREFOX, ERROR := fcw.BasicFirefox(profilePath, false, url...)
	if ERROR != nil {
		log.Println(ERROR)
		return
	}
	defer FIREFOX.Close()
	<-FIREFOX.Done()
}

// BrowseUsability launches a Firefox browser configured to use I2P and waits for it to exit.
// The profile is in "Usability" mode
func (i *I2PBrowser) BrowseUsability(url ...string) {
	profilePath, err := UnpackUsability(i.ProfileDir)
	if err != nil {
		log.Println(err)
		return
	}
	FIREFOX, ERROR := fcw.BasicFirefox(profilePath, false, url...)
	if ERROR != nil {
		log.Println(ERROR)
		return
	}
	defer FIREFOX.Close()
	<-FIREFOX.Done()
}

// BrowseApp launches a Firefox browser configured to use I2P and waits for it to exit.
// The profile is in "Usability" mode
func (i *I2PBrowser) BrowseApp(url ...string) {
	profilePath, err := UnpackUsability(i.ProfileDir)
	if err != nil {
		log.Println(err)
		return
	}
	FIREFOX, ERROR := fcw.WebAppFirefox(profilePath, false, true, url...)
	if ERROR != nil {
		log.Println(ERROR)
		return
	}
	defer FIREFOX.Close()
	<-FIREFOX.Done()
}
