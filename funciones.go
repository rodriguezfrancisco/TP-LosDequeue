package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	gocsv "github.com/gocarina/gocsv"
)

func agregarEjercicioARutina(rutina *Rutina, ejercicio Ejercicio) {
	rutina.Ejercicios = append(rutina.Ejercicios, ejercicio)
	rutina.DuracionTotal += ejercicio.Duracion
	rutina.PuntosCardioTotales += ejercicio.PuntosCardio
	rutina.PuntosFlexibilidadTotales += ejercicio.PuntosFlexibilidad
	rutina.PuntosFuerzaTotales += ejercicio.PuntosFuerza

}

// Función para agregar un ejercicio a una rutina existente
func agregarEjercicioARutinaExistente(nombreDeRut int) {

	rutinasFile, err := os.OpenFile("rutinas.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer rutinasFile.Close()
	rutinasCsv := []*RutinaCsv{}
	if err := gocsv.UnmarshalFile(rutinasFile, &rutinasCsv); err != nil { // Load rutinas from file
		panic(err)
	}
	rutinaSeleccionada := rutinasCsv[nombreDeRut]

	// Obtener la rutina seleccionada
	// rutinaSeleccionada := &rutinasL[nombreDeRut]
	var rutinaNueva Rutina
	// Listar categorías disponibles para seleccionar un ejercicio
	greenPrintf("\nSeleccione una categoría para elegir un ejercicio:\n")
	for categoria, ejercicios := range categorias {
		bluePrintf("- %s (%d ejercicios)\n", categoria, len(ejercicios))
	}

	// Escanear la selección de categoría
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	categoriaSeleccionada := scanner.Text()
	categoriaSeleccionada = strings.ToLower(categoriaSeleccionada)
	categoriaSeleccionada = strings.ReplaceAll(categoriaSeleccionada, " ", "")

	// Verificar si la categoría existe
	ejercicios, ok := categorias[categoriaSeleccionada]
	if !ok {
		redPrintf("\nCategoría no válida.\n")
		return
	}

	// Listar ejercicios en la categoría seleccionada
	greenPrintf("\nSeleccione un ejercicio de la categoría %v : (inserte el numero de ejercicio)\n", categoriaSeleccionada)
	for i, ejercicio := range ejercicios {
		bluePrintf("%d.", i+1)
		fmt.Printf("%s (Duración: %d Segundos, Intensidad: %s)\n", ejercicio.Nombre, ejercicio.Duracion, ejercicio.Intensidad)
	}

	// Escanear la selección de ejercicio
	scanner.Scan()
	seleccionEjercicioStr := scanner.Text()
	seleccionEjercicio, err := strconv.Atoi(seleccionEjercicioStr)
	if err != nil || seleccionEjercicio < 1 || seleccionEjercicio > len(ejercicios) {
		redPrintf("\nSelección de ejercicio inválida. Introduzca el numero de ejercicio correctamente\n")
		return
	}

	// Obtener el ejercicio seleccionado
	ejercicioSeleccionado := ejercicios[seleccionEjercicio-1]
	rutinaNueva.DuracionTotal = rutinaSeleccionada.DuracionTotal
	rutinaNueva.PuntosCardioTotales = rutinaSeleccionada.PuntosCardioTotales
	rutinaNueva.PuntosFuerzaTotales = rutinaSeleccionada.PuntosFuerzaTotales
	rutinaNueva.PuntosFlexibilidadTotales = rutinaSeleccionada.PuntosFlexibilidadTotales
	rutinaNueva.Ejercicios = rutinaSeleccionada.Ejercicios
	rutinaNueva.NombreDeRutina = rutinaSeleccionada.NombreDeRutina + "*"
	// Agregar el ejercicio a la rutina seleccionada
	agregarEjercicioARutina(&rutinaNueva, ejercicioSeleccionado)

	fmt.Printf("\nEjercicio '%s' agregado. La rutina modificada se llama '%s'. Duración total de la rutina: %d Segundos.\n",
		ejercicioSeleccionado.Nombre, rutinaNueva.NombreDeRutina, rutinaNueva.DuracionTotal)

	//Start save in csv
	rutinasFile, err = os.OpenFile("rutinas.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer rutinasFile.Close()
	rutinasCsv = []*RutinaCsv{}
	if err := gocsv.UnmarshalFile(rutinasFile, &rutinasCsv); err != nil { // Load rutinas from file
		panic(err)
	}

	if _, err := rutinasFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}

	rutinasCreadas := fmt.Sprintf("%v", len(rutinasCsv)+1)

	rutinasCsv = append(rutinasCsv, &RutinaCsv{Id: rutinasCreadas, NombreDeRutina: rutinaNueva.NombreDeRutina, Ejercicios: rutinaNueva.Ejercicios, DuracionTotal: rutinaNueva.DuracionTotal, PuntosCardioTotales: rutinaNueva.PuntosCardioTotales, PuntosFlexibilidadTotales: rutinaNueva.PuntosFlexibilidadTotales, PuntosFuerzaTotales: rutinaNueva.PuntosFuerzaTotales}) // Add rutinas

	err = gocsv.MarshalFile(&rutinasCsv, rutinasFile) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}
}

