package dataset

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type YearDatasetReader1 struct{}

func filenameSuffix1(year uint) string {
	if year <= 2016 {
		return fmt.Sprintf("_%v.csv", year)
	} else {
		return fmt.Sprintf("-%v.csv", year)
	}
}

func (*YearDatasetReader1) ReadCharacteristics(year uint, dataPath string) (accidents []*Accident, err error) {
	var delimiter rune

	if year == 2009 {
		delimiter = '\t'
	} else {
		delimiter = ','
	}

	path := filepath.Join(dataPath, fmt.Sprint(year), fmt.Sprintf("caracteristiques%v", filenameSuffix1(year)))

	convertRow := func(row map[string]string) (*Accident, error) {
		idAccident, err := readColumn(row, "Num_Acc", path)

		if err != nil {
			return nil, err
		}

		jourStr, err := readColumn(row, "jour", path)

		if err != nil {
			return nil, err
		}

		jourInt, err := strconv.Atoi(jourStr)

		if err != nil {
			return nil, fmt.Errorf(
				"can't parse column 'jour' with value '%v' for accident %v in %v",
				jourStr,
				idAccident,
				path,
			)
		}

		moisStr, err := readColumn(row, "mois", path)

		if err != nil {
			return nil, err
		}

		moisInt, err := strconv.Atoi(moisStr)

		if err != nil {
			return nil, fmt.Errorf(
				"can't parse column 'mois' with value '%v' for accident %v in %v",
				moisStr,
				idAccident,
				path,
			)
		}

		anStr, err := readColumn(row, "an", path)

		if err != nil {
			return nil, err
		}

		anInt, err := strconv.Atoi(anStr)

		if err != nil {
			return nil, fmt.Errorf(
				"can't parse column 'an' with value '%v' for accident %v in %v",
				anStr,
				idAccident,
				path,
			)
		}

		correctedAnInt := anInt + 2000

		hrmn, err := readColumn(row, "hrmn", path)

		if err != nil {
			return nil, err
		}

		var hourStr string
		var minuteStr string

		switch len(hrmn) {
		case 1, 2:
			hourStr = "0"
			minuteStr = hrmn

		case 3:
			hourStr = hrmn[0:1]
			minuteStr = hrmn[1:]

		case 4:
			hourStr = hrmn[0:2]
			minuteStr = hrmn[2:]
		}

		var hour int

		hour, err = strconv.Atoi(hourStr)

		if err != nil {
			return nil, fmt.Errorf(
				"can't parse column 'hrmn' with value '%v' for accident %v in %v",
				hrmn,
				idAccident,
				path,
			)
		}

		minute, err := strconv.Atoi(minuteStr)

		if err != nil {
			return nil, fmt.Errorf(
				"can't parse column 'hrmn' with value '%v' for accident %v in %v",
				hrmn,
				idAccident,
				path,
			)
		}

		isoTime := fmt.Sprintf("%02d:%02d", hour, minute)

		département, err := readColumn(row, "dep", path)

		if err != nil {
			return nil, err
		}

		var correctedDépartment string

		switch département {
		case "201":
			correctedDépartment = "2A"

		case "202":
			correctedDépartment = "2B"

		default:
			if len(département) == 3 && département[2] == '0' {
				correctedDépartment = département[0:2]
			} else {
				correctedDépartment = département
			}
		}

		communeStr, err := readColumn(row, "com", path)

		if err != nil {
			return nil, err
		}

		var commune *int

		if communeStr != "" {
			communeInt, err := strconv.Atoi(communeStr)

			if err != nil {
				return nil, fmt.Errorf(
					"can't parse column 'com' with value '%v' for accident %v in %v",
					communeStr,
					idAccident,
					path,
				)
			}

			commune = &communeInt
		}

		adresse, err := readColumn(row, "adr", path)

		if err != nil {
			return nil, err
		}

		latitudeStr, err := readColumn(row, "lat", path)

		if err != nil {
			return nil, err
		}

		latitude := parseFixedPointLatLong(latitudeStr)

		longitudeStr, err := readColumn(row, "long", path)

		if err != nil {
			return nil, err
		}

		longitude := parseFixedPointLatLong(longitudeStr)

		return &Accident{
			IdAccident:  idAccident,
			Date:        fmt.Sprintf("%04d-%02d-%02dT%v", correctedAnInt, moisInt, jourInt, isoTime),
			Département: correctedDépartment,
			Commune:     commune,
			Adresse:     adresse,
			Latitude:    latitude,
			Longitude:   longitude,
		}, nil
	}

	return readCsvFile(path, delimiter, convertRow)
}

