-- Insert sample hero types
INSERT INTO hero_types (id, name, rarity, base_hp, base_atk, description, image_url)
VALUES 
('hero_type_001', 'Warrior', 'rare', 500, 50, 'A strong warrior with balanced stats', 'heroes/warrior.png'),
('hero_type_002', 'Mage', 'rare', 300, 70, 'A powerful mage with high attack', 'heroes/mage.png'),
('hero_type_003', 'Archer', 'rare', 350, 60, 'A skilled archer with good range', 'heroes/archer.png'),
('hero_type_004', 'Knight', 'epic', 600, 55, 'A heavily armored knight with high HP', 'heroes/knight.png'),
('hero_type_005', 'Assassin', 'epic', 400, 75, 'A stealthy assassin with high attack', 'heroes/assassin.png');

-- Insert sample skills
INSERT INTO skills (id, hero_type_id, name, description, damage_multiplier, cooldown, targets_all)
VALUES 
('skill_001', 'hero_type_001', 'Mighty Slash', 'Deal 150% ATK to a single enemy', 1.5, 3, false),
('skill_002', 'hero_type_002', 'Fireball', 'Deal 120% ATK to all enemies', 1.2, 4, true),
('skill_003', 'hero_type_003', 'Quick Shot', 'Deal 130% ATK to a single enemy', 1.3, 2, false),
('skill_004', 'hero_type_004', 'Shield Bash', 'Deal 140% ATK to a single enemy', 1.4, 3, false),
('skill_005', 'hero_type_005', 'Backstab', 'Deal 160% ATK to a single enemy', 1.6, 3, false);

-- Insert sample stages
INSERT INTO stages (id, name, description, enemy_1, enemy_2, enemy_3, gold_reward, exp_reward)
VALUES 
('stage_001', 'Forest Path', 'A peaceful forest path with weak enemies', 'enemy_001', 'enemy_002', null, 100, 50),
('stage_002', 'Dark Cave', 'A dangerous cave with stronger enemies', 'enemy_003', 'enemy_004', 'enemy_005', 150, 75),
('stage_003', 'Mountain Pass', 'A treacherous mountain pass with powerful enemies', 'enemy_006', 'enemy_007', 'enemy_008', 200, 100);

-- Insert sample item templates
INSERT INTO item_templates (id, name, description, type, rarity, image_url, slot, atk_bonus, hp_bonus)
VALUES 
('item_template_001', 'Iron Sword', 'A basic iron sword', 'equipment', 'common', 'items/iron_sword.png', 'weapon', 10, 0),
('item_template_002', 'Steel Armor', 'Sturdy steel armor', 'equipment', 'common', 'items/steel_armor.png', 'armor', 0, 20),
('item_template_003', 'Silver Ring', 'A magical silver ring', 'equipment', 'uncommon', 'items/silver_ring.png', 'accessory', 5, 5),
('item_template_004', 'Health Potion', 'Restores 100 HP', 'consumable', 'common', 'items/health_potion.png', null, 0, 0),
('item_template_005', 'Iron Ore', 'Used for crafting weapons', 'material', 'common', 'items/iron_ore.png', null, 0, 0);

-- Insert sample mission templates
INSERT INTO mission_templates (id, title, description, type, requirement_type, target_value, gold_reward, gems_reward, experience_reward)
VALUES 
('mission_template_001', 'Win 3 Battles', 'Win 3 battles in any stage', 'daily', 'win_battles', 3, 100, 10, 50),
('mission_template_002', 'Level Up a Hero', 'Level up any hero', 'daily', 'level_up_hero', 1, 150, 15, 75),
('mission_template_003', 'Collect 5 Items', 'Collect any 5 items', 'weekly', 'collect_items', 5, 300, 30, 100);

-- Insert sample banners
INSERT INTO banners (id, name, description, type, image_url, start_time, end_time, standard_hero_rate, featured_hero_rate, guarantee_threshold, single_summon_cost, ten_summon_cost, cost_type, has_daily_free_summon)
VALUES 
('banner_001', 'Standard Summon', 'The standard summon banner with all heroes', 'standard', 'banners/standard.png', '2023-01-01 00:00:00', null, 0.03, 0.01, 100, 300, 2700, 'gem', true),
('banner_002', 'Knight & Assassin', 'Featured banner with Knight and Assassin', 'event', 'banners/knight_assassin.png', '2023-01-01 00:00:00', '2023-12-31 23:59:59', 0.02, 0.02, 80, 300, 2700, 'gem', false);

-- Insert banner featured heroes
INSERT INTO banner_featured_heroes (banner_id, hero_type_id)
VALUES 
('banner_002', 'hero_type_004'),
('banner_002', 'hero_type_005'); 