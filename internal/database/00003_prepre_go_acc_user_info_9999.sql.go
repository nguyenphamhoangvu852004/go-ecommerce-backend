// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: 00003_prepre_go_acc_user_info_9999.sql

package database

import (
	"context"
	"database/sql"
)

const addUserAutoUserId = `-- name: AddUserAutoUserId :execresult

INSERT INTO
  pre_go_acc_user_info_9999 (
    user_account,
    user_nickname,
    user_avatar,
    user_state,
    user_mobile,
    user_gender,
    user_birthday,
    user_email,
    user_is_authentication
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type AddUserAutoUserIdParams struct {
	UserAccount          string
	UserNickname         sql.NullString
	UserAvatar           sql.NullString
	UserState            uint8
	UserMobile           sql.NullString
	UserGender           sql.NullInt16
	UserBirthday         sql.NullTime
	UserEmail            sql.NullString
	UserIsAuthentication uint8
}

// -- name: UpdatePassword :exec
// UPDATE pre_go_acc_user_info_9999
// SET
//
//	user_password = ?
//
// WHERE
//
//	user_id = ?;
func (q *Queries) AddUserAutoUserId(ctx context.Context, arg AddUserAutoUserIdParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, addUserAutoUserId,
		arg.UserAccount,
		arg.UserNickname,
		arg.UserAvatar,
		arg.UserState,
		arg.UserMobile,
		arg.UserGender,
		arg.UserBirthday,
		arg.UserEmail,
		arg.UserIsAuthentication,
	)
}

const addUserHaveUserId = `-- name: AddUserHaveUserId :execresult
INSERT INTO
  pre_go_acc_user_info_9999 (
    user_id,
    user_account,
    user_nickname,
    user_avatar,
    user_state,
    user_mobile,
    user_gender,
    user_birthday,
    user_email,
    user_is_authentication
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type AddUserHaveUserIdParams struct {
	UserID               uint64
	UserAccount          string
	UserNickname         sql.NullString
	UserAvatar           sql.NullString
	UserState            uint8
	UserMobile           sql.NullString
	UserGender           sql.NullInt16
	UserBirthday         sql.NullTime
	UserEmail            sql.NullString
	UserIsAuthentication uint8
}

func (q *Queries) AddUserHaveUserId(ctx context.Context, arg AddUserHaveUserIdParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, addUserHaveUserId,
		arg.UserID,
		arg.UserAccount,
		arg.UserNickname,
		arg.UserAvatar,
		arg.UserState,
		arg.UserMobile,
		arg.UserGender,
		arg.UserBirthday,
		arg.UserEmail,
		arg.UserIsAuthentication,
	)
}

const editUserByUserId = `-- name: EditUserByUserId :execresult
UPDATE pre_go_acc_user_info_9999
SET
  user_nickname = ?,
  user_avatar = ?,
  user_mobile = ?,
  user_gender = ?,
  user_birthday = ?,
  user_email = ?,
  updated_at = NOW()
WHERE
  user_id = ?
  AND user_is_authentication = 1
`

type EditUserByUserIdParams struct {
	UserNickname sql.NullString
	UserAvatar   sql.NullString
	UserMobile   sql.NullString
	UserGender   sql.NullInt16
	UserBirthday sql.NullTime
	UserEmail    sql.NullString
	UserID       uint64
}

func (q *Queries) EditUserByUserId(ctx context.Context, arg EditUserByUserIdParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, editUserByUserId,
		arg.UserNickname,
		arg.UserAvatar,
		arg.UserMobile,
		arg.UserGender,
		arg.UserBirthday,
		arg.UserEmail,
		arg.UserID,
	)
}

const getUser = `-- name: GetUser :one
SELECT
  user_id,
  user_account,
  user_nickname,
  user_avatar,
  user_state,
  user_mobile,
  user_gender,
  user_birthday,
  user_email,
  user_is_authentication,
  created_at,
  updated_at
FROM
  pre_go_acc_user_info_9999
WHERE
  user_id = ?
LIMIT
  1
`

func (q *Queries) GetUser(ctx context.Context, userID uint64) (PreGoAccUserInfo9999, error) {
	row := q.db.QueryRowContext(ctx, getUser, userID)
	var i PreGoAccUserInfo9999
	err := row.Scan(
		&i.UserID,
		&i.UserAccount,
		&i.UserNickname,
		&i.UserAvatar,
		&i.UserState,
		&i.UserMobile,
		&i.UserGender,
		&i.UserBirthday,
		&i.UserEmail,
		&i.UserIsAuthentication,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByAccount = `-- name: GetUserByAccount :one
SELECT
  user_id,
  user_account,
  user_nickname,
  user_avatar,
  user_state,
  user_mobile,
  user_gender,
  user_birthday,
  user_email,
  user_is_authentication,
  created_at,
  updated_at
FROM
  pre_go_acc_user_info_9999
WHERE
  user_id IN (?)
`

func (q *Queries) GetUserByAccount(ctx context.Context, userID uint64) (PreGoAccUserInfo9999, error) {
	row := q.db.QueryRowContext(ctx, getUserByAccount, userID)
	var i PreGoAccUserInfo9999
	err := row.Scan(
		&i.UserID,
		&i.UserAccount,
		&i.UserNickname,
		&i.UserAvatar,
		&i.UserState,
		&i.UserMobile,
		&i.UserGender,
		&i.UserBirthday,
		&i.UserEmail,
		&i.UserIsAuthentication,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT
  user_id, user_account, user_nickname, user_avatar, user_state, user_mobile, user_gender, user_birthday, user_email, user_is_authentication, created_at, updated_at
FROM
  pre_go_acc_user_info_9999
WHERE
  user_account LIKE ?
  OR user_nickname LIKE ?
`

type ListUsersParams struct {
	UserAccount  string
	UserNickname sql.NullString
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]PreGoAccUserInfo9999, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.UserAccount, arg.UserNickname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PreGoAccUserInfo9999
	for rows.Next() {
		var i PreGoAccUserInfo9999
		if err := rows.Scan(
			&i.UserID,
			&i.UserAccount,
			&i.UserNickname,
			&i.UserAvatar,
			&i.UserState,
			&i.UserMobile,
			&i.UserGender,
			&i.UserBirthday,
			&i.UserEmail,
			&i.UserIsAuthentication,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsersLimit = `-- name: ListUsersLimit :many
SELECT
  user_id, user_account, user_nickname, user_avatar, user_state, user_mobile, user_gender, user_birthday, user_email, user_is_authentication, created_at, updated_at
FROM
  pre_go_acc_user_info_9999
LIMIT
  ?
OFFSET
  ?
`

type ListUsersLimitParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) ListUsersLimit(ctx context.Context, arg ListUsersLimitParams) ([]PreGoAccUserInfo9999, error) {
	rows, err := q.db.QueryContext(ctx, listUsersLimit, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PreGoAccUserInfo9999
	for rows.Next() {
		var i PreGoAccUserInfo9999
		if err := rows.Scan(
			&i.UserID,
			&i.UserAccount,
			&i.UserNickname,
			&i.UserAvatar,
			&i.UserState,
			&i.UserMobile,
			&i.UserGender,
			&i.UserBirthday,
			&i.UserEmail,
			&i.UserIsAuthentication,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeUser = `-- name: RemoveUser :exec
DELETE FROM pre_go_acc_user_info_9999
WHERE
  user_id = ?
`

func (q *Queries) RemoveUser(ctx context.Context, userID uint64) error {
	_, err := q.db.ExecContext(ctx, removeUser, userID)
	return err
}