func (*YearDatasetReader1) ReadPlaces(year uint, dataPath string) (places []*Lieu, err error) {
	delimiter := ','
	path := filepath.Join(dataPath, fmt.Sprint(year), fmt.Sprintf("lieux%v", filenameSuffix1(year)))

	convertRow := func(row map[string]string) (*Lieu, error) {
		idAccident, err := readColumn(row, "Num_Acc", path)

		if err != nil {
			fmt.Printf("row: %v\n", row)
			return nil, err
		}

		voieSpécialeStr, err := readColumn(row, "vosp", path)

		if err != nil {
			return nil, err
		}

		var voieSpéciale VoieSpéciale

		if voieSpécialeStr == "" {
			voieSpéciale = VoieSpécialeNonRenseignée
		} else {
			voieSpécialeInt, err := strconv.Atoi(voieSpécialeStr)

			if err != nil {
				return nil, fmt.Errorf(
					"can't parse column 'vosp' with value '%v' for accident %v in %v",
					voieSpécialeStr,
					idAccident,
					path,
				)
			}

			switch voieSpécialeInt {
			// 0 – Sans objet
			case 0:
				voieSpéciale = VoieSpécialeSansObjet

			// 1 – Piste cyclable
			case 1:
				voieSpéciale = PisteCyclable

			// 2 – Bande cyclable
			case 2:
				voieSpéciale = BandeCyclable

			// 3 – Voie réservée
			case 3:
				voieSpéciale = AutreVoieSpéciale

			default:
				return nil, fmt.Errorf(
					"can't parse column 'vosp' with value '%v' for accident %v in %v",
					voieSpécialeInt,
					idAccident,
					path,
				)
			}
		}

		return &Lieu{
			IdAccident:   idAccident,
			VoieSpéciale: voieSpéciale,
		}, nil
	}

	return readCsvFile(path, delimiter, convertRow)
}

func (*YearDatasetReader1) ReadVehicles(year uint, dataPath string) (vehicles []*Véhicule, err error) {
	delimiter := ','
	path := filepath.Join(dataPath, fmt.Sprint(year), fmt.Sprintf("vehicules%v", filenameSuffix1(year)))

	convertRow := func(row map[string]string) (*Véhicule, error) {
		idAccident, err := readColumn(row, "Num_Acc", path)

		if err != nil {
			return nil, err
		}

		idVéhicule, err := readColumn(row, "num_veh", path)

		if err != nil {
			return nil, err
		}

		catégorieVéhiculeStr, err := readColumn(row, "catv", path)

		if err != nil {
			return nil, err
		}

		catégorieVéhiculeInt, err := strconv.Atoi(catégorieVéhiculeStr)

		if err != nil {
			return nil, fmt.Errorf(
				"can't parse column 'catv' with value '%v' for accident %v in %v",
				catégorieVéhiculeStr,
				idAccident,
				path,
			)
		}

		var catégorieVéhicule CatégorieVéhicule

		switch catégorieVéhiculeInt {
		// 00 – Indéterminable
		case 0:
			catégorieVéhicule = CatégorieVéhiculeIndéterminable

		// 01 – Bicyclette
		case 1:
			catégorieVéhicule = Bicyclette

		// 04 – Référence inutilisée depuis 2006 (scooter immatriculé)
		// 30 – Scooter < 50 cm3
		// 32 – Scooter > 50 cm3 et <= 125 cm3
		// 34 – Scooter > 125 cm3
		case 4, 30, 32, 34:
			catégorieVéhicule = Scooter

		// 05 – Référence inutilisée depuis 2006 (motocyclette)
		// 31 – Motocyclette > 50 cm3 et <= 125 cm3
		// 33 – Motocyclette > 125 cm3
		case 5, 31, 33:
			catégorieVéhicule = Motocyclette

		// 07 – VL seul
		// 08 – Référence inutilisée depuis 2006 (VL + caravane)
		// 09 – Référence inutilisée depuis 2006 (VL + remorque)
		case 7, 8, 9:
			catégorieVéhicule = VéhiculeLéger

		// 10 – VU seul 1,5T <= PTAC <= 3,5T avec ou sans remorque (anciennement VU seul 1,5T <= PTAC <= 3,5T)
		// 11 – Référence inutilisée depuis 2006 (VU (10) + caravane)
		// 12 – Référence inutilisée depuis 2006 (VU (10) + remorque)
		case 10, 11, 12:
			catégorieVéhicule = VéhiculeUtilitaire

		// 13 – PL seul 3,5T <PTCA <= 7,5T
		// 14 – PL seul > 7,5T
		// 15 – PL > 3,5T + remorque
		case 13, 14, 15:
			catégorieVéhicule = PoidsLourd

		// 37 – Autobus
		case 37:
			catégorieVéhicule = Autobus

		// 38 – Autocar
		case 38:
			catégorieVéhicule = Autocar

		// 39 – Train
		case 39:
			catégorieVéhicule = Train

		// 40 – Tramway
		case 40:
			catégorieVéhicule = Tramway

		default:
			catégorieVéhicule = AutreVéhicule
		}

		return &Véhicule{
			IdVéhicule:        idVéhicule,
			IdAccident:        idAccident,
			CatégorieVéhicule: catégorieVéhicule,
		}, nil
	}

	return readCsvFile(path, delimiter, convertRow)
}

