-- oauth2_scope
INSERT INTO `oauth2_scope` (`id`, `created_at`, `updated_at`, `name`, `description`, `code`, `type`) VALUES ('1', '2020-06-29 12:28:59', '2020-06-29 12:28:59', 'openid', '用户ID', 'openid', 'basic');
INSERT INTO `oauth2_scope` (`id`, `created_at`, `updated_at`, `name`, `description`, `code`, `type`) VALUES ('2', '2020-06-29 12:28:59', '2020-06-29 12:28:59', 'profile', '用户资料', 'profile', 'basic');
INSERT INTO `oauth2_scope` (`id`, `created_at`, `updated_at`, `name`, `description`, `code`, `type`) VALUES ('3', '2020-06-29 12:28:59', '2020-06-29 12:28:59', 'email', '用户邮箱', 'email', 'basic');
INSERT INTO `oauth2_scope` (`id`, `created_at`, `updated_at`, `name`, `description`, `code`, `type`) VALUES ('4', '2020-06-29 12:28:59', '2020-06-29 12:28:59', 'phone', '用户手机号', 'phone', 'basic');
-- oauth2_client
INSERT INTO `oauth2_client` (`client_id`, `client_secret`, `redirect_uri`) VALUES (1000, "99799a6b-a289-4099-b4ad-b42603c17ffc", "http://localhost:9000/auth/callback");
-- oauth2_client_info
INSERT INTO `oauth2_client_info` (`client_id`, `name`, `profile`, `description`, `website`) VALUE (1000, "Nilorg Naas", "https://dss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=218375221,1552855610&fm=111&gp=0.jpg", "NilOrg认证授权服务平台", "http://localhost:9000");
-- oauth2_client_scope
INSERT INTO `oauth2_client_scope` (`id`, `created_at`, `updated_at`, `oauth2_client_id`, `scope_code`) VALUES ('1', '2020-06-29 14:35:49', '2020-06-29 14:35:49', '1000', 'openid');
INSERT INTO `oauth2_client_scope` (`id`, `created_at`, `updated_at`, `oauth2_client_id`, `scope_code`) VALUES ('2', '2020-06-29 14:35:49', '2020-06-29 14:35:49', '1000', 'profile');
INSERT INTO `oauth2_client_scope` (`id`, `created_at`, `updated_at`, `oauth2_client_id`, `scope_code`) VALUES ('3', '2020-06-29 14:35:49', '2020-06-29 14:35:49', '1000', 'email');
INSERT INTO `oauth2_client_scope` (`id`, `created_at`, `updated_at`, `oauth2_client_id`, `scope_code`) VALUES ('4', '2020-06-29 14:35:49', '2020-06-29 14:35:49', '1000', 'phone');

-- 创建casbin_rule索引
CREATE INDEX `idx_casbin_rule_v1` `casbin_rule` (v1) COMMENT '' ALGORITHM DEFAULT LOCK DEFAULT
