package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Response Test Suite")
}

var _ = Describe("API Response", func() {

	Context("Must maintain the message format", func() {

		When("generate a success response", func() {

			It("should return a success response with the correct format", func() {

				// Arrange

				w := httptest.NewRecorder()
				data := map[string]interface{}{
					"userName": "fakeUser",
					"email":    "fakeUser@email.com",
				}

				// Act

				Success(w, "User created successfully", data)

				var response SuccessResponse

				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					Fail("Failed to decode JSON response")
				}

				// Assert

				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(response.Status).To(Equal("success"))
				Expect(response.Message).To(Equal("User created successfully"))
				Expect(response.Data).To(Equal(data))
			})
		})

		When("generate a error response", func() {

			It("should return a error response with the correct format", func() {

				// Arrange

				w := httptest.NewRecorder()

				// Act

				Error(w, "Fail to create user", http.StatusBadRequest)

				var response SuccessResponse

				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					Fail("Failed to decode JSON response")
				}

				// Assert

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				Expect(response.Status).To(Equal("error"))
				Expect(response.Message).To(Equal("Fail to create user"))
			})
		})
	})
})
