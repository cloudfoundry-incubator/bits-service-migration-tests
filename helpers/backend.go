package helpers

import (
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var DEFAULT_TIMEOUT = 30 * time.Second

func SetBackend(appName string) {
	config := helpers.LoadConfig()
	if config.Backend == "diego" {
		EnableDiego(appName)
	} else if config.Backend == "dea" {
		DisableDiego(appName)
	}
}

func EnableDiego(appName string) {
	guid := GetAppGuid(appName)
	Eventually(cf.Cf("curl", "/v2/apps/"+guid, "-X", "PUT", "-d", `{"diego": true}`), DEFAULT_TIMEOUT).Should(Exit(0))
}

func DisableDiego(appName string) {
	guid := GetAppGuid(appName)
	Eventually(cf.Cf("curl", "/v2/apps/"+guid, "-X", "PUT", "-d", `{"diego": false}`), DEFAULT_TIMEOUT).Should(Exit(0))
}

func GetAppGuid(appName string) string {
	cfApp := cf.Cf("app", appName, "--guid")
	Eventually(cfApp, DEFAULT_TIMEOUT).Should(Exit(0))

	appGuid := strings.TrimSpace(string(cfApp.Out.Contents()))
	Expect(appGuid).NotTo(Equal(""))
	return appGuid
}
