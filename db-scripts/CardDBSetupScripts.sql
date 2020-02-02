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

CREATE TABLE mtg_cards.all_tokens (
	token_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	uuid CHAR(36) NOT NULL,
	token_hash CHAR(32) NOT NULL,
	artist VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 54
	border_color VARCHAR(30) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 10
	card_number VARCHAR(20) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 9
	card_power VARCHAR(10) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 3
	card_type VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 46
	color_identity SET('B', 'G', 'R', 'U', 'W') NULL,
	color_indicator SET('B', 'G', 'R', 'U', 'W') NULL,
	colors SET('B', 'G', 'R', 'U', 'W') NULL,
	is_online_only BOOLEAN NOT NULL DEFAULT FALSE,
	layout VARCHAR(25) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 9
	loyalty VARCHAR(20) NULL COLLATE utf8mb4_general_ci, #Max existing len: 5
	name VARCHAR(500) NULL COLLATE utf8mb4_general_ci, #Max existing len: 141
	scryfall_id CHAR(36) NOT NULL,
	scryfall_illustration_id CHAR(36) NULL,
	scryfall_oracle_id CHAR(36) NOT NULL,
	set_id INT NOT NULL,
	side CHAR(1) NULL,
	text VARCHAR(1500) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 770
	toughness VARCHAR(10) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 3
	watermark VARCHAR(50) NOT NULL COLLATE utf8mb4_general_ci, #Max existing len: 21
	UNIQUE INDEX uuid_index (uuid),
	INDEX token_hash_index (token_hash),
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
("Artifact"), ("Autobot"), ("Card"), ("Character"), ("Conspiracy"), ("Creature"),
("Eaturecray"), ("Elemental"), ("Elite"), ("Emblem"), ("Enchantment"), ("Ever"),
("Hero"), ("Instant"), ("Land"), ("Phenomenon"), ("Plane"), ("Planeswalker"),
("Scariest"), ("Scheme"), ("See"), ("Sorcery"), ("Specter"), ("Summon"),
("Token"), ("Tribal"), ("Vanguard"), ("Wolf"), ("You’ll"), ("instant");

CREATE TABLE mtg_cards.card_base_types (
	base_type_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	base_type_option_id INT NOT NULL
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.token_base_types (
	base_type_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	token_id INT NOT NULL,
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
("Abian"), ("Adventure"), ("Advisor"), ("Aetherborn"), ("Ajani"), ("Alara"),
("Alicorn"), ("Alien"), ("Ally"), ("Aminatou"), ("Angel"), ("Angrath"),
("Antelope"), ("Ape"), ("Arcane"), ("Archer"), ("Archon"), ("Arkhos"),
("Arlinn"), ("Army"), ("Artificer"), ("Ashiok"), ("Assassin"), ("Assembly-Worker"),
("Atog"), ("Aura"), ("Aurochs"), ("Autobot"), ("Avatar"), ("Azgol"),
("Azra"), ("Baddest,"), ("Badger"), ("Barbarian"), ("Basilisk"), ("Bat"),
("Bear"), ("Beast"), ("Beaver"), ("Beeble"), ("Beholder"), ("Belenon"),
("Berserker"), ("Biggest,"), ("Bird"), ("Boar"), ("Bolas"), ("Bolas’s Meditation Realm"),
("Bot"), ("Brainiac"), ("Bringer"), ("Brushwagg"), ("Bureaucrat"), ("Calix"),
("Camel"), ("Carrier"), ("Cartouche"), ("Cat"), ("Centaur"), ("Cephalid"),
("Chandra"), ("Chicken"), ("Child"), ("Chimera"), ("Citizen"), ("Clamfolk"),
("Cleric"), ("Cloud"), ("Clue"), ("Cockatrice"), ("Construct"), ("Contraption"),
("Cow"), ("Crab"), ("Crocodile"), ("Curse"), ("Cyborg"), ("Cyclops"),
("Dack"), ("Daretti"), ("Dauthi"), ("Davriel"), ("Deer"), ("Demigod"),
("Demon"), ("Desert"), ("Designer"), ("Devil"), ("Dinosaur"), ("Djinn"),
("Dog"), ("Dominaria"), ("Domri"), ("Donkey"), ("Dovin"), ("Dragon"),
("Drake"), ("Dreadnought"), ("Drone"), ("Druid"), ("Dryad"), ("Duck"),
("Dungeon"), ("Dwarf"), ("Efreet"), ("Egg"), ("Elder"), ("Eldrazi"),
("Elemental"), ("Elemental?"), ("Elephant"), ("Elf"), ("Elk"), ("Elspeth"),
("Elves"), ("Equilor"), ("Equipment"), ("Ergamon"), ("Estrid"), ("Etiquette"),
("Eye"), ("Fabacin"), ("Faerie"), ("Ferret"), ("Fire"), ("Fish"),
("Flagbearer"), ("Food"), ("Forest"), ("Fortification"), ("Fox"), ("Freyalise"),
("Frog"), ("Fungus"), ("Gamer"), ("Gargoyle"), ("Garruk"), ("Gate"),
("Germ"), ("Giant"), ("Gideon"), ("Gnome"), ("Goat"), ("Goblin"),
("God"), ("Golem"), ("Gorgon"), ("Gremlin"), ("Griffin"), ("Gus"),
("Hag"), ("Harpy"), ("Hatificer"), ("Head"), ("Hellion"), ("Hero"),
("Hippo"), ("Hippogriff"), ("Homarid"), ("Homunculus"), ("Horror"), ("Horse"),
("Hound"), ("Huatli"), ("Human"), ("Hydra"), ("Hyena"), ("Igpay"),
("Illusion"), ("Imp"), ("Incarnation"), ("Innistrad"), ("Insect"), ("Inzerva"),
("Iquatana"), ("Ir"), ("Island"), ("Jace"), ("Jackal"), ("Jaya"),
("Jellyfish"), ("Juggernaut"), ("Kaldheim"), ("Kamigawa"), ("Kangaroo"), ("Karn"),
("Karsus"), ("Kasmina"), ("Kavu"), ("Kaya"), ("Kephalai"), ("Key"),
("Killbot"), ("Kinshala"), ("Kiora"), ("Kirin"), ("Kithkin"), ("Knight"),
("Kobold"), ("Kolbahan"), ("Kor"), ("Koth"), ("Kraken"), ("Kyneth"),
("Lady"), ("Lair"), ("Lamia"), ("Lammasu"), ("Leech"), ("Legend"),
("Leviathan"), ("Lhurgoyf"), ("Licid"), ("Liliana"), ("Lizard"), ("Lobster"),
("Locus"), ("Lorwyn"), ("Luvion"), ("Mammoth"), ("Manticore"), ("Master"),
("Masticore"), ("Mercadia"), ("Mercenary"), ("Merfolk"), ("Metathran"), ("Mime"),
("Mine"), ("Minion"), ("Minotaur"), ("Mirrodin"), ("Moag"), ("Mode"),
("Mole"), ("Monger"), ("Mongoose"), ("Mongseng"), ("Monk"), ("Monkey"),
("Moonfolk"), ("Mountain"), ("Mouse"), ("Mummy"), ("Muraganda"), ("Mutant"),
("Myr"), ("Mystic"), ("Naga"), ("Nahiri"), ("Narset"), ("Nastiest,"),
("Nautilus"), ("Nephilim"), ("New Phyrexia"), ("Nightmare"), ("Nightstalker"), ("Ninja"),
("Nissa"), ("Nixilis"), ("Noble"), ("Noggle"), ("Nomad"), ("Nymph"),
("Octopus"), ("Ogre"), ("Oko"), ("Ooze"), ("Orc"), ("Orgg"),
("Ouphe"), ("Ox"), ("Oyster"), ("Pangolin"), ("Paratrooper"), ("Peasant"),
("Pegasus"), ("Penguin"), ("Pentavite"), ("Pest"), ("Phelddagrif"), ("Phoenix"),
("Phyrexia"), ("Phyrexian"), ("Pilot"), ("Pirate"), ("Plains"), ("Plant"),
("Power-Plant"), ("Praetor"), ("Processor"), ("Proper"), ("Pyrulea"), ("Rabbit"),
("Rabiah"), ("Raccoon"), ("Ral"), ("Rat"), ("Rath"), ("Ravnica"),
("Rebel"), ("Reflection"), ("Regatha"), ("Reveler"), ("Rhino"), ("Rigger"),
("Rogue"), ("Rowan"), ("Rukh"), ("Sable"), ("Saga"), ("Saheeli"),
("Salamander"), ("Samurai"), ("Samut"), ("Saproling"), ("Sarkhan"), ("Satyr"),
("Scarecrow"), ("Scientist"), ("Scion"), ("Scorpion"), ("Scout"), ("Sculpture"),
("Segovia"), ("Serf"), ("Serpent"), ("Serra"), ("Serra’s Realm"), ("Servo"),
("Shade"), ("Shadowmoor"), ("Shaman"), ("Shandalar"), ("Shapeshifter"), ("Sheep"),
("Ship"), ("Shrine"), ("Siren"), ("Skeleton"), ("Slith"), ("Sliver"),
("Slug"), ("Snake"), ("Soldier"), ("Soltari"), ("Sorin"), ("Spawn"),
("Specter"), ("Spellshaper"), ("Sphinx"), ("Spider"), ("Spike"), ("Spirit"),
("Sponge"), ("Spy"), ("Squid"), ("Squirrel"), ("Starfish"), ("Surrakar"),
("Survivor"), ("Swamp"), ("Tamiyo"), ("Teferi"), ("Tentacle"), ("Teyo"),
("Tezzeret"), ("Thalakos"), ("The"), ("Thopter"), ("Thrull"), ("Tibalt"),
("Tower"), ("Townsfolk"), ("Trap"), ("Treasure"), ("Treefolk"), ("Trilobite"),
("Triskelavite"), ("Troll"), ("Turtle"), ("Ugin"), ("Ulgrotha"), ("Unicorn"),
("Urza"), ("Urza’s"), ("Valla"), ("Vampire"), ("Vampyre"), ("Vedalken"),
("Vehicle"), ("Venser"), ("Viashino"), ("Villain"), ("Vivien"), ("Volver"),
("Vraska"), ("Vryn"), ("Waiter"), ("Wall"), ("Warlock"), ("Warrior"),
("Weird"), ("Werewolf"), ("Whale"), ("Wildfire"), ("Will"), ("Windgrace"),
("Wizard"), ("Wolf"), ("Wolverine"), ("Wombat"), ("Worm"), ("Wraith"),
("Wrenn"), ("Wrestler"), ("Wurm"), ("Xenagos"), ("Xerex"), ("Yanggu"),
("Yanling"), ("Yeti"), ("Zendikar"), ("Zombie"), ("Zubera"), ("and/or"),
("of");

CREATE TABLE mtg_cards.card_subtypes (
	subtype_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	subtype_option_id INT NOT NULL
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.token_subtypes (
	subtype_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	token_id INT NOT NULL,
	subtype_option_id INT NOT NULL
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.card_supertype_options (
	supertype_option_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	supertype_option VARCHAR(50) NOT NULL COLLATE utf8mb4_general_ci #Max existing len: 9
) DEFAULT COLLATE utf8mb4_bin;

INSERT INTO mtg_cards.card_supertype_options
(supertype_option)
VALUES
("Basic"), ("Host"), ("Legendary"), ("Ongoing"), ("Snow"), ("World");

CREATE TABLE mtg_cards.card_supertypes (
	supertype_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_id INT NOT NULL,
	supertype_option_id INT NOT NULL
) DEFAULT COLLATE utf8mb4_bin;

CREATE TABLE mtg_cards.token_supertypes (
	supertype_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	token_id INT NOT NULL,
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

CREATE TABLE mtg_cards.token_reverse_related (
	reverse_related_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	token_id INT NOT NULL,
	reverse_related_card VARCHAR(200) NOT NULL COLLATE utf8mb4_general_ci #Max existing len: 57
) DEFAULT COLLATE utf8mb4_bin;
