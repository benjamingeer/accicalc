package dataset

import "encoding/json"

type Jsonable interface {
	AsJson() (string, error)
}

type VoieSpéciale int

const (
	VoieSpécialeNonRenseignée = iota
	VoieSpécialeSansObjet
	PisteCyclable
	BandeCyclable
	AutreVoieSpéciale
)

func (voieSpéciale VoieSpéciale) String() string {
	return [...]string{
		"Non renseignée",
		"Sans objet",
		"Piste cyclable",
		"Bande cyclable",
		"Voie réservée",
	}[voieSpéciale]
}

func (voieSpéciale VoieSpéciale) MarshalJSON() ([]byte, error) {
	return json.Marshal(voieSpéciale.String())
}

type CatégorieUsager int

const (
	Conducteur CatégorieUsager = iota
	Passager
	Piéton
)

func (catégorieUsager CatégorieUsager) String() string {
	return [...]string{
		"Conducteur",
		"Passager",
		"Piéton",
	}[catégorieUsager]
}

func (catégorieUsager CatégorieUsager) MarshalJSON() ([]byte, error) {
	return json.Marshal(catégorieUsager.String())
}

type Gravité int

const (
	GravitéNonRenseignée = iota
	Indemne
	Tué
	BlesséHospitalisé
	BlesséLéger
)

func (gravité Gravité) String() string {
	return [...]string{
		"Non renseigné",
		"Indemne",
		"Tué",
		"Blessé hospitalisé",
		"Blessé léger",
	}[gravité]
}

func (gravité Gravité) MarshalJSON() ([]byte, error) {
	return json.Marshal(gravité.String())
}

type Sexe int

const (
	SexeNonRenseigné = iota
	Masculin
	Féminin
)

func (sexe Sexe) String() string {
	return [...]string{
		"Non renseigné",
		"Masculin",
		"Féminin",
	}[sexe]
}

func (sexe Sexe) MarshalJSON() ([]byte, error) {
	return json.Marshal(sexe.String())
}

type CatégorieVéhicule int

const (
	CatégorieVéhiculeIndéterminable CatégorieVéhicule = iota
	Bicyclette
	Scooter
	Motocyclette
	VéhiculeLéger
	VéhiculeUtilitaire
	PoidsLourd
	Autobus
	Autocar
	Train
	Tramway
	AutreVéhicule
)

func (catégorieVéhicule CatégorieVéhicule) String() string {
	return [...]string{
		"Indéterminable",
		"Bicyclette",
		"Scooter",
		"Motocyclette",
		"Véhicule léger",
		"Véhicule utilitaire",
		"Poids lourd",
		"Autobus",
		"Autocar",
		"Train",
		"Tramway",
		"Autre véhicule",
	}[catégorieVéhicule]
}

func (catégorieVéhicule CatégorieVéhicule) MarshalJSON() ([]byte, error) {
	return json.Marshal(catégorieVéhicule.String())
}

type Lieu struct {
	IdAccident   string
	VoieSpéciale VoieSpéciale
}

func (lieu Lieu) AsJson() (string, error) {
	return ToJson(lieu)
}

type Usager struct {
	IdVéhicule      string
	IdAccident      string
	CatégorieUsager CatégorieUsager
	Gravité         Gravité
	Sexe            Sexe
	AnnéeNaissance  int
}

func (usager Usager) AsJson() (string, error) {
	return ToJson(usager)
}

type Véhicule struct {
	IdVéhicule        string
	IdAccident        string
	CatégorieVéhicule CatégorieVéhicule
	Usagers           []*Usager
}

func (véhicule Véhicule) AsJson() (string, error) {
	return ToJson(véhicule)
}

type Accident struct {
	IdAccident    string
	Date          string
	Département   string
	Commune       *int
	Adresse       string
	Latitude      string
	Longitude     string
	Lieu          *Lieu
	Véhicules     []*Véhicule
	AutresUsagers []*Usager // Users not associated with a vehicle
}

func (accident Accident) AsJson() (string, error) {
	return ToJson(accident)
}

func ToJson(obj any) (string, error) {
	jsonBytes, err := json.MarshalIndent(obj, "", "  ")

	if err != nil {
		return "", err
	} else {
		return string(jsonBytes), nil
	}
}