// Función para agregar un ejercicio a una categoría.
func agregarEjercicioACategoria(nombre, tipo string, duracion int, intensidad string, calorias int, descripcion string, PuntosCardio int, PuntosFlexibilidad int, PuntosFuerza int) {
	/*nombre,*/ tipo = /*strings.ToLower(nombre),*/ strings.ToLower(tipo)
	/*nombre,*/ tipo = /*strings.ReplaceAll(nombre, " ", ""),*/ strings.ReplaceAll(tipo, " ", "")
	ej := Ejercicio{Nombre: nombre, Duracion: duracion, Tipo: tipo, Intensidad: intensidad, Calorias: calorias, Descripcion: descripcion, PuntosCardio: PuntosCardio, PuntosFlexibilidad: PuntosFlexibilidad, PuntosFuerza: PuntosFuerza}
	categorias[tipo] = append(categorias[tipo], ej)
	// fmt.Println(": ", categorias)
}

// Función para solicitar los detalles de un nuevo ejercicio y agregarlo a una categoría.
func solicitarYAgregarEjercicio() {
	scanner := bufio.NewScanner(os.Stdin)

	// Solicitar detalles del ejercicio
	fmt.Print("Ingrese el nombre del ejercicio: ")
	scanner.Scan()
	nombre := scanner.Text()

	fmt.Print("Ingrese la duración del ejercicio (en Segundos): ")
	scanner.Scan()
	duracionStr := scanner.Text()
	duracion, err := strconv.Atoi(duracionStr)
	if err != nil {
		fmt.Println("Duración inválida. Por favor, ingrese un número válido.")
		return
	}

	fmt.Print("Ingrese el tipo de ejercicio (Flexibilidad, Cardio, Fuerza): ")
	scanner.Scan()
	tipo := scanner.Text()
	tipo = strings.ToLower(tipo)
	tipo = strings.ReplaceAll(tipo, " ", "")
	if tipo != "flexibilidad" && tipo != "cardio" && tipo != "fuerza" {
		for tipo != "flexibilidad" && tipo != "cardio" && tipo != "fuerza" {
			redPrintf("\nEl tipo de ejercicio ingresado no existe, debe ser uno de los existentes\n")
			fmt.Print("Ingrese el tipo de ejercicio (Flexibilidad, Cardio, Fuerza): ")
			scanner.Scan()
			tipo = scanner.Text()
			tipo = strings.ToLower(tipo)
			tipo = strings.ReplaceAll(tipo, " ", "")
		}
	}
	fmt.Print("Ingrese la intensidad del ejercicio (Baja, Media, Alta): ")
	scanner.Scan()
	intensidad := scanner.Text()
	intensidad = strings.ToLower(intensidad)
	intensidad = strings.ReplaceAll(intensidad, " ", "")
	if intensidad != "baja" && intensidad != "media" && intensidad != "alta" {
		for intensidad != "baja" && intensidad != "media" && intensidad != "alta" {
			redPrintf("\nLa intensidad del ejercicio ingresado no existe, debe ser una de las existentes\n")
			fmt.Print("Ingrese la intensidad del ejercicio (Baja, Media, Alta): ")
			scanner.Scan()
			intensidad = scanner.Text()
			intensidad = strings.ToLower(intensidad)
			intensidad = strings.ReplaceAll(intensidad, " ", "")
		}
	}

	fmt.Print("Ingrese la cantidad de calorias quemadas: ")
	scanner.Scan()
	caloriasStr := scanner.Text()
	calorias, err := strconv.Atoi(caloriasStr)
	if err != nil {
		fmt.Println("Descripcion invalida")
		return
	}

	fmt.Print("Ingrese la descripción del ejercicio: ")
	scanner.Scan()
	descripcion := scanner.Text()

	fmt.Print("Ingrese los puntos de cardio del ejercicio (en numeros (del 0 al 100) o tendrá que hacer todo devuelta): ")
	scanner.Scan()
	PuntosCardiostr := scanner.Text()
	PuntosCardio, err := strconv.Atoi(PuntosCardiostr)
	if err != nil || PuntosCardio < 0 || PuntosCardio > 100 {
		redPrintf("\nPuntuacion invalida. Por favor, ingrese un número válido.\n")
		return
	}

	fmt.Print("Ingrese los puntos de flexibilidad del ejercicio (en numeros (del 0 al 100) o tendrá que hacer todo devuelta): ")
	scanner.Scan()
	PuntosFlexibilidadStr := scanner.Text()
	PuntosFlexibilidad, err := strconv.Atoi(PuntosFlexibilidadStr)
	if err != nil || PuntosFlexibilidad < 0 || PuntosFlexibilidad > 100 {
		redPrintf("\nPuntuacion inválida. Por favor, ingrese un número válido.\n")
		return
	}

	fmt.Print("Ingrese los puntos de fuerza del ejercicio (en numeros (del 0 al 100) o tendrá que hacer todo devuelta): ")
	scanner.Scan()
	PuntosFuerzaStr := scanner.Text()
	PuntosFuerza, err := strconv.Atoi(PuntosFuerzaStr)
	if err != nil || PuntosFuerza < 0 || PuntosFuerza > 100 {
		redPrintf("\nPuntuacion inválida. Por favor, ingrese un número válido.\n")
		return
	}

	// Llamar a agregarEjercicioACategoria para agregar el ejercicio
	agregarEjercicioACategoria(nombre, tipo, duracion, intensidad, calorias, descripcion, PuntosCardio, PuntosFlexibilidad, PuntosFuerza)

	// Confirmación
	fmt.Printf("Ejercicio '%s' agregado correctamente a la categoría '%s'.\n", nombre, tipo)
}

