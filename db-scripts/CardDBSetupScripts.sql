CREATE TABLE mtg-cards.atomic_card_data (
	card_data_hash CHAR(16) PRIMARY KEY,
	color_identity SET('B', 'G', 'R', 'U', 'W') NULL,
	color_indicator SET('B', 'G', 'R', 'U', 'W') NULL,
	colors SET('B', 'G', 'R', 'U', 'W') NULL,
	converted_mana_cost FLOAT NOT NULL,
	edhrec_rank INTEGER NULL,
	face_converted_mana_cost FLOAT NOT NULL,
	hand VARCHAR(10) NULL, --Max existing len: 2
	is_reserved BOOLEAN NULL,
	layout VARCHAR(25) NOT NULL, --Max existing len: 9
	life VARCHAR(10) NULL, --Max existing len: 3
	loyalty VARCHAR(20) NULL, --Max existing len: 5
	mana_cost VARCHAR(100) NOT NULL, --Max existing len: 45
	mtgstocks_id INTEGER NOT NULL,
	name VARCHAR(500) NULL, --Max existing len: 141
	card_power VARCHAR(10) NOT NULL, --Max existing len: 3
	scryfall_oracle_id CHAR(36) NOT NULL,
	side CHAR(1) NULL,
	text VARCHAR(1500) NOT NULL, --Max existing len: 770
	toughness VARCHAR(10) NOT NULL, --Max existing len: 3
	card_type VARCHAR(100) NOT NULL --Max existing len: 46
)

CREATE TABLE mtg-cards.all_cards (
	uuid CHAR(36) PRIMARY KEY,
	atomic_card_data_hash CHAR(16) NOT NULL,
	artist VARCHAR(100) NOT NULL, --Max existing len: 54
	border_color VARCHAR(30) NOT NULL, --Max existing len: 10
	card_number VARCHAR(20) NOT NULL, --Max existing len: 9
	scryfall_id CHAR(36) NOT NULL,
	watermark VARCHAR(50) NOT NULL, --Max existing len: 21
	frame_version VARCHAR(20) NOT NULL, --Max existing len: 6
	mcm_id INT NOT NULL,
	mcm_meta_id INT NOT NULL,
	multiverse_id INT NOT NULL,
	original_text VARCHAR(1500) NOT NULL, --Max existing len: 760
	original_type VARCHAR(100) NOT NULL, --Max existing len: 52
	rarity VARCHAR(20) NOT NULL, --Max existing len: 8
	tcgplayer_product_id INT NOT NULL,
	duel_deck CHAR(1) NULL,
	flavor_text VARCHAR(403) NULL, --Max existing len: 1000
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
	is_starter BOOLEAN NOT NULL DEFAULT FALSE,
	is_story_spotlight BOOLEAN NOT NULL DEFAULT FALSE,
	is_textless BOOLEAN NOT NULL DEFAULT FALSE,
	is_timeshifted BOOLEAN NOT NULL DEFAULT FALSE,
	mtg_arena_id INT NULL,
	mtgo_foil_id INT NULL,
	mtgo_id INT NULL,
	scryfall_illustration_id CHAR(36) NULL
)

CREATE TABLE mtg-cards.frame_effects (
	frame_effect_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_uuid CHAR(36) NOT NULL,
	frame_effect VARCHAR(50) NOT NULL --Max existing len: 22
)

CREATE TABLE mtg-cards.variations (
	variation_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_uuid CHAR(36) NOT NULL,
	variation_uuid CHAR(36) NOT NULL
)

CREATE TABLE mtg-cards.other_faces (
	other_face_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_uuid CHAR(36) NOT NULL,
	other_face_uuid CHAR(36) NOT NULL
)

