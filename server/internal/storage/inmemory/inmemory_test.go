package inmemory

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"github.com/ykode/srp_demo/server/internal/domain"
	//"github.com/ykode/srp_demo/server/internal/query"
	"fmt"
	"testing"
)

const user1 = "some-user-name"
const user2 = "another-user-name"
const password = "some-weak-pass"
const salt_hex = "dd64796d747d9470d2d6645c0b00c0fead013d28c7aff8a089423c38774dc74825efbe5ba051ce7cfeb2e5cd199fb76d371a8a5eda857d71c44aeeaac115a4bd"
const v_hex = "dbfa60d2d08d981c6ec187cf201a8c0ad8b29e0f5dc2d2e8ff73c48fc5cfcf7f97fe62b7edd9929e2ed0e2d3a73a92cfd08ea7897e76065d92684a5a174a7ab882f5b2ffd0fd4b0fe53b0869ea93ba9e8d50fb70beb0c2aafef8cc9a74b7b2fde51c07b1362b88a0345628bee345cde4c15e4a59e4dc48394c5096ded921503"

func TestIdentityStorage(t *testing.T) {
	v, _ := hex.DecodeString(v_hex)
	salt, _ := hex.DecodeString(salt_hex)

	id1, _ := domain.NewIdentity(user1, salt, v)
	// id2, _ := domain.NewIdentity(user2, salt, v)

	t.Run("TestNonExisting", func(t *testing.T) {
		idStorage := NewInMemoryIdentityStorage()
		r := <-idStorage.FindIdentityByUserName(user1)

		assert.Nil(t, r.Result)
		assert.Error(t, ErrorNotFound, r.Err)
	})

	t.Run("TestInsertAndQuery", func(t *testing.T) {

		idStorage := NewInMemoryIdentityStorage()
		err := <-idStorage.SaveIdentity(id1)

		assert.NoError(t, err)

		r := <-idStorage.FindIdentityByUserName(user1)

		assert.NoError(t, r.Err)
		assert.IsType(t, &domain.Identity{}, r.Result)
		id := r.Result.(*domain.Identity)

		assert.Equal(t, *id1, *id)
	})

	t.Run("TestParallelInsertAndQuery", func(t *testing.T) {
		//idStorage := NewInMemoryIdentityStorage()
		var c1 chan error
		// var r1, r2 <-chan query.Result

		go func() {
			//<-idStorage.SaveIdentity(id1)
			c1 <- nil
		}()

		//go func() {
		//	c2 = idStorage.SaveIdentity(id2)
		//}()

		l := <-c1
		fmt.Println(l)
		/*
			go func() {
				<-c1

				r1 = idStorage.FindIdentityByUserName(user1)

			}()

			go func() {
				<-c2

				r2 = idStorage.FindIdentityByUserName(user2)
			}()

			idR1 := <-r1
			idR2 := <-r2

			assert.NoError(t, idR1.Err)
			assert.NoError(t, idR2.Err)

			assert.IsType(t, &domain.Identity{}, idR1.Result)
			assert.IsType(t, &domain.Identity{}, idR2.Result)
		*/
	})
}
