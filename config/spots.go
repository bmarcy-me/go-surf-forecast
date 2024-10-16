package config

type SpotConfig struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Lat       string `json:"lat"`
	Long      string `json:"long"`
	Direction int    `json:"direction"`
}

var SpotConfigs = []SpotConfig{
	{
		Id:        1,
		Name:      "Plage de Gros Joncs - Ile de Ré",
		Lat:       "46.1740867",
		Long:      "-1.3853837",
		Direction: 220,
	},
	{
		Id:        2,
		Name:      "Pointe du Lizay - Ile de Ré",
		Lat:       "46.257935",
		Long:      "-1.518474",
		Direction: 320,
	},
	{
		Id:        3,
		Name:      "Plage de Vert Bois - Ile d'Oléron",
		Lat:       "45.874214",
		Long:      "-1.263475",
		Direction: 260,
	},
}
