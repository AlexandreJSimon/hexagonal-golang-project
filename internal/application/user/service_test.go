package user_test

import (
	"context"
	"errors"
	"testing"

	userApp "github.com/AlexandreJSimon/hexagonal-golang-project/internal/application/user"
	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/domain/user"
	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/security"
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
	var userService *userApp.UserService
	var hasher = security.BcryptHasher{}

	fakeName := "Fake User"
	fakeUsername := "fakeUser"
	fakeUserEmail := "fakeUser@email.com"
	fakeUserPassword := "fakePassword"
	fakeUserPasswordHashed, _ := hasher.Hash(fakeUserPassword)

	BeforeEach(func() {
		mockCtrl := gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()

		userRepoMock = mocks.NewMockUserRepository(mockCtrl)
		userService = userApp.NewUserService(userApp.UserServiceInput{
			UserRepository: userRepoMock,
			PasswordHasher: hasher,
		})

	})

	Context("interactions with the user service", func() {

		When("create a user", func() {

			It("should save user with sucess", func() {

				// Arrange

				input := userApp.CreateUserInput{
					Name:     fakeName,
					Username: fakeUsername,
					Email:    fakeUserEmail,
					Password: fakeUserPassword,
				}

				var savedUser *user.User
				userRepoMock.EXPECT().
					Save(gomock.Any()).
					DoAndReturn(func(u *user.User) error {
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
				Expect(hasher.Compare(fakeUserPassword, savedUser.Password)).To(BeTrue())
				Expect(savedUser.Role).To(Equal(user.Viewer))
				Expect(savedUser.Status).To(Equal(user.Active))

				Expect(savedUser.ID).NotTo(BeEmpty())
				Expect(savedUser.CreatedAt).NotTo(BeZero())

			})

			It("should return an error when saving user fails", func() {
				// Arrange

				input := userApp.CreateUserInput{
					Name:     fakeName,
					Username: fakeUsername,
					Email:    fakeUserEmail,
					Password: fakeUserPasswordHashed,
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

				fakeUserPasswordHashed2, _ := hasher.Hash("fakePassword2")
				fakeUserPasswordHashed3, _ := hasher.Hash("fakePassword3")

				fakeUsers := []*user.User{
					user.NewUser(fakeName, fakeUsername, fakeUserEmail, fakeUserPasswordHashed),
					user.NewUser("Fake User 2", "fakeUser2", "fakeUser2@email.com", fakeUserPasswordHashed2),
					user.NewUser("Fake User 3", "fakeUser3", "fakeUser3@email.com", fakeUserPasswordHashed3),
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

				fakeUser := user.NewUser(fakeName, fakeUsername, fakeUserEmail, fakeUserPasswordHashed)
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

				fakeUser := user.NewUser(fakeName, fakeUsername, fakeUserEmail, fakeUserPasswordHashed)

				updateInput := userApp.UpdateUserInput{
					Name:     "Updated Name",
					Username: "updatedUsername",
					Email:    "updatedEmail@email.com",
					Password: "newPassword",
				}

				userRepoMock.EXPECT().
					GetByID(fakeUser.ID).
					Return(fakeUser, nil)

				var savedUser *user.User
				userRepoMock.EXPECT().
					Update(gomock.Any()).
					DoAndReturn(func(user *user.User) error {
						savedUser = user
						return nil
					})
				// Act

				err := userService.UpdateUser(context.Background(), fakeUser.ID, updateInput)

				// Assert

				Expect(err).NotTo(HaveOccurred())

				Expect(savedUser.Name).To(Equal(updateInput.Name))
				Expect(savedUser.Username).To(Equal(updateInput.Username))
				Expect(savedUser.Email).To(Equal(updateInput.Email))
				Expect(hasher.Compare("newPassword", savedUser.Password)).To(BeTrue())
				Expect(savedUser.Role).To(Equal(user.Viewer))
				Expect(savedUser.Status).To(Equal(user.Active))

				Expect(savedUser.ID).NotTo(BeEmpty())
				Expect(savedUser.CreatedAt).NotTo(BeZero())

			})

			It("should return an error when user not found", func() {

				// Arrange

				updateInput := userApp.UpdateUserInput{
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

				fakeUser := user.NewUser(fakeName, fakeUsername, fakeUserEmail, fakeUserPassword)

				updateInput := userApp.UpdateUserInput{
					Name:     "Updated Name",
					Username: "updatedUsername",
					Email:    "updatedEmail@email.com",
					Password: "newPassword",
				}

				userRepoMock.EXPECT().
					GetByID(fakeUser.ID).
					Return(fakeUser, nil)

				userRepoMock.EXPECT().
					Update(gomock.Any()).
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

				fakeUser := user.NewUser(fakeName, fakeUsername, fakeUserEmail, fakeUserPassword)

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
