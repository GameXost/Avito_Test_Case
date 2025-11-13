CREATE TABLE teams
(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users
(
    user_id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    team_id BIGINT NOT NULL REFERENCES teams(id),
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE pull_request
(
    pull_request_id VARCHAR(255) PRIMARY KEY,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) NOT NULL REFERENCES users(user_id),
    status VARCHAR(16) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    merged_at TIMESTAMP
);

CREATE TABLE pr_reviewers
(
    pr_id VARCHAR(255) NOT NULL REFERENCES pull_request(pull_request_id),
    reviewer_id VARCHAR(255) NOT NULL REFERENCES users(user_id),
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (pr_id, reviewer_id)

);