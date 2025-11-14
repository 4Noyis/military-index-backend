-- Insert technology categories
INSERT INTO tech_categories (name, description, icon_url) VALUES
('Aircraft', 'Military aircraft including fighters, bombers, and transport', '/icons/aircraft.svg'),
('Naval', 'Naval vessels and submarine technology', '/icons/naval.svg'),
('Ground Vehicles', 'Tanks, armored vehicles, and artillery', '/icons/ground.svg'),
('Missiles', 'Missile systems and rockets', '/icons/missiles.svg'),
('Electronics', 'Radar, communication, and electronic warfare', '/icons/electronics.svg'),
('Drones', 'Unmanned aerial vehicles and systems', '/icons/drones.svg'),
('Space & Satellite', 'Space-based military technology', '/icons/space.svg')
ON CONFLICT (name) DO NOTHING;

-- Insert countries
INSERT INTO countries (name, code, flag_url) VALUES
('Turkey', 'TUR', '/flags/turkey.svg'),
('United States', 'USA', '/flags/usa.svg'),
('Russia', 'RUS', '/flags/russia.svg'),
('China', 'CHN', '/flags/china.svg'),
('Germany', 'DEU', '/flags/germany.svg'),
('United Kingdom', 'GBR', '/flags/uk.svg'),
('France', 'FRA', '/flags/france.svg'),
('Israel', 'ISR', '/flags/israel.svg'),
('India', 'IND', '/flags/india.svg'),
('Japan', 'JPN', '/flags/japan.svg'),
('South Korea', 'KOR', '/flags/korea.svg'),
('Italy', 'ITA', '/flags/italy.svg')
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
