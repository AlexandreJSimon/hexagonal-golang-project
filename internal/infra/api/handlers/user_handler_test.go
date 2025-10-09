package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/application/services/user_service"
	domain "github.com/AlexandreJSimon/hexagonal-golang-project/internal/domain/user"
	"github.com/AlexandreJSimon/hexagonal-golang-project/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("User API Handler", func() {

	var userServiceMock *mocks.MockUserServiceProvider

	fakeName := "Fake User"
	fakeUsername := "fakeUser"
	fakeUserEmail := "fakeUser@email.com"
	fakeUserPassword := "fakePassword"

	var handlers *Handler

	BeforeEach(func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		userServiceMock = mocks.NewMockUserServiceProvider(mockCtrl)

		handlers = NewHandler(HandlerInput{userServiceMock})
	})

	Context("User API Handler", func() {

		Context("Create User", func() {

			When("pass all parameters", func() {

				It("should create a user successfully", func() {

					// Arrange

					data, err := json.Marshal(map[string]interface{}{
						"name":     fakeName,
						"userName": fakeUsername,
						"email":    fakeUserEmail,
						"password": fakeUserPassword,
					})
					if err != nil {
						Fail("Failed to encode JSON")
					}

					rr := httptest.NewRecorder()
					req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(data))
					req.Header.Set("Content-Type", "application/json")

					userServiceMock.EXPECT().CreateUser(req.Context(), user_service.CreateUserInput{
						Name:     fakeName,
						Username: fakeUsername,
						Email:    fakeUserEmail,
						Password: fakeUserPassword,
					}).Return("12345", nil)

					// Act

					handlers.CreateUser(rr, req)

					// Assert

					Expect(rr.Code).To(Equal(http.StatusOK))
					Expect(rr.Header().Get("Content-Type")).To(Equal("application/json"))
					Expect(rr.Body.String()).To(MatchJSON(
						`{
						"status": "success",
						"message": "User created successfully",
						"data": {
							"id": "12345",
							"name": "Fake User",
							"username": "fakeUser",
							"email": "fakeUser@email.com"
						}
					}`,
					))
				})
			})

			When("there was a failure in creating it", func() {

				It("should return an error", func() {

					// Arrange

					data, err := json.Marshal(map[string]interface{}{
						"name":     fakeName,
						"userName": fakeUsername,
						"email":    fakeUserEmail,
						"password": fakeUserPassword,
					})
					if err != nil {
						Fail("Failed to encode JSON")
					}

					rr := httptest.NewRecorder()

					req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(data))
					req.Header.Set("Content-Type", "application/json")

					userServiceMock.EXPECT().CreateUser(req.Context(), user_service.CreateUserInput{
						Name:     fakeName,
						Username: fakeUsername,
						Email:    fakeUserEmail,
						Password: fakeUserPassword,
					}).Return("", errors.New("any error"))

					// Act

					handlers.CreateUser(rr, req)

					// Assert

					Expect(rr.Code).To(Equal(http.StatusInternalServerError))
					Expect(rr.Header().Get("Content-Type")).To(Equal("application/json"))
					Expect(rr.Body.String()).To(MatchJSON(
						`{
						"status": "error",
						"message": "Failed to create user",
						"code": 500
					}`,
					))
				})
			})
		})

		Context("List Users", func() {

			When("user list is requested", func() {

				It("should list users successfully", func() {

					// Arrange

					rr := httptest.NewRecorder()
					req := httptest.NewRequest(http.MethodGet, "/users?limit=10&offset=0", nil)

					userServiceMock.EXPECT().ListUsers(req.Context(), 10, 0).Return([]*domain.User{
						{
							ID:       "12345",
							Name:     fakeName,
							Username: fakeUsername,
							Email:    fakeUserEmail,
						},
						{
							ID:       "54321",
							Name:     "Another User",
							Username: "anotherUser",
							Email:    "anotherUser@email.com",
						},
					}, nil)

					// Act

					handlers.ListUsers(rr, req)

					// Assert

					Expect(rr.Code).To(Equal(http.StatusOK))
					Expect(rr.Body.String()).To(MatchJSON(
						`{
						"status": "success",
						"message": "Users retrieved successfully",
						"data": [
							{
								"id": "12345",
								"name": "Fake User",
								"username": "fakeUser",
								"email": "fakeUser@email.com"
							},
							{
								"id": "54321",
								"name": "Another User",
								"username": "anotherUser",
								"email": "anotherUser@email.com"
							}
						]
					}`,
					))
				})
			})

			When("user list is requested and a problem occurs", func() {

				It("should return an error", func() {

					// Arrange

					rr := httptest.NewRecorder()
					req := httptest.NewRequest(http.MethodGet, "/users?limit=10&offset=0", nil)

					userServiceMock.EXPECT().ListUsers(req.Context(), 10, 0).Return(nil, errors.New("any error"))

					// Act

					handlers.ListUsers(rr, req)

					// Assert

					Expect(rr.Code).To(Equal(http.StatusInternalServerError))
					Expect(rr.Header().Get("Content-Type")).To(Equal("application/json"))
					Expect(rr.Body.String()).To(MatchJSON(
						`{
						"status": "error",
						"message": "Failed to list users",
						"code": 500
					}`,
					))
				})
			})
		})
	})
})
