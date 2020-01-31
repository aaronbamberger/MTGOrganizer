USE mtg_cards;

CREATE TABLE mtg_cards.all_cards (
	card_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	uuid CHAR(36) NOT NULL,
	card_hash CHAR(32) NOT NULL,
	artist VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 54
	border_color VARCHAR(30) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 10
	card_number VARCHAR(20) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 9
	card_power VARCHAR(10) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 3
	card_type VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 46
	color_identity SET('B', 'G', 'R', 'U', 'W') NULL,
	color_indicator SET('B', 'G', 'R', 'U', 'W') NULL,
	colors SET('B', 'G', 'R', 'U', 'W') NULL,
	converted_mana_cost FLOAT NOT NULL,
	duel_deck CHAR(1) NULL COLLATE utf8mb4_general_ci,
	edhrec_rank INTEGER NULL,
	face_converted_mana_cost FLOAT NOT NULL,
	flavor_text VARCHAR(2000) NULL COLLATE utf8mb4_general_ci, #Max existing len: 1000
	frame_version VARCHAR(20) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 6
	hand VARCHAR(10) NULL COLLATE utf8mb4_general_ci, #Max existing len: 2
	has_foil BOOLEAN NOT NULL DEFAULT FALSE,
	has_non_foil BOOLEAN NOT NULL DEFAULT FALSE,
	is_alternative BOOLEAN NOT NULL DEFAULT FALSE,
	is_arena BOOLEAN NOT NULL DEFAULT FALSE,
	is_full_art BOOLEAN NOT NULL DEFAULT FALSE,
	is_mtgo BOOLEAN NOT NULL DEFAULT FALSE,
	is_online_only BOOLEAN NOT NULL DEFAULT FALSE,
	is_oversized BOOLEAN NOT NULL DEFAULT FALSE,
	is_paper BOOLEAN NOT NULL DEFAULT FALSE,
	is_promo BOOLEAN NOT NULL DEFAULT FALSE,
	is_reprint BOOLEAN NOT NULL DEFAULT FALSE,
	is_reserved BOOLEAN NOT NULL,
	is_starter BOOLEAN NOT NULL DEFAULT FALSE,
	is_story_spotlight BOOLEAN NOT NULL DEFAULT FALSE,
	is_textless BOOLEAN NOT NULL DEFAULT FALSE,
	is_timeshifted BOOLEAN NOT NULL DEFAULT FALSE,
	layout VARCHAR(25) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 9
	life VARCHAR(10) NULL COLLATE utf8mb4_general_ci, #Max existing len: 3
	loyalty VARCHAR(20) NULL COLLATE utf8mb4_general_ci, #Max existing len: 5
	mana_cost VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 45
	mcm_id INT NOT NULL,
	mcm_meta_id INT NOT NULL,
	mtg_arena_id INT NULL,
	mtgo_foil_id INT NULL,
	mtgo_id INT NULL,
	mtgstocks_id INTEGER NOT NULL,
	multiverse_id INT NOT NULL,
	name VARCHAR(500) NULL COLLATE utf8mb4_general_ci, #Max existing len: 141
	original_text VARCHAR(1500) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 760
	original_type VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 52
	rarity VARCHAR(20) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 8
	scryfall_id CHAR(36) NOT NULL,
	scryfall_illustration_id CHAR(36) NULL,
	scryfall_oracle_id CHAR(36) NOT NULL,
	set_id INT NOT NULL,
	side CHAR(1) NULL,
	tcgplayer_product_id INT NOT NULL,
	text VARCHAR(1500) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 770
	toughness VARCHAR(10) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 3
	watermark VARCHAR(50) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 21
	UNIQUE INDEX uuid_index (uuid),
	INDEX card_hash_index (card_hash),
	INDEX name_index (name)
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.frame_effect_options (
	frame_effect_option_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	frame_effect_option VARCHAR(50) NOT NULL COLLATE utf8mb4_general_ci #Max existing len: 22
) DEFAULT COLLATE utf8mb4_bin;

INSERT INTO mtg_cards.frame_effect_options
(frame_effect_option)
VALUES
("legendary"), ("nyxtouched"), ("sunmoondfc"), ("extendedart"), ("devoid"),
("tombstone"), ("compasslanddfc"), ("showcase"), ("colorshifted"), ("originpwdfc"),
("mooneldrazidfc"), ("inverted"), ("draft"), ("miracle"), ("nyxborn"),
("waxingandwaningmoondfc");

CREATE TABLE mtg_cards.frame_effects (
	frame_effect_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	frame_effect_option_id INT NOT NULL
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.variations (
	variation_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	variation_uuid CHAR(36) NOT NULL
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.other_faces (
	other_face_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	other_face_uuid CHAR(36) NOT NULL
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.alternate_language_data (
	alt_lang_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	flavor_text VARCHAR(1000) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 574
	language VARCHAR(50) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 19
	multiverse_id INTEGER NOT NULL,
	name VARCHAR(300) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 93
	text VARCHAR(2000) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 1006
	card_type VARCHAR(200) NOT NULL COLLATE utf8mb4_general_ci #Max existing len: 96
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.leadership_skills (
	leadership_skill_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	leadership_format_id INT NOT NULL,
	leader_legal BOOLEAN NOT NULL
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.leadership_formats (
	leadership_format_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	leadership_format_name VARCHAR(30) NOT NULL COLLATE utf8mb4_general_ci
) DEFAULT COLLATE utf8mb4_bin;

INSERT INTO mtg_cards.leadership_formats
(leadership_format_name)
VALUES
("brawl"), ("commander"), ("oathbreaker");

CREATE TABLE mtg_cards.legalities (
	legality_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	game_format_id INT NOT NULL,
	legality_option_id INT NOT NULL
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.game_formats (
	game_format_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	game_format_name VARCHAR(30) NOT NULL COLLATE utf8mb4_general_ci
) DEFAULT COLLATE utf8mb4_bin;

INSERT INTO mtg_cards.game_formats
(game_format_name)
VALUES
("brawl"), ("commander"), ("duel"), ("future"), ("frontier"), ("legacy"), ("modern"),
("pauper"), ("penny"), ("pioneer"), ("standard"), ("vintage"), ("historic"),
("oldschool");

CREATE TABLE mtg_cards.legality_options (
	legality_option_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	legality_option_name VARCHAR(30) NOT NULL COLLATE utf8mb4_general_ci
) DEFAULT COLLATE utf8mb4_bin;

INSERT INTO mtg_cards.legality_options
(legality_option_name)
VALUES
("Legal"), ("Not Legal"), ("Restricted"), ("Banned");

CREATE TABLE mtg_cards.base_type_options (
	base_type_option_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	base_type_option VARCHAR(50) NOT NULL COLLATE utf8mb4_general_ci
) DEFAULT COLLATE utf8mb4_bin;

INSERT INTO mtg_cards.base_type_options
(base_type_option)
VALUES
("Creature"), ("Sorcery"), ("Instant"), ("Land"), ("Planeswalker"), ("Artifact"),
("Enchantment"), ("Tribal"), ("Scheme"), ("Hero"), ("Eaturecray"), ("Summon"),
("Plane"), ("Phenomenon"), ("Autobot"), ("Character"), ("Vanguard"), ("Conspiracy"),
("Scariest"), ("You'll"), ("Ever"), ("See"), ("instant"), ("Wolf"), ("Elemental"),
("Specter");

CREATE TABLE mtg_cards.base_types (
	base_type_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	base_type_option_id INT NOT NULL
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.sets (
	set_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	set_hash CHAR(32) NOT NULL,
	base_size INT NOT NULL,
	block_name VARCHAR(50) NULL COLLATE utf8mb4_general_ci, #Max existing len: 22
	code VARCHAR(20) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 6
	is_foreign_only BOOLEAN NOT NULL DEFAULT FALSE,
	is_foil_only BOOLEAN NOT NULL DEFAULT FALSE,
	is_online_only BOOLEAN NOT NULL DEFAULT FALSE,
	is_partial_preview BOOLEAN NOT NULL DEFAULT FALSE,
	keyrune_code VARCHAR(30) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 12
	mcm_name VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 39
	mcm_id INT NOT NULL,
	mtgo_code VARCHAR(20) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 4
	name VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 42
	parent_code VARCHAR(20) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 3
	release_date DATE NOT NULL,
	tcgplayer_group_id INT NOT NULL,
	total_set_size INT NOT NULL,
	set_type VARCHAR(50) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 16
	INDEX set_hash_index (set_hash),
	UNIQUE INDEX code_index (code)
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.card_printings (
	card_printing_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	set_code VARCHAR(20) NOT NULL COLLATE utf8mb4_general_ci #Max existing len: 6
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.purchase_sites (
	purchase_site_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	purchase_site_name VARCHAR(30) NOT NULL COLLATE utf8mb4_general_ci #Max existing len: 10
) DEFAULT COLLATE utf8mb4_bin;

INSERT INTO mtg_cards.purchase_sites
(purchase_site_name)
VALUES
("cardmarket"), ("tcgplayer"), ("mtgstocks");

CREATE TABLE mtg_cards.purchase_urls (
	purchase_url_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	purchase_site_id INT NOT NULL,
	purchase_url VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci #Max existing len: 42
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.rulings (
	ruling_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	ruling_date DATE NOT NULL,
	ruling_text VARCHAR(3000) NOT NULL COLLATE utf8mb4_general_ci #Max existing len: 1513
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.card_subtype_options (
	subtype_option_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	subtype_option VARCHAR(50) NOT NULL COLLATE utf8mb4_general_ci #Max existing len: 26
) DEFAULT COLLATE utf8mb4_bin;

INSERT INTO mtg_cards.card_subtype_options
(subtype_option)
VALUES
("Mammoth"), ("Beast"), ("Domri"), ("Wraith"), ("Gorgon"), ("Kinshala"), ("Licid"),
("Kangaroo"), ("Devil"), ("Monger"), ("Arkhos"), ("Brushwagg"), ("Igpay"), ("Hero"),
("Head"), ("Human"), ("Rhino"), ("Warrior"), ("Tezzeret"), ("Samurai"), ("Angrath"),
("Spawn"), ("Xenagos"), ("Construct"), ("Soldier"), ("Centaur"), ("Horse"), ("Food"),
("Kithkin"), ("Dack"), ("Wrestler"), ("Paratrooper"), ("Fire"), ("Juggernaut"),
("Antelope"), ("Pirate"), ("Ral"), ("Nastiest,"), ("Brainiac"), ("Davriel"),
("Illusion"), ("Cockatrice"), ("Curse"), ("Hound"), ("Ape"), ("Vraska"),
("Valla"), ("Vampire"), ("Ooze"), ("Demon"), ("Pegasus"), ("Clamfolk"), ("Forest"),
("Shade"), ("Slith"), ("Rigger"), ("Ir"), ("Koth"), ("Specter"), ("Oko"),
("Tamiyo"), ("Chicken"), ("Sphinx"), ("Phelddagrif"), ("Orgg"), ("Alicorn"),
("Reflection"), ("Spider"), ("Archon"), ("Djinn"), ("God"), ("Drone"), ("Tibalt"),
("Vampyre"), ("Freyalise"), ("Equipment"), ("Boar"), ("Monk"), ("Basilisk"),
("Aetherborn"), ("Azra"), ("Regatha"), ("Processor"), ("Dungeon"), ("Jellyfish"),
("Adventure"), ("Mercenary"), ("Daretti"), ("Skeleton"), ("Garruk"), ("Bolas"),
("Sponge"), ("Donkey"), ("Mime"), ("Elemental"), ("Ogre"), ("Viashino"),
("Nightstalker"), ("Bear"), ("Lamia"), ("Angel"), ("Whale"), ("Ugin"), ("Hatificer"),
("Homunculus"), ("Scout"), ("Dog"), ("Drake"), ("Orc"), ("Noggle"), ("Moag"),
("Power-Plant"), ("Gamer"), ("Citizen"), ("Tower"), ("Crocodile"), ("Giant"),
("Trap"), ("Luvion"), ("Muraganda"), ("Nephilim"), ("Cephalid"), ("The"),
("Designer"), ("Duck"), ("Equilor"), ("Deer"), ("Avatar"), ("Manticore"),
("Townsfolk"), ("Horror"), ("Treefolk"), ("Windgrace"), ("Rath"), ("of"), ("Swamp"),
("Berserker"), ("Wolf"), ("Yeti"), ("Squid"), ("Goat"), ("Teyo"), ("Iquatana"),
("Lizard"), ("Will"), ("Myr"), ("Lhurgoyf"), ("Nomad"), ("Cartouche"), ("Zendikar"),
("Ulgrotha"), ("Wombat"), ("Spy"), ("Monkey"), ("Starfish"), ("Kephalai"),
("Shrine"), ("Cloud"), ("Bat"), ("Mutant"), ("Elspeth"), ("Kamigawa"), ("Karsus"),
("Shandalar"),("Druid"), ("Dinosaur"), ("Crab"), ("Scarecrow"), ("Elk"), ("Chimera"),
("Mongoose"), ("Sable"), ("Lobster"), ("Surrakar"), ("Desert"), ("Nautilus"),
("Dominaria"), ("Gideon"), ("Mirrodin"), ("Dreadnought"), ("Homarid"), ("Kavu"),
("Pangolin"), ("Saga"), ("Rebel"), ("Oyster"), ("Lammasu"), ("Gate"), ("Vehicle"),
("Kor"), ("Estrid"), ("Trilobite"), ("Ravnica"), ("Wildfire"), ("Imp"),
("Flagbearer"), ("Beeble"), ("Elves"), ("Lady"), ("Rogue"), ("Rowan"), ("Lair"),
("Child"), ("Fish"), ("Ajani"), ("Chandra"), ("Kirin"), ("Bringer"), ("Zombie"),
("Spirit"), ("Faerie"), ("Teferi"), ("Urza"), ("Elder"), ("Wrenn"), ("Jaya"),
("Cyborg"), ("Elemental?"), ("Abian"), ("Dragon"), ("Vedalken"), ("Wurm"),
("Unicorn"), ("Ally"), ("Penguin"), ("Liliana"), ("Ouphe"), ("Aurochs"), ("Fabacin"),
("Mode"), ("Goblin"), ("Soltari"), ("Hag"), ("Octopus"), ("Worm"), ("Saheeli"),
("Locus"), ("Master"), ("Rat"), ("Artificer"), ("Assembly-Worker"), ("Praetor"),
("Samut"), ("Biggest,"), ("Kasmina"), ("Spike"), ("Shaman"), ("Turtle"),
("Incarnation"), ("Eldrazi"), ("Gremlin"), ("Hyena"), ("Hippo"),
("Bolas’s Meditation Realm"), ("Rabbit"), ("Kobold"), ("and/or"), ("Advisor"),
("Elephant"), ("Kiora"), ("Hellion"), ("Carrier"), ("Rabiah"), ("Vryn"), ("Thrull"),
("Wolverine"), ("Camel"), ("Aminatou"), ("Yanling"), ("Wizard"), ("Insect"),
("Spellshaper"), ("Atog"), ("Segovia"), ("Demigod"), ("Beaver"), ("Knight"),
("Elf"), ("Arlinn"), ("Arcane"), ("Pyrulea"), ("Mongseng"), ("Shadowmoor"),
("Bureaucrat"), ("Werewolf"), ("Harpy"), ("Sheep"), ("Kaya"), ("Kolbahan"),
("Phyrexia"), ("Autobot"), ("Gargoyle"), ("Phoenix"), ("Noble"), ("Siren"),
("Naga"), ("Alien"), ("Griffin"), ("Nahiri"), ("Narset"), ("Calix"), ("Waiter"),
("Island"), ("Cleric"), ("Jackal"), ("Serra’s Realm"), ("Leech"), ("Urza’s"),
("Legend"), ("Frog"), ("Hydra"), ("Assassin"), ("Moonfolk"), ("Venser"), ("Killbot"),
("Minotaur"), ("Azgol"), ("Ferret"), ("Sarkhan"), ("Shapeshifter"), ("Sorin"),
("Barbarian"), ("Dovin"),("Lorwyn"), ("Ergamon"), ("Archer"), ("Nymph"),
("Kaldheim"), ("Villain"), ("Zubera"), ("Fox"), ("Karn"), ("Xerex"), ("Etiquette"),
("Yanggu"), ("Squirrel"), ("Mystic"), ("Scientist"), ("Gnome"), ("Serpent"),
("Mole"), ("Nixilis"), ("Alara"), ("New Phyrexia"), ("Raccoon"), ("Key"),
("Belenon"), ("Merfolk"), ("Ox"), ("Troll"), ("Peasant"), ("Plant"), ("Efreet"),
("Serra"), ("Thalakos"), ("Jace"), ("Satyr"), ("Minion"), ("Kyneth"), ("Gus"),
("Contraption"), ("Beholder"), ("Innistrad"), ("Golem"), ("Nightmare"),
("Warlock"), ("Kraken"), ("Ninja"), ("Fortification"), ("Hippogriff"), ("Bot"),
("Baddest,"), ("Plains"), ("Wall"), ("Fungus"), ("Eye"), ("Volver"), ("Cow"),
("Ashiok"), ("Pest"), ("Leviathan"), ("Cat"), ("Sliver"), ("Salamander"), ("Weird"),
("Dauthi"), ("Vivien"), ("Mummy"), ("Egg"), ("Dryad"), ("Slug"), ("Scorpion"),
("Thopter"), ("Metathran"), ("Ship"), ("Masticore"), ("Inzerva"), ("Mercadia"),
("Phyrexian"), ("Bird"), ("Mountain"), ("Snake"), ("Dwarf"), ("Badger"), ("Pilot"),
("Aura"), ("Cyclops"), ("Nissa"), ("Huatli"), ("Mine"), ("Proper");

CREATE TABLE mtg_cards.card_subtypes (
	subtype_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	subtype_option_id INT NOT NULL
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.card_supertype_options (
	supertype_option_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	supertype_option VARCHAR(50) NOT NULL COLLATE utf8mb4_general_ci #Max existing len: 9
) DEFAULT COLLATE utf8mb4_bin;

INSERT INTO mtg_cards.card_supertype_options
(supertype_option)
VALUES
("Legendary"), ("Basic"), ("Snow"), ("Ongoing"), ("World");

CREATE TABLE mtg_cards.card_supertypes (
	supertype_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	supertype_option_id INT NOT NULL
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.set_translation_languages (
	set_translation_language_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	set_translation_language VARCHAR(50) NOT NULL COLLATE utf8mb4_general_ci #Max existing len: 19
) DEFAULT COLLATE utf8mb4_bin;

INSERT INTO mtg_cards.set_translation_languages
(set_translation_language)
VALUES
("Chinese Simplified"), ("Chinese Traditional"), ("French"), ("German"), ("Italian"),
("Japanese"), ("Korean"), ("Portuguese (Brazil)"), ("Russian"), ("Spanish");

CREATE TABLE mtg_cards.set_translations (
	set_translation_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	set_id INT NOT NULL,
	set_translation_language_id INT NOT NULL,
	set_translated_name VARCHAR(200) NOT NULL COLLATE utf8mb4_general_ci #Max existing len: 67
) DEFAULT COLLATE utf8mb4_bin;