// Función para crear una rutina personalizada.
func crearRutinaPersonalizada() Rutina {
	var rutina Rutina

	scanner := bufio.NewScanner(os.Stdin)
	continuar := true

	greenPrintf("\n\n\n\n\n\n\nIngrese un nombre para la rutina:")
	bluePrintf("(ej: Rutina para Elias)\n\n\n\n\n___________________________________\n")

	scanner.Scan()
	rutina.NombreDeRutina = scanner.Text()
	rutina.NombreDeRutina = strings.ToLower(rutina.NombreDeRutina)
	rutina.NombreDeRutina = strings.ReplaceAll(rutina.NombreDeRutina, " ", "")

	greenPrintf("\n\n\n\n\n\n\nCategorías disponibles:")
	for categoria, ejercicios := range categorias {
		bluePrintf("\n- %s ", strings.Title(categoria))
		redPrintf("(%d ejercicios)\n", len(ejercicios))
	}

	for continuar {
		greenPrintf("\n\n\n\n\n\n\nCategorías disponibles:")
		for categoria, ejercicios := range categorias {
			bluePrintf("\n- %s ", strings.Title(categoria))
			redPrintf("(%d ejercicios)\n", len(ejercicios))
		}

		greenPrintf("\n\n\nSeleccione una categoría para agregar a su rutina")
		bluePrintf("(ej: 'fuerza')")
		greenPrintf(" o escriba ")
		bluePrintf("'listo' ")
		greenPrintf("para finalizar: ")
		bluePrintf("\n\n\n___________________________________\n")
		scanner.Scan()
		categoria := scanner.Text()
		categoria = strings.ReplaceAll(categoria, " ", "")
		categoria = strings.ToLower(categoria)

		if categoria == "listo" {
			fmt.Println("\n\n\n\n\n\n\n\n\n\n\n.")
			continuar = false
			continue
		}

		ejercicios, ok := categorias[categoria]
		if !ok {
			redPrintf("\n\n\n\nLa categoría seleccionada no existe.\n\n Las categorias disponibles son:")
			for categoria, ejercicios := range categorias {
				bluePrintf("\n- %s ", strings.Title(categoria))
				redPrintf("(%d ejercicios)", len(ejercicios))
			}
			continue
		}

		greenPrintf("\n\n\n\n\nOrganizar la lista de ejercicios por:\n 1.Nombre \n 2.Duracion\n 3.Calorias\n\n\n___________________________________\n")
		scanner.Scan()
		filtrado := scanner.Text()
		filtrado = strings.ReplaceAll(filtrado, " ", "")

		switch strings.ToLower(filtrado) {
		case "1":
			greenPrintf("\n\n\n\n\nEjercicios disponibles en la categoría %s:\n\n", categoria)
			for i, ejercicio := range ejercicios {
				bluePrintf("%d.", i+1)
				fmt.Printf(" %s", ejercicio.Nombre)
				redPrintf(" (Duración: %d Segundos, Intensidad: %s)\n", ejercicio.Duracion, ejercicio.Intensidad)
			}
		case "2":
			greenPrintf("\n\n\n\n\nEjercicios disponibles en la categoría %s:\n\n", categoria)
			for i, ejercicio := range ejercicios {
				bluePrintf("%d.", i+1)
				fmt.Printf(" %s", ejercicio.Nombre)
				redPrintf(" (Duración: %d Segundos, Intensidad: %s)\n", ejercicio.Duracion, ejercicio.Intensidad)
			}
		case "3":
			greenPrintf("\n\n\n\n\nEjercicios disponibles en la categoría %s:\n\n", categoria)
			for i, ejercicio := range ejercicios {
				bluePrintf("%d.", i+1)
				fmt.Printf(" %s", ejercicio.Nombre)
				redPrintf(" (Duración: %d Segundos, Intensidad: %s)\n", ejercicio.Duracion, ejercicio.Intensidad)
			}
		case "nombre":
			greenPrintf("\n\n\n\n\nEjercicios disponibles en la categoría %s:\n\n", categoria)
			for i, ejercicio := range ejercicios {
				bluePrintf("%d.", i+1)
				fmt.Printf(" %s", ejercicio.Nombre)
				redPrintf(" (Duración: %d Segundos, Intensidad: %s)\n", ejercicio.Duracion, ejercicio.Intensidad)
			}
		case "duracion":
			greenPrintf("\n\n\n\n\nEjercicios disponibles en la categoría %s:\n\n", categoria)
			for i, ejercicio := range ejercicios {
				bluePrintf("%d.", i+1)
				fmt.Printf(" %s", ejercicio.Nombre)
				redPrintf(" (Duración: %d Segundos, Intensidad: %s)\n", ejercicio.Duracion, ejercicio.Intensidad)
			}
		case "calorias":
			greenPrintf("\n\n\n\n\nEjercicios disponibles en la categoría %s:\n\n", categoria)
			for i, ejercicio := range ejercicios {
				bluePrintf("%d.", i+1)
				fmt.Printf(" %s", ejercicio.Nombre)
				redPrintf(" (Duración: %d Segundos, Intensidad: %s)\n", ejercicio.Duracion, ejercicio.Intensidad)
			}
		}

		greenPrintf("\n\nSeleccione el ejercicio que desea agregar a su rutina o escriba ")
		bluePrintf("(0)")
		greenPrintf(" para cambiar de categoría: ")
		bluePrintf("\n\n\n___________________________________\n")
		scanner.Scan()
		ejercicioIndexStr := scanner.Text()
		ejercicioIndex, err := strconv.Atoi(ejercicioIndexStr)
		if err != nil || ejercicioIndex < 0 || ejercicioIndex > len(ejercicios) {
			fmt.Println("Número de ejercicio inválido.")
			continue
		}
		if ejercicioIndex == 0 {
			continue
		}

		ejercicioSeleccionado := ejercicios[ejercicioIndex-1]
		rutina.Ejercicios = append(rutina.Ejercicios, ejercicioSeleccionado)
		rutina.DuracionTotal += ejercicioSeleccionado.Duracion
		rutina.PuntosCardioTotales += ejercicioSeleccionado.PuntosCardio
		rutina.PuntosFlexibilidadTotales += ejercicioSeleccionado.PuntosFlexibilidad
		rutina.PuntosFuerzaTotales += ejercicioSeleccionado.PuntosFuerza

		fmt.Printf("\n\n\n\n\nSe ha agregado '%s' a su rutina.\n", ejercicioSeleccionado.Nombre)
	}

	fmt.Printf("La duración total de su rutina es de %d Segundos.\n", rutina.DuracionTotal)
	rutinasL = append(rutinasL, rutina)

	//Start save in csv
	rutinasFile, err := os.OpenFile("rutinas.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer rutinasFile.Close()
	rutinasCsv := []*RutinaCsv{}
	if err := gocsv.UnmarshalFile(rutinasFile, &rutinasCsv); err != nil { // Load rutinas from file
		panic(err)
	}

	if _, err := rutinasFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}

	rutinasCreadas := fmt.Sprintf("%v", len(rutinasCsv)+1)

	rutinasCsv = append(rutinasCsv, &RutinaCsv{Id: rutinasCreadas, NombreDeRutina: rutina.NombreDeRutina, Ejercicios: rutina.Ejercicios, DuracionTotal: rutina.DuracionTotal, PuntosCardioTotales: rutina.PuntosCardioTotales, PuntosFlexibilidadTotales: rutina.PuntosFlexibilidadTotales, PuntosFuerzaTotales: rutina.PuntosFuerzaTotales}) // Add rutinas

	err = gocsv.MarshalFile(&rutinasCsv, rutinasFile) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}
	return rutina
}

