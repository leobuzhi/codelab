package main

import (
	"fmt"
	"testing"

	"github.com/leobuzhi/codelab/friends/model"
	"github.com/stretchr/testify/assert"
)

func findLikePerson(where string) []model.LikePerson {
	var ret []model.LikePerson
	db.Find(&ret, where)
	return ret
}

func findFriend(where string) []model.Friend {
	var ret []model.Friend
	db.Find(&ret, where)
	return ret
}

type params struct {
	person1 int
	person2 int
}

func TestAttention(t *testing.T) {
	tcs := []struct {
		params      []params
		err         error
		likePersons []model.LikePerson
		friends     []model.Friend
	}{
		{
			[]params{
				params{1000, 1000},
			},
			fmt.Errorf("userID1 and userID2 should not be equal, userID1: %d, userID2: %d", 1000, 1000),
			[]model.LikePerson{
				{ID: 1, UserID: 1000, LikerID: 2000, RelationShip: 3},
			},
			[]model.Friend{
				{ID: 1, Person1: 1000, Person2: 2000},
			},
		},
		{
			[]params{
				params{1000, 2000},
				params{1000, 2000},
			},
			nil,
			[]model.LikePerson{
				{ID: 1, UserID: 1000, LikerID: 2000, RelationShip: 1},
			},
			[]model.Friend{},
		},
		{
			[]params{
				params{2000, 1000},
				params{2000, 1000},
			},
			nil,
			[]model.LikePerson{
				{ID: 1, UserID: 1000, LikerID: 2000, RelationShip: 2},
			},
			[]model.Friend{},
		},
		{
			[]params{
				params{1000, 2000},
				params{2000, 1000},
				params{1000, 3000},
			},
			nil,
			[]model.LikePerson{
				{ID: 1, UserID: 1000, LikerID: 2000, RelationShip: 3},
				// NOTE(leobuzhi): ID is 3
				{ID: 3, UserID: 1000, LikerID: 3000, RelationShip: 1},
			},
			[]model.Friend{
				{ID: 1, Person1: 1000, Person2: 2000},
			},
		},
		{
			[]params{
				params{1000, 2000},
				params{2000, 1000},
			},
			nil,
			[]model.LikePerson{
				{ID: 1, UserID: 1000, LikerID: 2000, RelationShip: 3},
			},
			[]model.Friend{
				{ID: 1, Person1: 1000, Person2: 2000},
			},
		},
	}

	for _, tc := range tcs {
		setupDB("testdb")
		for _, param := range tc.params {
			assert.Equal(t, tc.err, attention(param.person1, param.person2))
		}

		if tc.err == nil {
			assert.Equal(t, tc.likePersons, findLikePerson(""))
			assert.Equal(t, tc.friends, findFriend(""))
		}
		teardown()
	}
}
