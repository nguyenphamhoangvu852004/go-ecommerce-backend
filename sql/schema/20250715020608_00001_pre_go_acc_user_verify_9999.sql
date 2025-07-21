-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_acc_user_verify_9999 (
    verify_id int NOT NULL AUTO_INCREMENT,
    verify_otp varchar(6) NOT NULL,
    verify_key varchar(255) NOT NULL,
    verify_key_hash varchar(255) NOT NULL,
    verify_type int DEFAULT '1',
    is_verified int DEFAULT '0',
    is_deleted int DEFAULT '0',
    verify_created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    verify_updated_at timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (verify_id),
    UNIQUE KEY unique_verify_key (verify_key),
    KEY idx_verify_otp (verify_otp)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='pre_go_acc_user_verify_9999';

-- SET FOREIGN_KEY_CHECKS = 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `pre_go_acc_user_verify_9999`;
-- +goose StatementEnd
