-- Insert technology categories
INSERT INTO tech_categories (name, description, icon_url) VALUES
('Aircraft', 'Military aircraft including fighters, bombers, and transport', 'https://api.iconify.design/mdi/airplane.svg'),
('Naval', 'Naval vessels and submarine technology', 'https://api.iconify.design/mdi/ferry.svg'),
('Ground Vehicles', 'Tanks, armored vehicles, and artillery', 'https://api.iconify.design/mdi/tank.svg'),
('Missiles', 'Missile systems and rockets', 'https://api.iconify.design/mdi/rocket-launch.svg'),
('Electronics', 'Radar, communication, and electronic warfare', 'https://api.iconify.design/mdi/radar.svg'),
('Drones', 'Unmanned aerial vehicles and systems', 'https://api.iconify.design/mdi/quadcopter.svg'),
('Space & Satellite', 'Space-based military technology', 'https://api.iconify.design/mdi/satellite-variant.svg')
ON CONFLICT (name) DO NOTHING;

-- Insert countries
INSERT INTO countries (name, code, flag_url) VALUES
('Turkey', 'TUR', 'https://flagcdn.com/w320/tr.png'),
('United States', 'USA', 'https://flagcdn.com/w320/us.png'),
('Russia', 'RUS', 'https://flagcdn.com/w320/ru.png'),
('China', 'CHN', 'https://flagcdn.com/w320/cn.png'),
('Germany', 'DEU', 'https://flagcdn.com/w320/de.png'),
('United Kingdom', 'GBR', 'https://flagcdn.com/w320/gb.png'),
('France', 'FRA', 'https://flagcdn.com/w320/fr.png'),
('Israel', 'ISR', 'https://flagcdn.com/w320/il.png'),
('India', 'IND', 'https://flagcdn.com/w320/in.png'),
('Japan', 'JPN', 'https://flagcdn.com/w320/jp.png'),
('South Korea', 'KOR', 'https://flagcdn.com/w320/kr.png'),
('Italy', 'ITA', 'https://flagcdn.com/w320/it.png')
ON CONFLICT (code) DO NOTHING;


-- Sample Turkish military technologies
INSERT INTO technologies (
    country_id, category_id, name, description, designer,
    year_developed, year_deployed, manufacturer, unit_cost,
    mass, length, width, height, status, image_url
) VALUES
(
    (SELECT id FROM countries WHERE code = 'TUR'),
    (SELECT id FROM tech_categories WHERE name = 'Aircraft'),
    'KAAN (TF-X)',
    'Fifth-generation stealth multirole fighter aircraft',
    'TAI',
    2023,
    2028,
    'Turkish Aerospace Industries',
    100000000,
    27000,
    2100,
    1400,
    550,
    'prototype',
    '/images/kaan.jpg'
),
(
    (SELECT id FROM countries WHERE code = 'TUR'),
    (SELECT id FROM tech_categories WHERE name = 'Drones'),
    'Bayraktar TB2',
    'Medium-altitude long-endurance tactical unmanned combat aerial vehicle',
    'Baykar',
    2014,
    2014,
    'Baykar Makina',
    5000000,
    700,
    640,
    1200,
    220,
    'current',
    '/images/tb2.jpg'
),
(
    (SELECT id FROM countries WHERE code = 'TUR'),
    (SELECT id FROM tech_categories WHERE name = 'Drones'),
    'Bayraktar Akıncı',
    'High-altitude long-endurance unmanned combat aerial vehicle',
    'Baykar',
    2019,
    2021,
    'Baykar Makina',
    NULL,
    1350,
    1250,
    2000,
    440,
    'current',
    '/images/akinci.jpg'
),
(
    (SELECT id FROM countries WHERE code = 'TUR'),
    (SELECT id FROM tech_categories WHERE name = 'Ground Vehicles'),
    'Altay MBT',
    'Main battle tank with advanced composite armor',
    'Otokar/FNSS',
    2013,
    2024,
    'BMC',
    NULL,
    65000,
    770,
    390,
    250,
    'current',
    '/images/altay.jpg'
),
(
    (SELECT id FROM countries WHERE code = 'TUR'),
    (SELECT id FROM tech_categories WHERE name = 'Naval'),
    'TCG Anadolu',
    'Landing helicopter dock (LHD) amphibious assault ship',
    'Sedef Shipbuilding',
    2015,
    2023,
    'Sedef Shipyard',
    NULL,
    27000000,
    23100,
    3200,
    5800,
    'current',
    '/images/anadolu.jpg'
),
(
    (SELECT id FROM countries WHERE code = 'TUR'),
    (SELECT id FROM tech_categories WHERE name = 'Aircraft'),
    'Hürjet',
    'Advanced jet trainer and light attack aircraft',
    'TAI',
    2022,
    2025,
    'Turkish Aerospace Industries',
    NULL,
    5500,
    1380,
    930,
    480,
    'prototype',
    '/images/hurjet.jpg'
);