func consultaRutinaCreada(n int) string {
	rutinasFile, err := os.OpenFile("rutinas.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer rutinasFile.Close()
	rutinasCsv := []*RutinaCsv{}
	if err := gocsv.UnmarshalFile(rutinasFile, &rutinasCsv); err != nil { // Load rutinas from file
		panic(err)
	}
	rutinaString := rutinasCsv[n].String()
	return rutinaString
}

func modificarRutina(nombreDeRut int) {
	for {
		greenPrintf("¿Qué modificación querés hacerle?\n")
		bluePrintf("1.")
		fmt.Print("Agregar ejercicio\n")
		bluePrintf("2.")
		fmt.Print("Editar ejercicio\n")
		bluePrintf("3.")
		fmt.Print("Eliminar ejercicio\n")
		bluePrintf("4.")
		fmt.Print("Volver\n")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		comando := scanner.Text()
		comando = strings.ReplaceAll(comando, " ", "")

		switch strings.ToLower(comando) {
		case "agregarejercicio":
			agregarEjercicioARutinaExistente(nombreDeRut)
		case "1":
			agregarEjercicioARutinaExistente(nombreDeRut)
		case "editarejercicio":
			redPrintf("todavia no disponible. Pruebe eliminando y agregando el ejercicio ya modificado\n")
			continue
		case "2":
			redPrintf("todavia no disponible. Pruebe eliminando y agregando el ejercicio ya modificado\n")
			continue
		case "eliminarejercicio":
			rutinasFile, err := os.OpenFile("rutinas.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				panic(err)
			}
			defer rutinasFile.Close()
			rutinasCsv := []*RutinaCsv{}
			if err := gocsv.UnmarshalFile(rutinasFile, &rutinasCsv); err != nil { // Load rutinas from file
				panic(err)
			}
			rutinaSeleccionada := rutinasCsv[nombreDeRut]
			greenPrintf("¿Qué ejercicio querés eliminar?(insertar el numero)\n")
			for i, rut := range rutinaSeleccionada.Ejercicios {
				bluePrintf("%v.", i+1)
				fmt.Printf("%v\n", rut)
			}
			scanner.Scan()
			seleccionIn := scanner.Text()
			seleccion, _ := strconv.Atoi(seleccionIn)
			if len(rutinaSeleccionada.Ejercicios) < seleccion || len(rutinaSeleccionada.Ejercicios) == 0 {
				fmt.Println("El ejercicio seleccionado no existe")
				continue
			}
			if seleccion == 0 {
				redPrintf("Introduzca correctamente el numero de ejercicio\n")
				continue
			}
			rutinaSeleccionada.DuracionTotal -= rutinaSeleccionada.Ejercicios[seleccion-1].Duracion
			rutinaSeleccionada.PuntosCardioTotales -= rutinaSeleccionada.Ejercicios[seleccion-1].PuntosCardio
			rutinaSeleccionada.PuntosFlexibilidadTotales -= rutinaSeleccionada.Ejercicios[seleccion-1].PuntosFlexibilidad
			rutinaSeleccionada.PuntosFuerzaTotales -= rutinaSeleccionada.Ejercicios[seleccion-1].PuntosFuerza
			j := rutinaSeleccionada.Ejercicios[:seleccion-1]
			j = append(j, rutinaSeleccionada.Ejercicios[seleccion:]...)
			rutinaSeleccionada.Ejercicios = j
			rutinaSeleccionada.NombreDeRutina = rutinaSeleccionada.NombreDeRutina + "*"
			//Start save in csv
			rutinasFile, err = os.OpenFile("rutinas.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				panic(err)
			}
			defer rutinasFile.Close()
			rutinasCsv = []*RutinaCsv{}
			if err := gocsv.UnmarshalFile(rutinasFile, &rutinasCsv); err != nil { // Load rutinas from file
				panic(err)
			}

			if _, err := rutinasFile.Seek(0, 0); err != nil { // Go to the start of the file
				panic(err)
			}

			rutinasCreadas := fmt.Sprintf("%v", len(rutinasCsv)+1)

			rutinasCsv = append(rutinasCsv, &RutinaCsv{Id: rutinasCreadas, NombreDeRutina: rutinaSeleccionada.NombreDeRutina, Ejercicios: rutinaSeleccionada.Ejercicios, DuracionTotal: rutinaSeleccionada.DuracionTotal, PuntosCardioTotales: rutinaSeleccionada.PuntosCardioTotales, PuntosFlexibilidadTotales: rutinaSeleccionada.PuntosFlexibilidadTotales, PuntosFuerzaTotales: rutinaSeleccionada.PuntosFuerzaTotales}) // Add rutinas

			err = gocsv.MarshalFile(&rutinasCsv, rutinasFile) // Use this to save the CSV back to the file
			if err != nil {
				panic(err)
			}
			bluePrintf("Ejercicio eliminado. Nueva rutina %v creada\n", rutinaSeleccionada.NombreDeRutina)
			return
		case "3":
			rutinasFile, err := os.OpenFile("rutinas.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				panic(err)
			}
			defer rutinasFile.Close()
			rutinasCsv := []*RutinaCsv{}
			if err := gocsv.UnmarshalFile(rutinasFile, &rutinasCsv); err != nil { // Load rutinas from file
				panic(err)
			}
			rutinaSeleccionada := rutinasCsv[nombreDeRut]
			greenPrintf("¿Qué ejercicio querés eliminar?(insertar el numero)\n")
			for i, rut := range rutinaSeleccionada.Ejercicios {
				bluePrintf("%v.", i+1)
				fmt.Printf("%v\n", rut)
			}
			scanner.Scan()
			seleccionIn := scanner.Text()
			seleccion, _ := strconv.Atoi(seleccionIn)
			if len(rutinaSeleccionada.Ejercicios) < seleccion || len(rutinaSeleccionada.Ejercicios) == 0 {
				fmt.Println("El ejercicio seleccionado no existe")
				continue
			}
			if seleccion == 0 {
				redPrintf("Introduzca correctamente el numero de ejercicio\n")
				continue
			}
			rutinaSeleccionada.DuracionTotal -= rutinaSeleccionada.Ejercicios[seleccion-1].Duracion
			rutinaSeleccionada.PuntosCardioTotales -= rutinaSeleccionada.Ejercicios[seleccion-1].PuntosCardio
			rutinaSeleccionada.PuntosFlexibilidadTotales -= rutinaSeleccionada.Ejercicios[seleccion-1].PuntosFlexibilidad
			rutinaSeleccionada.PuntosFuerzaTotales -= rutinaSeleccionada.Ejercicios[seleccion-1].PuntosFuerza
			j := rutinaSeleccionada.Ejercicios[:seleccion-1]
			j = append(j, rutinaSeleccionada.Ejercicios[seleccion:]...)
			rutinaSeleccionada.Ejercicios = j
			rutinaSeleccionada.NombreDeRutina = rutinaSeleccionada.NombreDeRutina + "*"
			//Start save in csv
			rutinasFile, err = os.OpenFile("rutinas.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				panic(err)
			}
			defer rutinasFile.Close()
			rutinasCsv = []*RutinaCsv{}
			if err := gocsv.UnmarshalFile(rutinasFile, &rutinasCsv); err != nil { // Load rutinas from file
				panic(err)
			}

			if _, err := rutinasFile.Seek(0, 0); err != nil { // Go to the start of the file
				panic(err)
			}

			rutinasCreadas := fmt.Sprintf("%v", len(rutinasCsv)+1)

			rutinasCsv = append(rutinasCsv, &RutinaCsv{Id: rutinasCreadas, NombreDeRutina: rutinaSeleccionada.NombreDeRutina, Ejercicios: rutinaSeleccionada.Ejercicios, DuracionTotal: rutinaSeleccionada.DuracionTotal, PuntosCardioTotales: rutinaSeleccionada.PuntosCardioTotales, PuntosFlexibilidadTotales: rutinaSeleccionada.PuntosFlexibilidadTotales, PuntosFuerzaTotales: rutinaSeleccionada.PuntosFuerzaTotales}) // Add rutinas

			err = gocsv.MarshalFile(&rutinasCsv, rutinasFile) // Use this to save the CSV back to the file
			if err != nil {
				panic(err)
			}
			bluePrintf("Ejercicio eliminado. Nueva rutina %v creada\n", rutinaSeleccionada.NombreDeRutina)
			return
		case "volver":
			bluePrintf("Cambios guardados\n")
			return

		case "4":
			bluePrintf("Cambios guardados\n")
			return
		}
		return
	}
}

func seleccionarEjerciciosParaTiempoMaximo(categoria string, tiempoDisponible int) []Ejercicio {
	ejercicios, ok := categorias[categoria]
	if !ok {
		fmt.Println("Categoría no válida.")
		return nil
	}

	// Ordenar los ejercicios por duración de menor a mayor
	sort.Slice(ejercicios, func(i, j int) bool {
		return ejercicios[i].Duracion < ejercicios[j].Duracion
	})

	var ejerciciosSeleccionados []Ejercicio
	var tiempoTotal int

	for _, ejercicio := range ejercicios {
		if tiempoTotal+ejercicio.Duracion <= tiempoDisponible {
			ejerciciosSeleccionados = append(ejerciciosSeleccionados, ejercicio)
			tiempoTotal += ejercicio.Duracion
		} else {
			break
		}
	}

	return ejerciciosSeleccionados
}

func agregarEjerciciosMaximosARutina(nombreDeRut string /*nombreDeRut int,*/, categoria string, tiempoDisponible int) {

	// rutinaSeleccionada := &rutinasL[nombreDeRut]

	var rutinaSeleccionada Rutina
	rutinaSeleccionada.NombreDeRutina = nombreDeRut

	ejerciciosSeleccionados := seleccionarEjerciciosParaTiempoMaximo(categoria, tiempoDisponible)
	if ejerciciosSeleccionados == nil {
		return
	}

	for _, ejercicio := range ejerciciosSeleccionados {
		agregarEjercicioARutina(&rutinaSeleccionada, ejercicio)
	}

	fmt.Printf("\nSe han agregado %d ejercicios a la rutina '%s' de la categoría '%s' con un tiempo total de %d segundos.\n",
		len(ejerciciosSeleccionados), rutinaSeleccionada.NombreDeRutina, categoria, rutinaSeleccionada.DuracionTotal)

	rutinasL = append(rutinasL, rutinaSeleccionada)
	//Start save in csv
	rutinasFile, err := os.OpenFile("rutinas.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer rutinasFile.Close()
	rutinasCsv := []*RutinaCsv{}
	if err := gocsv.UnmarshalFile(rutinasFile, &rutinasCsv); err != nil { // Load rutinas from file
		panic(err)
	}

	if _, err := rutinasFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}

	rutinasCreadas := fmt.Sprintf("%v", len(rutinasCsv)+1)

	rutinasCsv = append(rutinasCsv, &RutinaCsv{Id: rutinasCreadas, NombreDeRutina: rutinaSeleccionada.NombreDeRutina, Ejercicios: rutinaSeleccionada.Ejercicios, DuracionTotal: rutinaSeleccionada.DuracionTotal}) // Add rutinas

	err = gocsv.MarshalFile(&rutinasCsv, rutinasFile) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}
}

