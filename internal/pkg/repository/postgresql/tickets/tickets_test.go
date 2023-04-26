package tickets

import (
	"context"
	"github.com/stretchr/testify/assert"
	"homework-week-5/internal/pkg/db"
	"homework-week-5/internal/pkg/repository"
	"log"
	"testing"
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
	ticketRepo := NewTickets(database)

	defer database.GetPool(ctx).Close()

	tt := []struct {
		test_name string
		user_id   int
		cost      int
		place     int
		error     bool
	}{
		{
			test_name: "success",
			user_id:   1,
			cost:      200,
			place:     2325,
			error:     false,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.test_name, func(t *testing.T) {
			//act
			_, err := ticketRepo.Add(ctx, &repository.Ticket{
				UserID: tc.user_id,
				Cost:   tc.cost,
				Place:  tc.place,
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

func Test_Update(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	//arrange
	ticketRepo := NewTickets(database)

	defer database.GetPool(ctx).Close()

	id, _ := ticketRepo.Add(ctx, &repository.Ticket{
		UserID: 1,
		Cost:   138,
		Place:  232,
	})

	tt := []struct {
		test_name string
		id        int
		user_id   int
		cost      int
		place     int
		error     bool
	}{
		{
			test_name: "success",
			id:        int(id),
			user_id:   1,
			cost:      200,
			place:     2325,
			error:     false,
		},
		{
			test_name: "fail",
			id:        -1,
			user_id:   1,
			cost:      200,
			place:     2325,
			error:     false,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.test_name, func(t *testing.T) {
			//act
			_, err := ticketRepo.Update(ctx, &repository.Ticket{
				ID:     tc.id,
				UserID: tc.user_id,
				Cost:   tc.cost,
				Place:  tc.place,
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
	ticketRepo := NewTickets(database)

	defer database.GetPool(ctx).Close()

	id, _ := ticketRepo.Add(ctx, &repository.Ticket{
		UserID: 1,
		Cost:   138,
		Place:  232,
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
			error:     false,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.test_name, func(t *testing.T) {
			//act
			_, err := ticketRepo.GetById(ctx, tc.id)

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

func Test_GetByUserId(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	//arrange
	ticketRepo := NewTickets(database)

	defer database.GetPool(ctx).Close()

	id, _ := ticketRepo.Add(ctx, &repository.Ticket{
		UserID: 12,
		Cost:   138,
		Place:  232,
	})

	tt := []struct {
		test_name string
		id        int
		user_id   int
		error     bool
	}{
		{
			test_name: "success",
			id:        int(id),
			user_id:   12,
			error:     false,
		},
		{
			test_name: "fail",
			id:        int(id),
			user_id:   -1,
			error:     false,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.test_name, func(t *testing.T) {
			//act
			_, err := ticketRepo.GetByUserId(ctx, tc.user_id)

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
	ticketRepo := NewTickets(database)

	defer database.GetPool(ctx).Close()

	id, _ := ticketRepo.Add(ctx, &repository.Ticket{
		UserID: 1,
		Cost:   138,
		Place:  232,
	})

	tt := []struct {
		test_name string
		id        int
		error     bool
	}{
		{
			test_name: "fail",
			id:        -1,
			error:     false,
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
			_, err := ticketRepo.Delete(ctx, tc.id)

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
