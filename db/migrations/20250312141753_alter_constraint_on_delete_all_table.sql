-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_roles DROP CONSTRAINT fk_user_roles_user_id;
ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE user_roles DROP CONSTRAINT fk_user_roles_role_id;
ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE application_permissions DROP CONSTRAINT fk_application_permissions_application_id;
ALTER TABLE application_permissions ADD CONSTRAINT fk_application_permissions_application_id FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE application_permissions DROP CONSTRAINT fk_application_permissions_permission_id;
ALTER TABLE application_permissions ADD CONSTRAINT fk_application_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE role_permissions DROP CONSTRAINT fk_role_permissions_role_id;
ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE role_permissions DROP CONSTRAINT fk_role_permissions_permission_id;
ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE user_profiles DROP CONSTRAINT fk_user_profiles_user_id;
ALTER TABLE user_profiles ADD CONSTRAINT fk_user_profiles_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE role_app_permissions DROP CONSTRAINT fk_role_app_permissions_role_id;
ALTER TABLE role_app_permissions ADD CONSTRAINT fk_role_app_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE role_app_permissions DROP CONSTRAINT fk_role_app_permissions_app_permission_id;
ALTER TABLE role_app_permissions ADD CONSTRAINT fk_role_app_permissions_app_permission_id FOREIGN KEY (app_permission_id) REFERENCES application_permissions(id) ON DELETE RESTRICT ON UPDATE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_roles DROP CONSTRAINT fk_user_roles_user_id;
ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE user_roles DROP CONSTRAINT fk_user_roles_role_id;
ALTER TABLE user_roles ADD CONSTRAINT fk_user_roles_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE application_permissions DROP CONSTRAINT fk_application_permissions_application_id;
ALTER TABLE application_permissions ADD CONSTRAINT fk_application_permissions_application_id FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE application_permissions DROP CONSTRAINT fk_application_permissions_permission_id;
ALTER TABLE application_permissions ADD CONSTRAINT fk_application_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE role_permissions DROP CONSTRAINT fk_role_permissions_role_id;
ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE role_permissions DROP CONSTRAINT fk_role_permissions_permission_id;
ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE user_profiles DROP CONSTRAINT fk_user_profiles_user_id;
ALTER TABLE user_profiles ADD CONSTRAINT fk_user_profiles_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE role_app_permissions DROP CONSTRAINT fk_role_app_permissions_role_id;
ALTER TABLE role_app_permissions ADD CONSTRAINT fk_role_app_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE role_app_permissions DROP CONSTRAINT fk_role_app_permissions_app_permission_id;
ALTER TABLE role_app_permissions ADD CONSTRAINT fk_role_app_permissions_app_permission_id FOREIGN KEY (app_permission_id) REFERENCES application_permissions(id) ON DELETE CASCADE ON UPDATE CASCADE;
-- +goose StatementEnd
