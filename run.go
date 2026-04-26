// Package goi2pbrowser is a package which can be used to manage an I2P browsing
// profile using a pre-configured, common profile which is used by the I2P Easy-Install
// bundle and the i2p.plugins.firefox profile manager. It is a Go clone of
// i2p.plugins.firefox for use in native applications.
package goi2pbrowser

// BrowseStrict launches a Firefox browser configured to use I2P and waits for it to exit.
// The profile is in "Strict" mode
func BrowseStrict(profileDir string, url ...string) error {
	i2pBrowser, err := NewI2PBrowser(profileDir)
	if err != nil {
		return err
	}
	i2pBrowser.BrowseStrict(url...)
	return nil
}

// BrowseUsability launches a Firefox browser configured to use I2P and waits for it to exit.
// The profile is in "Usability" mode
func BrowseUsability(profileDir string, url ...string) error {
	i2pBrowser, err := NewI2PBrowser(profileDir)
	if err != nil {
		return err
	}
	i2pBrowser.BrowseUsability(url...)
	return nil
}

// BrowseApp launches a Firefox browser configured to use I2P and waits for it to exit.
// The profile is in "Usability" mode
func BrowseApp(profileDir string, url ...string) error {
	i2pBrowser, err := NewI2PBrowser(profileDir)
	if err != nil {
		return err
	}
	i2pBrowser.BrowseUsability(url...)
	return nil
}
