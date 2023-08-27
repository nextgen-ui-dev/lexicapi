package friend

import "gopkg.in/guregu/null.v4"

type FriendDetail struct {
	Friend
	Name     null.String `json:"name"`
	Email    null.String `json:"email"`
	ImageUrl null.String `json:"image_url"`
}
