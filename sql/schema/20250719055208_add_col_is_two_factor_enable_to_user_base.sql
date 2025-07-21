-- +goose Up
-- +goose StatementBegin
ALTER TABLE pre_go_acc_user_base_9999 
ADD COLUMN is_two_factor_enable int(1) NOT NULL DEFAULT 0 Comment 'Is two factor authentication enable';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
Alter Table pre_go_acc_user_base_9999
Drop COLUMN is_two_factor_enable;
-- +goose StatementEnd
