package github

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
)

const UriUser = "https://api.github.com/user"

type UserPublicProfile struct {
	Login             string    `json:"login"`
	ID                int64     `json:"id"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GravatarID        string    `json:"gravatar_id"`
	URL               string    `json:"url"`
	HtmlURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             string    `json:"email"`
	Hireable          bool      `json:"hireable"`
	Bio               string    `json:"bio"`
	TwitterUsername   string    `json:"twitter_username"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func FetchGithubUserPublicProfile(token string) (UserPublicProfile, error) {
	write := func(err error) (UserPublicProfile, error) {
		return UserPublicProfile{}, err
	}

	req, err := http.NewRequest(http.MethodGet, UriUser, strings.NewReader(""))
	if err != nil {
		return write(util.NewErrorWithUnknown("fails to new request"))
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return write(util.NewErrorWithUnknown("fails to fetch github user public profile"))
	}

	if res.StatusCode == http.StatusForbidden {
		// TODO[feat]: should we inactive token?
		return write(util.NewError(errors.ErrThirdPartyForbidden, "github access token has been invalid"))
	}

	if res.StatusCode != http.StatusOK {
		return write(util.NewError(errors.ErrThirdPartyUnknown, "fails to fetch github user public profile"))
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return write(util.NewErrorWithUnknown("fails to fetch github user public profile"))
	}

	payload := UserPublicProfile{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return write(util.NewErrorWithUnknown("fails to parse github user public profile"))
	}

	return payload, nil
}
