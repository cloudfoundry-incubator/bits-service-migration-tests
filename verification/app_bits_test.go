package verification_test

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	. "github.com/cloudfoundry-incubator/bits-service-migration-tests/helpers"
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("AppBits", func() {
	It("uses previously cached app_bits", func() {
		resp, err := makeAuthorizedCfRequest("GET", "https://"+config.ApiEndpoint+"/v2/apps/"+GetAppGuid(AppBitsAppName)+"/download")
		Expect(err).ToNot(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		tempFile, err := ioutil.TempFile("", "")
		Expect(err).ToNot(HaveOccurred())
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		_, err = io.Copy(tempFile, resp.Body)
		Expect(err).ToNot(HaveOccurred())

		tmpDir, err := ioutil.TempDir("", "unziped-app-path")
		Expect(err).ToNot(HaveOccurred())
		Unzip(tempFile.Name(), tmpDir)
		resourceMatchBody := string(ResourceMatchBody(tmpDir))

		resourceMatches := cf.Cf("curl", "-X", "PUT", "/v2/resource_match", "-d", resourceMatchBody).Wait(defaultTimeout)

		Expect(resourceMatches).To(Exit(0))
		Expect(resourceMatches).To(Say(resourceMatchBody))
	})
})

// Not using cf curl here, because it doesn't do more than 1 redirect.
// But the newest bits-service-client requires more than 1 redirect in most cases.
func makeAuthorizedCfRequest(method string, url string) (*http.Response, error) {
	curlOAuth := cf.Cf("oauth-token").Wait(defaultTimeout)
	Expect(curlOAuth.ExitCode()).To(Equal(0))
	bearerToken := strings.TrimRight(string(curlOAuth.Out.Contents()), "\n")

	r, err := http.NewRequest(method, url, nil)
	Expect(err).ToNot(HaveOccurred())
	r.Header.Set("Authorization", bearerToken)

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	return httpClient.Do(r)
}