CREATE TABLE mtg-cards.alternate_language_data (
	alt_lang_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_data_hash CHAR(16) NOT NULL,
	flavor_text VARCHAR(1000) NOT NULL, --Max existing len: 574
	language VARCHAR(50) NOT NULL, --Max existing len: 19
	multiverse_id INTEGER NOT NULL,
	name VARCHAR(300) NOT NULL, --Max existing len: 93
	text VARCHAR(2000) NOT NULL, --Max existing len: 1006
	card_type VARCHAR(200) NOT NULL -- Max existing len: 96
)

CREATE TABLE mtg-cards.leadership_skills (
	card_data_hash CHAR(16) PRIMARY KEY,
	brawl_leader_legal BOOLEAN NOT NULL,
	commander_leader_legal BOOLEAN NOT NULL,
	oathbreaker_leader_legal BOOLEAN NOT NULL
)

CREATE TABLE mtg-cards.legalities (
	legality_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_data_hash CHAR(16),
	game_format_id INT NOT NULL,
	legality_option_id INT NOT NULL
)

CREATE TABLE mtg-cards.game_formats (
	game_format_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	game_format_name VARCHAR(30) NOT NULL
)

CREATE TABLE mtg-cards.legality_options (
	legality_option_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	legality_option_name VARCHAR(30) NOT NULL
)

CREATE TABLE mtg-cards.sets (
	set_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	set_hash CHAR(16) NOT NULL,
	base_size INT NOT NULL,
	block_name VARCHAR(50) NULL, --Max existing len: 22
	code VARCHAR(20) NOT NULL, --Max existing len: 6
	is_foreign_only BOOLEAN NOT NULL DEFAULT FALSE,
	is_foil_only BOOLEAN NOT NULL DEFAULT FALSE,
	is_online_only BOOLEAN NOT NULL DEFAULT FALSE,
	is_partial_preview BOOLEAN NOT NULL DEFAULT FALSE,
	keyrune_code VARCHAR(30) NOT NULL, --Max existing len: 12
	mcm_name VARCHAR(100) NOT NULL, --Max existing len: 39
	mcm_id INT NOT NULL,
	mtgo_code VARCHAR(20) NOT NULL, --Max existing len: 4
	name VARCHAR(100) NOT NULL, --Max existing len: 42
	parent_code VARCHAR(20) NOT NULL, --Max existing len: 3
	release_date DATE NOT NULL,
	tcgplayer_group_id INT NOT NULL,
	total_set_size INT NOT NULL,
	set_type VARCHAR(50) NOT NULL --Max existing len: 16
)

CREATE TABLE mtg-cards.card_printings (

	card_printing_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_data_hash CHAR(16) NOT NULL,
	set_id INT NOT NULL
)

CREATE TABLE mtg-cards.purchase_sites (
	purchase_site_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	purchase_site_name VARCHAR(30) NOT NULL
)

CREATE TABLE mtg-cards.purchase_urls (
	purchase_url_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	purchase_site_id INT NOT NULL,
	purchase_url VARCHAR(100) NOT NULL --Max existing len: 42
)

CREATE TABLE mtg-cards.rulings (
	ruling_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_data_hash CHAR(16) NOT NULL,
	ruling_date DATE NOT NULL,
	ruling_text VARCHAR(3000) NOT NULL --Max existing len: 1513
)

CREATE TABLE mtg-cards.card_subtypes (
	subtype_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_data_hash CHAR(16) NOT NULL,
	card_subtype VARCHAR(50) NOT NULL --Max existing len: 26
)

CREATE TABLE mtg-cards.card_supertypes (
	supertype_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	card_data_hash CHAR(16) NOT NULL,
	card_supertype VARCHAR(50) NOT NULL --Max existing len: 9
)

CREATE TABLE mtg-cards.set_translation_languages (
	set_translation_language_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	set_translation_language VARCHAR(50) NOT NULL
)

CREATE TABLE mtg-cards.set_translations (
	set_translation_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	set_id INT NOT NULL,
	set_translation_language_id INT NOT NULL,
	set_translated_name VARCHAR(200) NOT NULL --Max existing len: 67
)
