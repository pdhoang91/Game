-- Create Users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create PlayerResources table
CREATE TABLE IF NOT EXISTS player_resources (
    user_id VARCHAR(36) PRIMARY KEY,
    gold INT NOT NULL DEFAULT 0,
    premium_currency INT NOT NULL DEFAULT 0,
    last_idle_claim TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create HeroTypes table
CREATE TABLE IF NOT EXISTS hero_types (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    rarity VARCHAR(20) NOT NULL,
    base_hp INT NOT NULL,
    base_atk INT NOT NULL,
    description TEXT,
    image_url VARCHAR(255)
);

-- Create Skills table
CREATE TABLE IF NOT EXISTS skills (
    id VARCHAR(36) PRIMARY KEY,
    hero_type_id VARCHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    damage_multiplier FLOAT NOT NULL,
    cooldown INT NOT NULL,
    targets_all BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (hero_type_id) REFERENCES hero_types(id) ON DELETE CASCADE
);

-- Create Heroes table
CREATE TABLE IF NOT EXISTS heroes (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    hero_type_id VARCHAR(36) NOT NULL,
    level INT NOT NULL DEFAULT 1,
    experience INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (hero_type_id) REFERENCES hero_types(id)
);

-- Create Teams table
CREATE TABLE IF NOT EXISTS teams (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    position_1 VARCHAR(36),
    position_2 VARCHAR(36),
    position_3 VARCHAR(36),
    position_4 VARCHAR(36),
    position_5 VARCHAR(36),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (position_1) REFERENCES heroes(id) ON DELETE SET NULL,
    FOREIGN KEY (position_2) REFERENCES heroes(id) ON DELETE SET NULL,
    FOREIGN KEY (position_3) REFERENCES heroes(id) ON DELETE SET NULL,
    FOREIGN KEY (position_4) REFERENCES heroes(id) ON DELETE SET NULL,
    FOREIGN KEY (position_5) REFERENCES heroes(id) ON DELETE SET NULL
);

-- Create Stages table
CREATE TABLE IF NOT EXISTS stages (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    enemy_1 VARCHAR(36),
    enemy_2 VARCHAR(36),
    enemy_3 VARCHAR(36),
    enemy_4 VARCHAR(36),
    enemy_5 VARCHAR(36),
    gold_reward INT NOT NULL DEFAULT 0,
    exp_reward INT NOT NULL DEFAULT 0
);

-- Create BattleResults table
CREATE TABLE IF NOT EXISTS battle_results (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    team_id VARCHAR(36) NOT NULL,
    stage_id VARCHAR(36) NOT NULL,
    result VARCHAR(20) NOT NULL,
    rewards_json TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (stage_id) REFERENCES stages(id)
);

-- Create ItemTemplates table
CREATE TABLE IF NOT EXISTS item_templates (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL,
    rarity VARCHAR(20) NOT NULL,
    image_url VARCHAR(255),
    slot VARCHAR(20),
    atk_bonus INT DEFAULT 0,
    hp_bonus INT DEFAULT 0,
    effect VARCHAR(50),
    effect_value INT DEFAULT 0
);

-- Create Items table
CREATE TABLE IF NOT EXISTS items (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    item_template_id VARCHAR(36) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    equipped_to_hero_id VARCHAR(36),
    acquired_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (item_template_id) REFERENCES item_templates(id),
    FOREIGN KEY (equipped_to_hero_id) REFERENCES heroes(id) ON DELETE SET NULL
);

-- Create MissionTemplates table
CREATE TABLE IF NOT EXISTS mission_templates (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL,
    requirement_type VARCHAR(50) NOT NULL,
    target_value INT NOT NULL,
    target_id VARCHAR(36),
    gold_reward INT NOT NULL DEFAULT 0,
    gems_reward INT NOT NULL DEFAULT 0,
    experience_reward INT NOT NULL DEFAULT 0
);

-- Create MissionItemRewards table
CREATE TABLE IF NOT EXISTS mission_item_rewards (
    mission_template_id VARCHAR(36) NOT NULL,
    item_template_id VARCHAR(36) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    PRIMARY KEY (mission_template_id, item_template_id),
    FOREIGN KEY (mission_template_id) REFERENCES mission_templates(id) ON DELETE CASCADE,
    FOREIGN KEY (item_template_id) REFERENCES item_templates(id)
);

-- Create Missions table
CREATE TABLE IF NOT EXISTS missions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    mission_template_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL,
    current_value INT NOT NULL DEFAULT 0,
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    claimed_at TIMESTAMP,
    expires_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (mission_template_id) REFERENCES mission_templates(id)
);

-- Create Banners table
CREATE TABLE IF NOT EXISTS banners (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL,
    image_url VARCHAR(255),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    standard_hero_rate FLOAT NOT NULL,
    featured_hero_rate FLOAT NOT NULL,
    guarantee_threshold INT NOT NULL,
    single_summon_cost INT NOT NULL,
    ten_summon_cost INT NOT NULL,
    cost_type VARCHAR(20) NOT NULL,
    has_daily_free_summon BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create BannerFeaturedHeroes table
CREATE TABLE IF NOT EXISTS banner_featured_heroes (
    banner_id VARCHAR(36) NOT NULL,
    hero_type_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (banner_id, hero_type_id),
    FOREIGN KEY (banner_id) REFERENCES banners(id) ON DELETE CASCADE,
    FOREIGN KEY (hero_type_id) REFERENCES hero_types(id)
);

-- Create BannerFeaturedItems table
CREATE TABLE IF NOT EXISTS banner_featured_items (
    banner_id VARCHAR(36) NOT NULL,
    item_template_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (banner_id, item_template_id),
    FOREIGN KEY (banner_id) REFERENCES banners(id) ON DELETE CASCADE,
    FOREIGN KEY (item_template_id) REFERENCES item_templates(id)
);

-- Create SummonSessions table
CREATE TABLE IF NOT EXISTS summon_sessions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    banner_id VARCHAR(36) NOT NULL,
    pull_count INT NOT NULL DEFAULT 0,
    last_legendary_at INT NOT NULL DEFAULT 0,
    has_guarantee BOOLEAN NOT NULL DEFAULT FALSE,
    last_free_summon TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (banner_id) REFERENCES banners(id)
);

-- Create SummonResults table
CREATE TABLE IF NOT EXISTS summon_results (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    banner_id VARCHAR(36) NOT NULL,
    result_type VARCHAR(20) NOT NULL,
    result_id VARCHAR(36) NOT NULL,
    rarity VARCHAR(20) NOT NULL,
    is_featured BOOLEAN NOT NULL DEFAULT FALSE,
    is_pity_break BOOLEAN NOT NULL DEFAULT FALSE,
    pull_number INT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (banner_id) REFERENCES banners(id)
); 