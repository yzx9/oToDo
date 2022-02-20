# oToDo API

This package contain apis of oToDo.

All APIs require authentication unless otherwise specified, we are
using `access_token` and `refresh_token` for authentication. You needs
to login to obtain tokens, then include authorization in your `Request.Headers`:

```http
Authorization: Bearer <access_token>
```

And safely save `refresh_token`, e.g.: `localStroage`

## Session

### Login

> Login with password

POST /api/sessions

#### Request

| Param    | Type   | Description |
| -------- | ------ | ----------- |
| userName | String | "admin"     |
| password | String | "admin123"  |

#### Response

| Param        | Type   | Description                     |
| ------------ | ------ | ------------------------------- |
| accessToken  | String | Json web token                  |
| expiresIn    | Int    | Access token expiration seconds |
| tokenType    | String | Only `bearer` now               |
| refreshToken | String | Should be save SAFELY           |

#### Remark

- This is a public api

### Logout

> Logout and unactive refresh token

DELETE /api/sessions

#### Request

| Param | Type | Description |
| ----- | ---- | ----------- |
| -     | -    | -           |

#### Response

| Param   | Type   | Description |
| ------- | ------ | ----------- |
| message | String | "see you"   |

### New Access Token (Active)

> Get new access token by refresh token

POST /api/sessions/current/tokens

#### Request

| Param        | Type   | Description          |
| ------------ | ------ | -------------------- |
| refreshToken | String | From [Login](#Login) |

#### Response

| Param       | Type   | Description         |
| ----------- | ------ | ------------------- |
| accessToken | String | See [Login](#Login) |
| expiresIn   | Int    | See [Login](#Login) |
| tokenType   | String | See [Login](#Login) |

### New Access Token (Passive)

POST /api/\*

#### Details

The system may include an new access token in your non-public request,
check the response headers if there exists `Authorization`, you can
update your local token so that we dont need to refresh token
frequently. The smuggled token still follows your request format:

```http
Authorization: Bearer <access_token>
```

### Test Access Token

> Test your access token, also for passsive refresh token via timer

GET /api/sessions

### Github OAuth Creater

> Login by github, generate an redirect uri

GET /api/sessions/oauth/github

#### Request

| Param | Type | Description |
| ----- | ---- | ----------- |
| -     | -    | -           |

#### Response

| Param       | Type   | Description      |
| ----------- | ------ | ---------------- |
| redirectURI | String | Redirect user to |

### Github OAuth Creater

> Login by github, send code to server after user authentized

POST /api/sessions/oauth/github

#### Request

| Param | Type   | Description |
| ----- | ------ | ----------- |
| code  | String | -           |
| state | String | -           |

#### Response

| Param        | Type   | Description         |
| ------------ | ------ | ------------------- |
| accessToken  | String | See [Login](#Login) |
| expiresIn    | Int    | See [Login](#Login) |
| tokenType    | String | See [Login](#Login) |
| refreshToken | String | See [Login](#Login) |

#### Remark

- Github should take params `code` and `state`, simply forward to server
- You should reject the access without `code` and `state`, such as redirect to welcome page

## User

## Todo

## Todo List

## Todo List Folder

## Sharing

## File
