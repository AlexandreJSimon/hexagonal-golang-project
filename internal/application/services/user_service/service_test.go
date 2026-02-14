package user_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/application/services/user_service"
	domain "github.com/AlexandreJSimon/hexagonal-golang-project/internal/domain"
	"github.com/AlexandreJSimon/hexagonal-golang-project/mocks"
	"go.uber.org/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Service Test Suite")
}

var _ = Describe("Test suite for testing User Service behaviors", func() {

	var userRepoMock *mocks.MockUserRepository
	var userService *user_service.UserService

	fakeName := "Fake User"
	fakeUsername := "fakeUser"
	fakeUserEmail := "fakeUser@email.com"
	fakeUserPassword := "fakePassword"

	BeforeEach(func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		userRepoMock = mocks.NewMockUserRepository(mockCtrl)
		userService = user_service.NewUserService(user_service.UserServiceInput{
			UserRepository: userRepoMock,
		})

	})

	Context("interactions with the user service", func() {

		When("create a user", func() {

			It("should save user with sucess", func() {

				// Arrange

				input := user_service.CreateUserInput{
					Name:     fakeName,
					Username: fakeUsername,
					Email:    fakeUserEmail,
					Password: fakeUserPassword,
				}

				var savedUser *domain.User
				userRepoMock.EXPECT().
					Save(gomock.Any()).
					DoAndReturn(func(u *domain.User) error {
						savedUser = u
						return nil
					})

				// Act

				_, err := userService.CreateUser(context.Background(), input)

				// Assert

				Expect(err).NotTo(HaveOccurred())

				Expect(savedUser.Name).To(Equal(input.Name))
				Expect(savedUser.Username).To(Equal(input.Username))
				Expect(savedUser.Email).To(Equal(input.Email))
				Expect(savedUser.Password).To(Equal(input.Password))
				Expect(savedUser.Role).To(Equal(domain.Viewer))
				Expect(savedUser.Status).To(Equal(domain.Active))

				Expect(savedUser.ID).NotTo(BeEmpty())
				Expect(savedUser.CreatedAt).NotTo(BeZero())

			})

			It("should return an error when saving user fails", func() {
				// Arrange

				input := user_service.CreateUserInput{
					Name:     fakeName,
					Username: fakeUsername,
					Email:    fakeUserEmail,
					Password: fakeUserPassword,
				}

				userRepoMock.EXPECT().
					Save(gomock.Any()).
					Return(errors.New("any error"))

				// Act

				_, err := userService.CreateUser(context.Background(), input)

				// Assert

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(errors.New("any error")))
			})
		})

		When("list users", func() {

			It("should return a list of users", func() {

				// Arrange

				fakeUsers := []*domain.User{
					domain.NewUser(fakeName, fakeUsername, fakeUserEmail, fakeUserPassword),
					domain.NewUser("Fake User 2", "fakeUser2", "fakeUser2@email.com", "fakePassword2"),
					domain.NewUser("Fake User 3", "fakeUser3", "fakeUser3@email.com", "fakePassword3"),
				}

				userRepoMock.EXPECT().
					List(10, 0).
					Return(fakeUsers, nil)

				// Act

				users, err := userService.ListUsers(context.Background(), 10, 0)

				// Assert

				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(Equal(fakeUsers))
			})

			It("should return an error when listing users fails", func() {

				// Arrange

				userRepoMock.EXPECT().
					List(10, 0).
					Return(nil, errors.New("any error"))

				// Act

				users, err := userService.ListUsers(context.Background(), 10, 0)

				// Assert

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(errors.New("any error")))
				Expect(users).To(BeNil())
			})
		})

		When("get a single user", func() {

			It("should return a user", func() {

				// Arrange

				fakeUser := domain.NewUser(fakeName, fakeUsername, fakeUserEmail, fakeUserPassword)
				userRepoMock.EXPECT().
					GetByID(fakeUser.ID).
					Return(fakeUser, nil)

				// Act

				user, err := userService.GetUserByID(context.Background(), fakeUser.ID)

				// Assert

				Expect(err).NotTo(HaveOccurred())
				Expect(user).To(Equal(fakeUser))
			})

			It("should return an error when user not found", func() {

				// Arrange

				userRepoMock.EXPECT().
					GetByID("non-existing-id").
					Return(nil, errors.New("user not found"))

				// Act

				user, err := userService.GetUserByID(context.Background(), "non-existing-id")

				// Assert

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(errors.New("user not found")))
				Expect(user).To(BeNil())
			})
		})

		When("update a user", func() {

			It("should update user with success", func() {

				// Arrange

				fakeUser := domain.NewUser(fakeName, fakeUsername, fakeUserEmail, fakeUserPassword)

				updateInput := user_service.UpdateUserInput{
					Name:     "Updated Name",
					Username: "updatedUsername",
					Email:    "updatedEmail@email.com",
					Password: "newPassword",
				}

				fakeUserUpdated := &domain.User{
					ID:        fakeUser.ID,
					Name:      updateInput.Name,
					Username:  updateInput.Username,
					Email:     updateInput.Email,
					Password:  updateInput.Password,
					Role:      fakeUser.Role,
					Status:    fakeUser.Status,
					CreatedAt: fakeUser.CreatedAt,
					UpdatedAt: fakeUser.UpdatedAt,
				}

				userRepoMock.EXPECT().
					GetByID(fakeUser.ID).
					Return(fakeUser, nil)

				userRepoMock.EXPECT().
					Update(fakeUserUpdated).
					Return(nil)

				// Act

				err := userService.UpdateUser(context.Background(), fakeUser.ID, updateInput)

				// Assert

				Expect(err).NotTo(HaveOccurred())
			})

			It("should return an error when user not found", func() {

				// Arrange

				updateInput := user_service.UpdateUserInput{
					Name:     "Updated Name",
					Username: "updatedUsername",
					Email:    "updatedEmail@email.com",
					Password: "newPassword",
				}

				userRepoMock.EXPECT().
					GetByID("non-existing-id").
					Return(nil, errors.New("user not found"))

				// Act

				err := userService.UpdateUser(context.Background(), "non-existing-id", updateInput)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(errors.New("user not found")))
			})

			It("should return an error when updating user fails", func() {

				// Arrange

				fakeUser := domain.NewUser(fakeName, fakeUsername, fakeUserEmail, fakeUserPassword)

				updateInput := user_service.UpdateUserInput{
					Name:     "Updated Name",
					Username: "updatedUsername",
					Email:    "updatedEmail@email.com",
					Password: "newPassword",
				}

				fakeUserUpdated := &domain.User{
					ID:        fakeUser.ID,
					Name:      updateInput.Name,
					Username:  updateInput.Username,
					Email:     updateInput.Email,
					Password:  updateInput.Password,
					Role:      fakeUser.Role,
					Status:    fakeUser.Status,
					CreatedAt: fakeUser.CreatedAt,
					UpdatedAt: fakeUser.UpdatedAt,
				}

				userRepoMock.EXPECT().
					GetByID(fakeUser.ID).
					Return(fakeUser, nil)

				userRepoMock.EXPECT().
					Update(fakeUserUpdated).
					Return(errors.New("any error"))

				// Act

				err := userService.UpdateUser(context.Background(), fakeUser.ID, updateInput)

				// Assert

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(errors.New("any error")))
			})
		})

		When("delete a user", func() {

			It("should delete user with success", func() {

				// Arrange

				fakeUser := domain.NewUser(fakeName, fakeUsername, fakeUserEmail, fakeUserPassword)

				userRepoMock.EXPECT().
					Delete(fakeUser.ID).
					Return(nil)

				// Act

				err := userService.DeleteUser(context.Background(), fakeUser.ID)

				// Assert

				Expect(err).NotTo(HaveOccurred())
			})

			It("should return an error when user not found", func() {

				// Arrange

				userRepoMock.EXPECT().
					Delete("non-existing-id").
					Return(errors.New("user not found"))

				// Act

				err := userService.DeleteUser(context.Background(), "non-existing-id")

				// Assert

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(errors.New("user not found")))
			})
		})

		When("count users", func() {

			It("should return the count of users", func() {

				// Arrange

				expectedCount := 5
				userRepoMock.EXPECT().
					Count().
					Return(expectedCount, nil)

				// Act

				count, err := userService.CountUsers(context.Background())

				// Assert

				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(expectedCount))
			})

			It("should return an error when counting users fails", func() {

				// Arrange

				userRepoMock.EXPECT().
					Count().
					Return(0, errors.New("any error"))

				// Act

				count, err := userService.CountUsers(context.Background())

				// Assert

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(errors.New("any error")))
				Expect(count).To(Equal(0))
			})
		})
	})
})
