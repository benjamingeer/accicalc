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

		if err != nil {
			return nil, err
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

		d??partement, err := readColumn(row, "dep", path)

		if err != nil {
			return nil, err
		}

		var correctedD??partment string

		switch d??partement {
		case "201":
			correctedD??partment = "2A"

		case "202":
			correctedD??partment = "2B"

		default:
			if len(d??partement) == 3 && d??partement[2] == '0' {
				correctedD??partment = d??partement[0:2]
			} else {
				correctedD??partment = d??partement
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

		latitude := parseFixedPointLatLong(latitudeStr, idAccident, "lat")

		longitudeStr, err := readColumn(row, "long", path)

		if err != nil {
			return nil, err
		}

		longitude := parseFixedPointLatLong(longitudeStr, idAccident, "long")

		return &Accident{
			IdAccident:  idAccident,
			Date:        fmt.Sprintf("%04d-%02d-%02dT%v", correctedAnInt, moisInt, jourInt, isoTime),
			D??partement: correctedD??partment,
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

		voieSp??cialeStr, err := readColumn(row, "vosp", path)

		if err != nil {
			return nil, err
		}

		var voieSp??ciale VoieSp??ciale

		if voieSp??cialeStr == "" {
			voieSp??ciale = VoieSp??cialeNonRenseign??e
		} else {
			voieSp??cialeInt, err := strconv.Atoi(voieSp??cialeStr)

			if err != nil {
				return nil, fmt.Errorf(
					"can't parse column 'vosp' with value '%v' for accident %v in %v",
					voieSp??cialeStr,
					idAccident,
					path,
				)
			}

			switch voieSp??cialeInt {
			// 0 ??? Sans objet
			case 0:
				voieSp??ciale = VoieSp??cialeSansObjet

			// 1 ??? Piste cyclable
			case 1:
				voieSp??ciale = PisteCyclable

			// 2 ??? Bande cyclable
			case 2:
				voieSp??ciale = BandeCyclable

			// 3 ??? Voie r??serv??e
			case 3:
				voieSp??ciale = AutreVoieSp??ciale

			default:
				return nil, fmt.Errorf(
					"can't parse column 'vosp' with value '%v' for accident %v in %v",
					voieSp??cialeInt,
					idAccident,
					path,
				)
			}
		}

		return &Lieu{
			IdAccident:   idAccident,
			VoieSp??ciale: voieSp??ciale,
		}, nil
	}

	return readCsvFile(path, delimiter, convertRow)
}

func (*YearDatasetReader1) ReadVehicles(year uint, dataPath string) (vehicles []*V??hicule, err error) {
	delimiter := ','
	path := filepath.Join(dataPath, fmt.Sprint(year), fmt.Sprintf("vehicules%v", filenameSuffix1(year)))

	convertRow := func(row map[string]string) (*V??hicule, error) {
		idAccident, err := readColumn(row, "Num_Acc", path)

		if err != nil {
			return nil, err
		}

		idV??hicule, err := readColumn(row, "num_veh", path)

		if err != nil {
			return nil, err
		}

		cat??gorieV??hiculeStr, err := readColumn(row, "catv", path)

		if err != nil {
			return nil, err
		}

		cat??gorieV??hiculeInt, err := strconv.Atoi(cat??gorieV??hiculeStr)

		if err != nil {
			return nil, fmt.Errorf(
				"can't parse column 'catv' with value '%v' for accident %v in %v",
				cat??gorieV??hiculeStr,
				idAccident,
				path,
			)
		}

		var cat??gorieV??hicule Cat??gorieV??hicule

		switch cat??gorieV??hiculeInt {
		// 00 ??? Ind??terminable
		case 0:
			cat??gorieV??hicule = Cat??gorieV??hiculeInd??terminable

		// 01 ??? Bicyclette
		case 1:
			cat??gorieV??hicule = Bicyclette

		// 04 ??? R??f??rence inutilis??e depuis 2006 (scooter immatricul??)
		// 30 ??? Scooter < 50 cm3
		// 32 ??? Scooter > 50 cm3 et <= 125 cm3
		// 34 ??? Scooter > 125 cm3
		case 4, 30, 32, 34:
			cat??gorieV??hicule = Scooter

		// 05 ??? R??f??rence inutilis??e depuis 2006 (motocyclette)
		// 31 ??? Motocyclette > 50 cm3 et <= 125 cm3
		// 33 ??? Motocyclette > 125 cm3
		case 5, 31, 33:
			cat??gorieV??hicule = Motocyclette

		// 07 ??? VL seul
		// 08 ??? R??f??rence inutilis??e depuis 2006 (VL + caravane)
		// 09 ??? R??f??rence inutilis??e depuis 2006 (VL + remorque)
		case 7, 8, 9:
			cat??gorieV??hicule = V??hiculeL??ger

		// 10 ??? VU seul 1,5T <= PTAC <= 3,5T avec ou sans remorque (anciennement VU seul 1,5T <= PTAC <= 3,5T)
		// 11 ??? R??f??rence inutilis??e depuis 2006 (VU (10) + caravane)
		// 12 ??? R??f??rence inutilis??e depuis 2006 (VU (10) + remorque)
		case 10, 11, 12:
			cat??gorieV??hicule = V??hiculeUtilitaire

		// 13 ??? PL seul 3,5T <PTCA <= 7,5T
		// 14 ??? PL seul > 7,5T
		// 15 ??? PL > 3,5T + remorque
		case 13, 14, 15:
			cat??gorieV??hicule = PoidsLourd

		// 37 ??? Autobus
		case 37:
			cat??gorieV??hicule = Autobus

		// 38 ??? Autocar
		case 38:
			cat??gorieV??hicule = Autocar

		// 39 ??? Train
		case 39:
			cat??gorieV??hicule = Train

		// 40 ??? Tramway
		case 40:
			cat??gorieV??hicule = Tramway

		default:
			cat??gorieV??hicule = AutreV??hicule
		}

		return &V??hicule{
			IdV??hicule:        idV??hicule,
			IdAccident:        idAccident,
			Cat??gorieV??hicule: cat??gorieV??hicule,
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

		idV??hicule, err := readColumn(row, "num_veh", path)

		if err != nil {
			return nil, err
		}

		cat??gorieUsagerStr, err := readColumn(row, "catu", path)

		if err != nil {
			return nil, err
		}

		cat??gorieUsagerInt, err := strconv.Atoi(cat??gorieUsagerStr)

		if err != nil {
			return nil, fmt.Errorf(
				"can't parse column 'catu' with value '%v' for accident %v in %v",
				cat??gorieUsagerStr,
				idAccident,
				path,
			)
		}

		var cat??gorieUsager Cat??gorieUsager

		switch cat??gorieUsagerInt {
		// 1 ??? Conducteur
		// 4 - Pi??ton en roller ou en trottinette (cat??gorie d??plac??e, ?? partir de l???ann??e 2018, vers le fichier
		// "V??hicules" Cat??gorie du v??hicule : 99 - Autre v??hicule. Cette cat??gorie est d??sormais consid??r??e comme
		// un v??hicule : engin de d??placement personnel)
		case 1, 4:
			cat??gorieUsager = Conducteur

		// 2 ??? Passager
		case 2:
			cat??gorieUsager = Passager

		// 3 ??? Pi??ton
		case 3:
			cat??gorieUsager = Pi??ton

		default:
			return nil, fmt.Errorf(
				"can't parse column 'catu' with value '%v' for accident %v in %v",
				cat??gorieUsagerInt,
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
				cat??gorieUsagerStr,
				idAccident,
				path,
			)
		}

		var sexe Sexe

		switch sexeInt {
		// 1 ??? Masculin
		case 1:
			sexe = Masculin

		// 2 ??? F??minin
		case 2:
			sexe = F??minin

		default:
			// -1 seems to be used for this in the data
			sexe = SexeNonRenseign??
		}

		gravit??Str, err := readColumn(row, "grav", path)

		if err != nil {
			return nil, err
		}

		gravit??Int, err := strconv.Atoi(gravit??Str)

		if err != nil {
			return nil, fmt.Errorf(
				"can't parse column 'grav' with value '%v' for accident %v in %v",
				gravit??Str,
				idAccident,
				path,
			)
		}

		var gravit?? Gravit??

		switch gravit??Int {
		// 1 ??? Indemne
		case 1:
			gravit?? = Indemne

		// 2 ??? Tu??
		case 2:
			gravit?? = Tu??

		// 3 ??? Bless?? hospitalis??
		case 3:
			gravit?? = Bless??Hospitalis??

		// 4 ??? Bless?? l??ger
		case 4:
			gravit?? = Bless??L??ger

		default:
			// -1 seems to be used for this in the data
			gravit?? = Gravit??NonRenseign??e
		}

		ann??eNaissanceStr, err := readColumn(row, "an_nais", path)

		if err != nil {
			return nil, err
		}

		var ann??eNaissance int

		if ann??eNaissanceStr != "" {
			ann??eNaissance, err = strconv.Atoi(ann??eNaissanceStr)

			if err != nil {
				return nil, fmt.Errorf(
					"can't parse column 'an_nais' with value '%v' for accident %v in %v",
					ann??eNaissanceStr,
					idAccident,
					path,
				)
			}
		}

		return &Usager{
			IdV??hicule:      idV??hicule,
			IdAccident:      idAccident,
			Cat??gorieUsager: cat??gorieUsager,
			Sexe:            sexe,
			Gravit??:         gravit??,
			Ann??eNaissance:  ann??eNaissance,
		}, nil
	}

	return readCsvFile(path, delimiter, convertRow)
}

func parseFixedPointLatLong(latLongStr string, idAccident string, colName string) string {
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
