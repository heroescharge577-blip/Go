package main

import (
	"fmt"
	"math/rand"
	"net/http" // La librería que convierte Go en un servidor web
	"time"
)

// Las tres opciones posibles del juego
var opciones = []string{"piedra", "papel", "tijera"}

// elegirComputador elige una opción al azar para la máquina
func elegirComputador() string {
	rand.Seed(time.Now().UnixNano())
	return opciones[rand.Intn(3)]
}

// decidirGanador compara las dos elecciones y devuelve el resultado como texto
func decidirGanador(jugador, computador string) string {
	if jugador == computador {
		return "¡Empate!"
	}

	// Mapa de quién le gana a quién.
	// La clave le gana al valor.
	gana := map[string]string{
		"piedra": "tijera",
		"tijera": "papel",
		"papel":  "piedra",
	}

	// Si lo que vence al jugador es el computador, el computador gana
	if gana[jugador] == computador {
		return "¡Ganaste! 🏆"
	}

	return "¡Perdiste! 💀"
}

// manejarInicio responde cuando el jugador visita localhost:8080
// Su único trabajo es servir el archivo HTML
func manejarInicio(w http.ResponseWriter, r *http.Request) {
	// http.ServeFile busca el archivo y lo envía al navegador
	http.ServeFile(w, r, "static/index.html")
}

// manejarJugar responde cuando el jugador hace clic en un botón
// Aquí vive toda la lógica del juego
func manejarJugar(w http.ResponseWriter, r *http.Request) {
	jugador := r.FormValue("jugador")

	if jugador != "piedra" && jugador != "papel" && jugador != "tijera" {
		fmt.Fprintf(w, "Opción inválida")
		return
	}

	computador := elegirComputador()
	resultado := decidirGanador(jugador, computador)

	// Emojis para mostrar visualmente la elección
	emojis := map[string]string{
		"piedra":  "🪨",
		"papel":   "📄",
		"tijera":  "✂️",
	}

	// Go devuelve una página HTML completa con CSS incluido
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <title>Resultado</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: sans-serif;
            background: #1a1a2e;
            color: #eee;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            gap: 1.5rem;
        }
        .tarjeta {
            background: #16213e;
            border: 2px solid #444;
            border-radius: 20px;
            padding: 2.5rem 3rem;
            text-align: center;
            display: flex;
            flex-direction: column;
            gap: 1.2rem;
        }
        .versus {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 2rem;
            font-size: 4rem;
        }
        .versus span { font-size: 1rem; color: #aaa; display: block; margin-top: 6px; }
        .vs { font-size: 1.2rem; color: #e94560; font-weight: bold; }
        .resultado {
            font-size: 1.8rem;
            font-weight: bold;
            color: #e94560;
        }
        .boton {
            margin-top: 0.5rem;
            padding: 0.8rem 2rem;
            font-size: 1rem;
            background: #e94560;
            color: white;
            border: none;
            border-radius: 12px;
            cursor: pointer;
            text-decoration: none;
            display: inline-block;
        }
        .boton:hover { background: #c73652; }
    </style>
</head>
<body>
    <div class="tarjeta">
        <div class="versus">
            <div>%s<span>Tú</span></div>
            <div class="vs">VS</div>
            <div>%s<span>Computador</span></div>
        </div>
        <div class="resultado">%s</div>
        <a class="boton" href="/">Jugar de nuevo</a>
    </div>
</body>
</html>`, emojis[jugador], emojis[computador], resultado)
}

func main() {
	// Le decimos a Go: "cuando alguien visite /, ejecuta manejarInicio"
	http.HandleFunc("/", manejarInicio)

	// Y cuando visiten /jugar, ejecuta manejarJugar
	http.HandleFunc("/jugar", manejarJugar)

	fmt.Println("Servidor corriendo en http://localhost:8080")

	// ListenAndServe arranca el servidor en el puerto 8080.
	// El nil significa "usa la configuración por defecto".
	// Esta línea bloquea el programa – Go se queda escuchando peticiones.
	http.ListenAndServe(":8080", nil)
}