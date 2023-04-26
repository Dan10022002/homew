package users

import (
	"context"
	"log"
	"testing"

	"homework-week-5/internal/pkg/db"
	"homework-week-5/internal/pkg/repository"

	"github.com/stretchr/testify/assert"
)

func Test_Add(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	//arrange
	userRepo := NewUsers(database)

	defer database.GetPool(ctx).Close()

	tt := []struct {
		test_name string
		name      string
		surname   string
		age       int
		error     bool
	}{
		{
			test_name: "success",
			name:      "Ivan",
			surname:   "Ivanov",
			age:       25,
			error:     false,
		},
		{
			test_name: "fail",
			name:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab",
			surname:   "Ivanov",
			age:       25,
			error:     true,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.test_name, func(t *testing.T) {
			//act
			_, err := userRepo.Add(ctx, &repository.User{
				Name:    tc.name,
				Surname: tc.surname,
				Age:     tc.age,
			})

			//assert
			if tc.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
	database.Truncate(ctx)
}

func Test_GetById(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	//arrange
	userRepo := NewUsers(database)

	defer database.GetPool(ctx).Close()

	id, _ := userRepo.Add(ctx, &repository.User{
		Name:    "Ivan",
		Surname: "Ivanov",
		Age:     25,
	})

	tt := []struct {
		test_name string
		id        int
		error     bool
	}{
		{
			test_name: "success",
			id:        int(id),
			error:     false,
		},
		{
			test_name: "fail",
			id:        -1,
			error:     true,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.test_name, func(t *testing.T) {
			//act
			_, err := userRepo.GetById(ctx, tc.id)

			//assert
			if tc.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
	database.Truncate(ctx)
}

func Test_Update(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	//arrange
	userRepo := NewUsers(database)

	defer database.GetPool(ctx).Close()

	id, _ := userRepo.Add(ctx, &repository.User{
		Name:    "Ivan",
		Surname: "Ivanov",
		Age:     25,
	})

	tt := []struct {
		test_name string
		id        int
		name      string
		surname   string
		age       int
		error     bool
	}{
		{
			test_name: "success",
			id:        int(id),
			name:      "Semen",
			surname:   "Semenov",
			age:       30,
			error:     false,
		},
		{
			test_name: "fail",
			id:        int(id),
			name:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab",
			surname:   "Semenov",
			age:       30,
			error:     true,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.test_name, func(t *testing.T) {
			//act
			_, err := userRepo.Update(ctx, &repository.User{
				ID:      tc.id,
				Name:    tc.name,
				Surname: tc.surname,
				Age:     tc.age,
			})

			//assert
			if tc.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
	database.Truncate(ctx)
}

func Test_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	//arrange
	userRepo := NewUsers(database)

	defer database.GetPool(ctx).Close()

	id, _ := userRepo.Add(ctx, &repository.User{
		Name:    "Ivan",
		Surname: "Ivanov",
		Age:     25,
	})

	tt := []struct {
		test_name string
		id        int
		error     bool
	}{
		{
			test_name: "fail",
			id:        -1,
			error:     true,
		},
		{
			test_name: "success",
			id:        int(id),
			error:     false,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.test_name, func(t *testing.T) {
			//act
			_, err := userRepo.Delete(ctx, tc.id)

			//assert
			if tc.error {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
	database.Truncate(ctx)
}
