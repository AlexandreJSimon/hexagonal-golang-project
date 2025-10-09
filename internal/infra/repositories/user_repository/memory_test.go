package user_repository_test

import (
	"testing"

	domain "github.com/AlexandreJSimon/hexagonal-golang-project/internal/domain/user"
	"github.com/AlexandreJSimon/hexagonal-golang-project/internal/infra/repositories/user_repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Repository Test Suite")
}

var _ = Describe("Test suite for testing User Memory Repository behaviors", func() {

	var userRepo domain.UserRepository

	BeforeEach(func() {
		userRepo = user_repository.NewMemoryRepository()
	})

	Context("user operations in memory", func() {

		When("trying to save a user", func() {

			It("should save the user", func() {

				// Arrange

				user := domain.NewUser("John Doe", "johndoe", "johndoe@example.com", "password123")

				// Act

				err := userRepo.Save(user)
				userInMemory, _ := userRepo.GetByID(user.ID)

				// Assert

				Expect(err).To(BeNil())
				Expect(user).To(Equal(userInMemory))

				userRepo.Delete(user.ID)
			})
		})

		When("trying to get a user by ID", func() {

			It("should return the user", func() {

				// Arrange

				user := domain.NewUser("New John Doe", "newjohndoe", "newjohndoe@example.com", "newpassword123")
				userRepo.Save(user)

				// Act

				userInMemory, err := userRepo.GetByID(user.ID)

				// Assert

				Expect(err).To(BeNil())
				Expect(user).To(Equal(userInMemory))

				userRepo.Delete(user.ID)
			})
		})

		When("trying to update a user", func() {

			It("should update the user", func() {

				// Arrange

				user := domain.NewUser("Jane Smith", "janesmith", "janesmith@example.com", "securepass456")
				userRepo.Save(user)
				user.Name = "Jane Doe"
				user.Username = "janedoe"

				// Act

				err := userRepo.Update(user)
				userInMemory, _ := userRepo.GetByID(user.ID)

				// Assert

				Expect(err).To(BeNil())
				Expect(user).To(Equal(userInMemory))

				userRepo.Delete(user.ID)
			})
		})

		When("trying to delete a user", func() {

			It("should delete the user", func() {

				// Arrange

				user := domain.NewUser("Alice Johnson", "alicej", "alicej@example.com", "mypassword789")
				userRepo.Save(user)

				// Act

				err := userRepo.Delete(user.ID)
				userInMemory, errGet := userRepo.GetByID(user.ID)

				// Assert

				Expect(err).To(BeNil())
				Expect(userInMemory).To(BeNil())
				Expect(errGet).ToNot(BeNil())
				Expect(errGet.Error()).To(Equal("user not found"))
			})
		})

		When("trying to list users", func() {
			It("should return the list of users", func() {

				// Arrange

				var users []*domain.User

				users = append(users, domain.NewUser("Bob Brown", "bobb", "bobb@example.com", "randompass123"))
				userRepo.Save(users[0])

				users = append(users, domain.NewUser("Charlie Green", "charlieg", "charlieg@example.com", "securepass789"))
				userRepo.Save(users[1])

				// Act

				usersInMemory, err := userRepo.List(10, 0)
				count, _ := userRepo.Count()

				// Assert

				Expect(err).To(BeNil())
				Expect(usersInMemory).To(HaveLen(count))
				Expect(usersInMemory).To(Equal(users))

				userRepo.Delete(users[0].ID)
				userRepo.Delete(users[1].ID)
			})
		})

		When("trying to count users", func() {
			It("should return the count of users", func() {

				// Arrange

				var users []*domain.User

				users = append(users, domain.NewUser("Bob Brown", "bobb", "bobb@example.com", "randompass123"))
				userRepo.Save(users[0])

				users = append(users, domain.NewUser("Charlie Green", "charlieg", "charlieg@example.com", "securepass789"))
				userRepo.Save(users[1])

				users = append(users, domain.NewUser("Diana Prince", "dianap", "dianap@example.com", "wonderpass456"))
				userRepo.Save(users[2])

				// Act

				count, err := userRepo.Count()
				usersInMemory, _ := userRepo.List(10, 0)

				// Assert

				Expect(err).To(BeNil())
				Expect(count).To(Equal(len(usersInMemory)))
				Expect(count).To(Equal(3))

				userRepo.Delete(users[0].ID)
				userRepo.Delete(users[1].ID)
				userRepo.Delete(users[2].ID)
			})
		})
	})
})
