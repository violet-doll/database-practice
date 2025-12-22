-- 1. Modify courses table, add capacity control fields
ALTER TABLE courses 
ADD COLUMN capacity INT DEFAULT 50 COMMENT 'Maximum course capacity',
ADD COLUMN enrolled_count INT DEFAULT 0 COMMENT 'Current enrolled count';

-- 2. Create course_prerequisites table (course prerequisite relationship - many-to-many self-reference)
CREATE TABLE IF NOT EXISTS course_prerequisites (
    course_id INT NOT NULL COMMENT 'Current course ID',
    prereq_id INT NOT NULL COMMENT 'Prerequisite course ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (course_id, prereq_id),
    FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    FOREIGN KEY (prereq_id) REFERENCES courses(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Add index for performance optimization
CREATE INDEX idx_prereq_course ON course_prerequisites(prereq_id);

-- 3. Create enrollment stored procedure
DELIMITER //

CREATE PROCEDURE sp_enroll_student(
    IN p_student_id INT,
    IN p_course_id INT,
    OUT p_status INT,       -- 0: success, 1: failure
    OUT p_message VARCHAR(255)
)
BEGIN
    DECLARE v_capacity INT;
    DECLARE v_enrolled INT;
    DECLARE v_already_enrolled INT;
    DECLARE v_prereq_count INT;
    DECLARE v_prereq_met INT;

    -- Start transaction
    START TRANSACTION;

    -- 1. Check if already enrolled
    SELECT COUNT(*) INTO v_already_enrolled 
    FROM enrollments 
    WHERE student_id = p_student_id AND course_id = p_course_id AND deleted_at IS NULL;

    IF v_already_enrolled > 0 THEN
        SET p_status = 1;
        SET p_message = 'Already enrolled in this course';
        ROLLBACK;
    ELSE
        -- 2. Check prerequisite requirements
        -- Count how many prerequisites this course has
        SELECT COUNT(*) INTO v_prereq_count 
        FROM course_prerequisites 
        WHERE course_id = p_course_id;

        -- Count how many prerequisites the student has completed and passed (score >= 60)
        SELECT COUNT(*) INTO v_prereq_met
        FROM course_prerequisites cp
        JOIN enrollments e ON cp.prereq_id = e.course_id
        JOIN grades g ON e.id = g.enrollment_id
        WHERE cp.course_id = p_course_id 
          AND e.student_id = p_student_id
          AND g.score >= 60; 

        IF v_prereq_met < v_prereq_count THEN
            SET p_status = 1;
            SET p_message = 'Prerequisites not completed';
            ROLLBACK;
        ELSE
            -- 3. Check and lock course capacity
            SELECT capacity, enrolled_count INTO v_capacity, v_enrolled
            FROM courses
            WHERE id = p_course_id
            FOR UPDATE;

            IF v_enrolled >= v_capacity THEN
                SET p_status = 1;
                SET p_message = 'Course is full';
                ROLLBACK;
            ELSE
                -- 4. Execute enrollment
                INSERT INTO enrollments (created_at, updated_at, student_id, course_id)
                VALUES (NOW(), NOW(), p_student_id, p_course_id);

                UPDATE courses 
                SET enrolled_count = enrolled_count + 1 
                WHERE id = p_course_id;

                SET p_status = 0;
                SET p_message = 'Enrollment successful';
                COMMIT;
            END IF;
        END IF;
    END IF;
END //

DELIMITER ;