/*
	func seleccionarEjerciciosMinDuracion(rutina Rutina, categoria string, tiempoDisponible int, dificultad string, tipo string) []Ejercicio {
		ejercicios, ok := categorias[categoria]
		if !ok {
			fmt.Println("Categoría no válida.")
			return nil
		}

		// Filtrar ejercicios por tipo y dificultad
		var ejerciciosFiltrados []Ejercicio
		for _, ejercicio := range ejercicios {
			if ejercicio.Dificultad == dificultad && ejercicio.Tipo == tipo {
				ejerciciosFiltrados = append(ejerciciosFiltrados, ejercicio)
			}
		}

		// Ordenar los ejercicios por duración de menor a mayor
		sort.Slice(ejerciciosFiltrados, func(i, j int) bool {
			return ejerciciosFiltrados[i].Duracion < ejerciciosFiltrados[j].Duracion
		})

		var ejerciciosSeleccionados []Ejercicio
		var tiempoTotal int

		for _, ejercicio := range ejerciciosFiltrados {
			if tiempoTotal+ejercicio.Duracion <= tiempoDisponible {
				// Verificar si el ejercicio ya está en la rutina para evitar repeticiones
				if !ejercicioEnRutina(rutina, ejercicio) {
					ejerciciosSeleccionados = append(ejerciciosSeleccionados, ejercicio)
					tiempoTotal += ejercicio.Duracion
				}
			} else {
				break
			}
		}

		return ejerciciosSeleccionados
	}

// Función auxiliar para verificar si un ejercicio ya está en la rutina

	func ejercicioEnRutina(rutina Rutina, ejercicio Ejercicio) bool {
		for _, ej := range rutina.Ejercicios {
			if ej.Nombre == ejercicio.Nombre {
				return true
			}
		}
		return false
	}
*/
// Función para seleccionar los ejercicios más eficientes
func seleccionarEjerciciosEficientes(categoria string, caloriasDeseadas int) []Ejercicio {
	ejercicios, ok := categorias[categoria]
	if !ok {
		fmt.Println("Categoría no válida.")
		return []Ejercicio{}
	}

	// Ordenar los ejercicios por la mejor relación calorías/duración
	sort.Slice(ejercicios, func(i, j int) bool {
		ratioI := float64(ejercicios[i].Calorias) / float64(ejercicios[i].Duracion)
		ratioJ := float64(ejercicios[j].Calorias) / float64(ejercicios[j].Duracion)
		return ratioI > ratioJ
	})

	var ejerciciosSeleccionados []Ejercicio
	caloriasTotales := 0

	for _, ejercicio := range ejercicios {
		if ejercicio.Calorias == caloriasDeseadas {
			// Si un ejercicio por sí solo cumple el objetivo, seleccionarlo
			return []Ejercicio{ejercicio}
		}
		if caloriasTotales >= caloriasDeseadas {
			break
		}
		ejerciciosSeleccionados = append(ejerciciosSeleccionados, ejercicio)
		caloriasTotales += ejercicio.Calorias
	}

	if caloriasTotales < caloriasDeseadas {
		fmt.Println("No se encontraron suficientes ejercicios para alcanzar las calorías deseadas.")
	}

	return ejerciciosSeleccionados
}

