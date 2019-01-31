package verification_test

import (
	"net/http"

	. "github.com/cloudfoundry-incubator/bits-service-migration-tests/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verification/Droplets", func() {
	It("finds the previously created droplet", func() {
		resp, err := makeAuthorizedCfRequest("GET", "https://"+config.ApiEndpoint+"/v2/apps/"+GetAppGuid(DropletTestAppName)+"/droplet/download")
		Expect(resp.StatusCode, err).To(Equal(http.StatusOK))
	})
})