func (*YearDatasetReader1) ReadUsers(year uint, dataPath string) (users []*Usager, err error) {
	delimiter := ','
	path := filepath.Join(dataPath, fmt.Sprint(year), fmt.Sprintf("usagers%v", filenameSuffix1(year)))

	convertRow := func(row map[string]string) (*Usager, error) {
		idAccident, err := readColumn(row, "Num_Acc", path)

		if err != nil {
			return nil, err
		}

		idVéhicule, err := readColumn(row, "num_veh", path)

		if err != nil {
			return nil, err
		}

		catégorieUsagerStr, err := readColumn(row, "catu", path)

		if err != nil {
			return nil, err
		}

		catégorieUsagerInt, err := strconv.Atoi(catégorieUsagerStr)

		if err != nil {
			return nil, fmt.Errorf(
				"can't parse column 'catu' with value '%v' for accident %v in %v",
				catégorieUsagerStr,
				idAccident,
				path,
			)
		}

		var catégorieUsager CatégorieUsager

		switch catégorieUsagerInt {
		// 1 – Conducteur
		// 4 - Piéton en roller ou en trottinette (catégorie déplacée, à partir de l’année 2018, vers le fichier
		// "Véhicules" Catégorie du véhicule : 99 - Autre véhicule. Cette catégorie est désormais considérée comme
		// un véhicule : engin de déplacement personnel)
		case 1, 4:
			catégorieUsager = Conducteur

		// 2 – Passager
		case 2:
			catégorieUsager = Passager

		// 3 – Piéton
		case 3:
			catégorieUsager = Piéton

		default:
			return nil, fmt.Errorf(
				"can't parse column 'catu' with value '%v' for accident %v in %v",
				catégorieUsagerInt,
				idAccident,
				path,
			)
		}

		sexeStr, err := readColumn(row, "sexe", path)

		if err != nil {
			return nil, err
		}

		sexeInt, err := strconv.Atoi(sexeStr)

		if err != nil {
			return nil, fmt.Errorf(
				"can't parse column 'sexe' with value '%v' for accident %v in %v",
				catégorieUsagerStr,
				idAccident,
				path,
			)
		}

		var sexe Sexe

		switch sexeInt {
		// 1 – Masculin
		case 1:
			sexe = Masculin

		// 2 – Féminin
		case 2:
			sexe = Féminin

		default:
			// -1 seems to be used for this in the data
			sexe = SexeNonRenseigné
		}

		gravitéStr, err := readColumn(row, "grav", path)

		if err != nil {
			return nil, err
		}

		gravitéInt, err := strconv.Atoi(gravitéStr)

		if err != nil {
			return nil, fmt.Errorf(
				"can't parse column 'grav' with value '%v' for accident %v in %v",
				gravitéStr,
				idAccident,
				path,
			)
		}

		var gravité Gravité

		switch gravitéInt {
		// 1 – Indemne
		case 1:
			gravité = Indemne

		// 2 – Tué
		case 2:
			gravité = Tué

		// 3 – Blessé hospitalisé
		case 3:
			gravité = BlesséHospitalisé

		// 4 – Blessé léger
		case 4:
			gravité = BlesséLéger

		default:
			// -1 seems to be used for this in the data
			gravité = GravitéNonRenseignée
		}

		annéeNaissanceStr, err := readColumn(row, "an_nais", path)

		if err != nil {
			return nil, err
		}

		var annéeNaissance int

		if annéeNaissanceStr != "" {
			annéeNaissance, err = strconv.Atoi(annéeNaissanceStr)

			if err != nil {
				return nil, fmt.Errorf(
					"can't parse column 'an_nais' with value '%v' for accident %v in %v",
					annéeNaissanceStr,
					idAccident,
					path,
				)
			}
		}

		return &Usager{
			IdVéhicule:      idVéhicule,
			IdAccident:      idAccident,
			CatégorieUsager: catégorieUsager,
			Sexe:            sexe,
			Gravité:         gravité,
			AnnéeNaissance:  annéeNaissance,
		}, nil
	}

	return readCsvFile(path, delimiter, convertRow)
}

func parseFixedPointLatLong(latLongStr string) string {
	if len(latLongStr) < 2 {
		return ""
	} else {
		latLongWithComma := latLongStr[0:2] + "," + latLongStr[2:]
		latLongWithCommaTrimmedRight := strings.TrimRight(latLongWithComma, "0")

		if latLongWithCommaTrimmedRight[0:1] == "0" {
			return latLongWithComma[1:]
		} else {
			return latLongWithComma
		}
	}
}