// Función para agregar los ejercicios eficientes a una rutina
func agregarEjerciciosEficientesARutina(nombreDeRut string, categoria string, caloriasDeseadas int) {
	var rutinaSeleccionada Rutina
	rutinaSeleccionada.NombreDeRutina = nombreDeRut

	ejerciciosEficientes := seleccionarEjerciciosEficientes(categoria, caloriasDeseadas)
	if len(ejerciciosEficientes) == 0 {
		fmt.Println("No se pudieron encontrar ejercicios eficientes.")
		return
	}

	for _, ejercicio := range ejerciciosEficientes {
		agregarEjercicioARutina(&rutinaSeleccionada, ejercicio)
	}

	fmt.Printf("\nSe han agregado ejercicios a la rutina '%s' de la categoría '%s' para quemar al menos %d calorías.\n",
		rutinaSeleccionada.NombreDeRutina, categoria, caloriasDeseadas)

	rutinasL = append(rutinasL, rutinaSeleccionada)

	// Guardar la rutina en CSV
	rutinasFile, err := os.OpenFile("rutinas.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer rutinasFile.Close()
	rutinasCsv := []*RutinaCsv{}
	if err := gocsv.UnmarshalFile(rutinasFile, &rutinasCsv); err != nil { // Cargar rutinas del archivo
		panic(err)
	}

	if _, err := rutinasFile.Seek(0, 0); err != nil { // Ir al inicio del archivo
		panic(err)
	}

	rutinasCreadas := fmt.Sprintf("%v", len(rutinasCsv)+1)

	rutinasCsv = append(rutinasCsv, &RutinaCsv{
		Id:                        rutinasCreadas,
		NombreDeRutina:            rutinaSeleccionada.NombreDeRutina,
		Ejercicios:                rutinaSeleccionada.Ejercicios,
		DuracionTotal:             rutinaSeleccionada.DuracionTotal,
		PuntosCardioTotales:       rutinaSeleccionada.PuntosCardioTotales,
		PuntosFlexibilidadTotales: rutinaSeleccionada.PuntosFlexibilidadTotales,
		PuntosFuerzaTotales:       rutinaSeleccionada.PuntosFuerzaTotales,
	}) // Agregar rutinas

	err = gocsv.MarshalFile(&rutinasCsv, rutinasFile) // Guardar el CSV en el archivo
	if err != nil {
		panic(err)
	}
}

func (r *RutinaCsv) String() string {
	ejerciciosStr := make([]string, len(r.Ejercicios))
	for i, e := range r.Ejercicios {
		ejerciciosStr[i] = fmt.Sprintf("{Nombre: %s, Duracion: %d, Tipo: %s, Intensidad: %s, Calorias: %d, Descripcion: %s, PuntosCardio: %d, PuntosFlexibilidad: %d, PuntosFuerza: %d}",
			e.Nombre, e.Duracion, e.Tipo, e.Intensidad, e.Calorias, e.Descripcion, e.PuntosCardio, e.PuntosFlexibilidad, e.PuntosFuerza)
	}
	return fmt.Sprintf("Id: %s, NombreDeRutina: %s, Ejercicios: [%s], DuracionTotal: %d, PuntosCardioTotales: %d, PuntosFlexibilidadTotales: %d, PuntosFuerzaTotales: %d",
		r.Id, r.NombreDeRutina, strings.Join(ejerciciosStr, ", "), r.DuracionTotal, r.PuntosCardioTotales, r.PuntosFlexibilidadTotales, r.PuntosFuerzaTotales)
}
