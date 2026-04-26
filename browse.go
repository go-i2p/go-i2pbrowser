package goi2pbrowser

import (
	"log"

	fcw "github.com/go-wbg/go-fpw"
)

// BrowseStrict launches a Firefox browser configured to use I2P and waits for it to exit.
// The profile is in "Strict" mode
func (i *I2PBrowser) BrowseStrict(url ...string) error {
	profilePath, err := UnpackBase(i.ProfileDir)
	if err != nil {
		log.Println(err)
		return err
	}
	FIREFOX, ERROR := fcw.BasicFirefox(profilePath, false, url...)
	if ERROR != nil {
		log.Println(ERROR)
		return ERROR
	}
	defer func() {
		if err := FIREFOX.Close(); err != nil {
			log.Printf("WARNING: failed to close Firefox browser: %v", err)
		}
	}()
	<-FIREFOX.Done()
	return nil
}

// BrowseUsability launches a Firefox browser configured to use I2P and waits for it to exit.
// The profile is in "Usability" mode
func (i *I2PBrowser) BrowseUsability(url ...string) error {
	profilePath, err := UnpackUsability(i.ProfileDir)
	if err != nil {
		log.Println(err)
		return err
	}
	FIREFOX, ERROR := fcw.BasicFirefox(profilePath, false, url...)
	if ERROR != nil {
		log.Println(ERROR)
		return ERROR
	}
	defer func() {
		if err := FIREFOX.Close(); err != nil {
			log.Printf("WARNING: failed to close Firefox browser: %v", err)
		}
	}()
	<-FIREFOX.Done()
	return nil
}

// BrowseApp launches a Firefox browser configured to use I2P and waits for it to exit.
// The profile is in "Usability" mode with webapp/kiosk-style window.
func (i *I2PBrowser) BrowseApp(url ...string) error {
	profilePath, err := UnpackUsability(i.ProfileDir)
	if err != nil {
		log.Println(err)
		return err
	}
	FIREFOX, ERROR := fcw.WebAppFirefox(profilePath, false, true, url...)
	if ERROR != nil {
		log.Println(ERROR)
		return ERROR
	}
	defer func() {
		if err := FIREFOX.Close(); err != nil {
			log.Printf("WARNING: failed to close Firefox browser: %v", err)
		}
	}()
	<-FIREFOX.Done()
	return nil
}
