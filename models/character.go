package models

type Character struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Element     string  `json:"element"`
	WeaponType  string  `json:"weapon_type"`
	Rarity      string  `json:"rarity"`
	Role        *string `json:"role"`
	Description *string `json:"description"`
	ReleaseDate string  `json:"release_date"`
	BaseAttack  int     `json:"base_attack"`
	BaseDefense int     `json:"base_defense"`
	BaseHealth  int     `json:"base_health"`
}
