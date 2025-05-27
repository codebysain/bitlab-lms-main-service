DO $$
    DECLARE
        v_realm_id TEXT;
        v_admin_id TEXT;
        v_teacher_id TEXT;
        v_student_id TEXT;
    BEGIN
        -- Get realm ID dynamically
        SELECT id INTO v_realm_id FROM realm WHERE name = 'lms';

        -- Create ROLE_ADMIN if not exists
        IF NOT EXISTS (
            SELECT 1 FROM keycloak_role
            WHERE name = 'ROLE_ADMIN' AND realm_id = v_realm_id
        ) THEN
            INSERT INTO keycloak_role (
                id, client_realm_constraint, client_role, description,
                name, realm_id, client
            ) VALUES (
                         gen_random_uuid(),
                         v_realm_id,
                         false,
                         'System Administrator',
                         'ROLE_ADMIN',
                         v_realm_id,
                         NULL
                     );
        END IF;

        -- Create ROLE_TEACHER if not exists
        IF NOT EXISTS (
            SELECT 1 FROM keycloak_role
            WHERE name = 'ROLE_TEACHER' AND realm_id = v_realm_id
        ) THEN
            INSERT INTO keycloak_role (
                id, client_realm_constraint, client_role, description,
                name, realm_id, client
            ) VALUES (
                         gen_random_uuid(),
                         v_realm_id,
                         false,
                         'Teacher role',
                         'ROLE_TEACHER',
                         v_realm_id,
                         NULL
                     );
        END IF;

        -- Create ROLE_STUDENT if not exists
        IF NOT EXISTS (
            SELECT 1 FROM keycloak_role
            WHERE name = 'ROLE_STUDENT' AND realm_id = v_realm_id
        ) THEN
            INSERT INTO keycloak_role (
                id, client_realm_constraint, client_role, description,
                name, realm_id, client
            ) VALUES (
                         gen_random_uuid(),
                         v_realm_id,
                         false,
                         'Student role',
                         'ROLE_STUDENT',
                         v_realm_id,
                         NULL
                     );
        END IF;

        -- Create admin user
        IF NOT EXISTS (
            SELECT 1 FROM user_entity WHERE username = 'admin' AND realm_id = v_realm_id
        ) THEN
            v_admin_id := gen_random_uuid();
            INSERT INTO user_entity (
                id, email, email_constraint, email_verified, enabled,
                federation_link, first_name, last_name, realm_id, username,
                created_timestamp, service_account_client_link, not_before
            ) VALUES (
                         v_admin_id,
                         'admin@lms.edu', 'admin@lms.edu', true, true,
                         NULL, 'System', 'Administrator', v_realm_id, 'admin',
                         EXTRACT(EPOCH FROM NOW()) * 1000, NULL, 0
                     );

            INSERT INTO user_role_mapping (role_id, user_id)
            SELECT id, v_admin_id FROM keycloak_role
            WHERE name = 'ROLE_ADMIN' AND realm_id = v_realm_id;
        END IF;

        -- Create teacher user
        IF NOT EXISTS (
            SELECT 1 FROM user_entity WHERE username = 'TeacherJohn' AND realm_id = v_realm_id
        ) THEN
            v_teacher_id := gen_random_uuid();
            INSERT INTO user_entity (
                id, email, email_constraint, email_verified, enabled,
                federation_link, first_name, last_name, realm_id, username,
                created_timestamp, service_account_client_link, not_before
            ) VALUES (
                         v_teacher_id,
                         'johnsmith@lms.edu', 'johnsmith@lms.edu', true, true,
                         NULL, 'John', 'Smith', v_realm_id, 'TeacherJohn',
                         EXTRACT(EPOCH FROM NOW()) * 1000, NULL, 0
                     );

            INSERT INTO user_role_mapping (role_id, user_id)
            SELECT id, v_teacher_id FROM keycloak_role
            WHERE name = 'ROLE_TEACHER' AND realm_id = v_realm_id;
        END IF;

        -- Create student user
        IF NOT EXISTS (
            SELECT 1 FROM user_entity WHERE username = 'StudentAlice' AND realm_id = v_realm_id
        ) THEN
            v_student_id := gen_random_uuid();
            INSERT INTO user_entity (
                id, email, email_constraint, email_verified, enabled,
                federation_link, first_name, last_name, realm_id, username,
                created_timestamp, service_account_client_link, not_before
            ) VALUES (
                         v_student_id,
                         'alicejohnson@lms.edu', 'alicejohnson@lms.edu', true, true,
                         NULL, 'Alice', 'Johnson', v_realm_id, 'StudentAlice',
                         EXTRACT(EPOCH FROM NOW()) * 1000, NULL, 0
                     );

            INSERT INTO user_role_mapping (role_id, user_id)
            SELECT id, v_student_id FROM keycloak_role
            WHERE name = 'ROLE_STUDENT' AND realm_id = v_realm_id;
        END IF;
    END $$;
