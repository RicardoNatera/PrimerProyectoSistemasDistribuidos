package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/sync/semaphore"
)

//Funciones referentes a informacion de los diferentes servicios
func RunCmd(cmd string) (result []byte, err error) {
	cmdPip := strings.Split(cmd, "|")
	if len(cmdPip) < 2 {
		return run(cmd).Output()
	}

	return runPipe(cmdPip)
}

func run(cmd string) *exec.Cmd {
	if strings.Contains(cmd, "findstr") || !strings.Contains(cmd, "find") {
		cmd = strings.Replace(cmd, `"`, "", -1)
		cmd = strings.TrimSpace(cmd)
	}
	cmdList := strings.Split(cmd, " ")

	return exec.Command(cmdList[0], cmdList[1:]...)
}

func runPipe(pip []string) (result []byte, err error) {
	var cmds []*exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmds = append(cmds, run(winPre+pip[0]))
	case "linux":
		cmds = append(cmds, run(LinuxPre+pip[0]))
	default:
		return nil, fmt.Errorf("Not supported by the system. %v", runtime.GOOS)
	}

	for i := 1; i < len(pip); i++ {
		cmds = append(cmds, run(pip[i]))
		cmds[i].Stdin, err = cmds[i-1].StdoutPipe()
		if err != nil {
			return nil, err
		}
	}

	end := len(cmds) - 1
	// cmds[end].Stdout = os.Stdout
	stdout, _ := cmds[end].StdoutPipe()

	for i := end; i > 0; i-- {
		cmds[i].Start()
	}
	cmds[0].Run()

	buf := make([]byte, 102400000)
	n, _ := stdout.Read(buf)

	err = cmds[end].Wait()

	return buf[:n], err
}

const (
	winPre   = "cmd /c "
	LinuxPre = "bash /c "
)

func info() string {
	var cadena string
	cmd := `netstat -n`
	result, err := RunCmd(cmd)
	if err != nil {
		log.Fatal(err, "hola")
	}
	var con_ESTABLISHED, con_CLOSE_WAIT, con_TIME_WAIT, con_LISTEN int
	b := string(result)
	cadena = cadena + "Informacion Completa" + "\n"
	cadena = cadena + b
	cadena = cadena + "Informacion Resumida" + "\n"
	con_ESTABLISHED = strings.Count(b, "ESTABLISHED")
	con_CLOSE_WAIT = strings.Count(b, "CLOSE_WAIT")
	con_TIME_WAIT = strings.Count(b, "TIME_WAIT")
	con_LISTEN = strings.Count(b, "LISTEN")
	cadena = cadena + "established : " + strconv.Itoa(con_ESTABLISHED) + "\n"
	cadena = cadena + "close_wait  : " + strconv.Itoa(con_CLOSE_WAIT) + "\n"
	cadena = cadena + "time_wait   : " + strconv.Itoa(con_TIME_WAIT) + "\n"
	cadena = cadena + "Listen      : " + strconv.Itoa(con_LISTEN) + "\n"

	cadena = cadena + "Estadisticas de todos los servicios de red" + "\n"

	cmd = `netstat -s`
	result, err = RunCmd(cmd)
	if err != nil {
		log.Fatal(err)
	}

	b = string(result)
	cadena = cadena + b + "\n"

	cadena = cadena + "Informacion De quien usa el puerto 8080" + "\n"

	cmd = `netstat -anlp |grep 8080 | grep LISTEN`
	result, err = RunCmd(cmd)
	if err != nil {
		cadena = cadena + "El puerto 8080 esta activo actualmente " + "\n"
	}

	b = string(result)
	cadena = cadena + b + "\n"
	cadena = cadena + "Informacion De quien usa el puerto 2002" + "\n"

	cmd = `netstat -anlp |grep 2002 | grep LISTEN`
	result, err = RunCmd(cmd)
	if err != nil {
		cadena = cadena + "El puerto 2002 esta activo actualmente " + "\n"

	}

	b = string(result)
	cadena = cadena + b + "\n"
	cadena = cadena + "Informacion De quien usa el puerto 2020" + "\n"

	cmd = `netstat -anlp |grep 2020 | grep LISTEN`
	result, err = RunCmd(cmd)
	if err != nil {
		cadena = cadena + "El puerto 2020 esta activo actualmente " + "\n"

	}

	b = string(result)
	cadena = cadena + b + "\n"
	return cadena

}

//-------------------------------------------------------------

//-----------------------------------

// Para las llamadas RPC
type Args struct {
	Action float64
	Valor  float64
}
type API int

func (a *API) Information(valor int, cadena *string) error {
	cmd := `netstat -n`
	result, err := RunCmd(cmd)
	if err != nil {
		log.Fatal(err)
	}
	var con_ESTABLISHED, con_CLOSE_WAIT, con_TIME_WAIT, con_LISTEN int
	b := string(result)
	*cadena = *cadena + "Informacion Completa" + "\n"
	*cadena = *cadena + b
	*cadena = *cadena + "Informacion Resumida" + "\n"
	con_ESTABLISHED = strings.Count(b, "ESTABLISHED")
	con_CLOSE_WAIT = strings.Count(b, "CLOSE_WAIT")
	con_TIME_WAIT = strings.Count(b, "TIME_WAIT")
	con_LISTEN = strings.Count(b, "LISTEN")
	*cadena = *cadena + "established : " + strconv.Itoa(con_ESTABLISHED) + "\n"
	*cadena = *cadena + "close_wait  : " + strconv.Itoa(con_CLOSE_WAIT) + "\n"
	*cadena = *cadena + "time_wait   : " + strconv.Itoa(con_TIME_WAIT) + "\n"
	*cadena = *cadena + "Listen      : " + strconv.Itoa(con_LISTEN) + "\n"

	*cadena = *cadena + "Estadisticas de todos los servicios de red" + "\n"

	cmd = `netstat -s`
	result, err = RunCmd(cmd)
	if err != nil {
		log.Fatal(err)
	}

	b = string(result)
	*cadena = *cadena + b + "\n"

	*cadena = *cadena + "Informacion De quien usa el puerto 8080" + "\n"

	cmd = `netstat -anlp |grep 8080 | grep LISTEN`
	result, err = RunCmd(cmd)
	if err != nil {
		*cadena = *cadena + "El puerto 8080 esta activo actualmente " + "\n"
	}

	b = string(result)
	*cadena = *cadena + b + "\n"
	*cadena = *cadena + "Informacion De quien usa el puerto 2002" + "\n"

	cmd = `netstat -anlp |grep 2002 | grep LISTEN`
	result, err = RunCmd(cmd)
	if err != nil {
		*cadena = *cadena + "El puerto 2002 esta activo actualmente " + "\n"

	}

	b = string(result)
	*cadena = *cadena + b + "\n"
	*cadena = *cadena + "Informacion De quien usa el puerto 2020" + "\n"

	cmd = `netstat -anlp |grep 2020 | grep LISTEN`
	result, err = RunCmd(cmd)
	if err != nil {
		*cadena = *cadena + "El puerto 2020 esta activo actualmente " + "\n"

	}

	b = string(result)
	*cadena = *cadena + b + "\n"

	return nil
}

func (a *API) Operation(args Args, reply *int64) error {

	var y int = int(args.Valor)
	t := time.Now().Format(time.RFC3339)
	mensaje := Mensaje{y, "cadena", t}
	if args.Action < 0 { //decrementar -1
		mensaje.Action = "dec"

	}
	if args.Action > 1 && args.Action < 3 { //incrementar 2
		mensaje.Action = "inc"
	}
	if args.Action > 5 { //reset 8
		mensaje.Action = "res"
	}

	if err := semContador.Acquire(ctxContador, 1); err != nil {
		fmt.Println(err)
	}
	wg.Add(1)
	go func() {

		queue.AddQueue(mensaje)
		semContador.Release(1)
		wg.Done()
	}()
	wg.Wait()

	wg.Add(1)
	go func() {
		mensaje, err := queue.GetQueue()
		if err != nil {
			fmt.Println(err)
		} else {

			operationCounter(mensaje.Action, int64(mensaje.Cant))
			*reply = value

		}
		wg.Done()
	}()
	wg.Wait()
	return nil

}

//------------------------------------

var value int64                                    //Contador
var semContador = semaphore.NewWeighted(int64(20)) //Semaforo del contador
var ctxContador = context.Background()             // Contexto de la aplicacion
var wg sync.WaitGroup                              //Para el manejo de grupo de hilos
var channelColaDeMensajes chan Mensaje             //canal de comunicacion entre goRoutines referentes a la cola de mensajes
var channelSemaforo chan Mensaje                   //canal de comunicacion entre goRoutines referentes al semaforo
//canal de intercambio cuando existe algun cambio en general, en la cola de mensajes o en el semaforo
//este canal fue pensado para mantener un control entre las respuestas de los servidores para que fueran lo mas rapida posible
var response chan string

var queue = &Queue{ //cola de mensajes
	maxSize: 1000,
	front:   -1,
	rear:    -1,
}

type Mensaje struct { //Estructura de los mensajes
	Cant   int
	Action string
	ID     string
}
type Queue struct { //Estructura de la cola de mensajes
	maxSize int
	array   [1000]Mensaje // Cola de simulación de matriz
	front   int           //  apuntar a la cabecera de la cola
	rear    int           //  apuntar al final de la cola
}

// Agregar datos a la cola
func (this *Queue) AddQueue(msj Mensaje) (err error) {

	if this.rear == this.maxSize-1 {
		return errors.New("queue full")
	}
	this.rear++ // parte trasera se mueve hacia atrás
	this.array[this.rear] = msj

	return
}

// Eliminar datos de la cola
func (this *Queue) GetQueue() (val Mensaje, err error) {
	// Primero determina si la cola está vacía
	if this.rear == this.front {
		return Mensaje{0, "0", "0"}, errors.New("Cola de Mensajes Vacia")
	}
	this.front++
	val = this.array[this.front]
	return val, err
}

func process(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}
		split := strings.Split(temp, ".")
		cadena := split[0] //inc.25=inc.
		aux := split[1]    //inc.25=25
		val, err := strconv.Atoi(aux)
		t := time.Now().Format(time.RFC3339)
		mensaje := Mensaje{val, cadena, t}

		if err := semContador.Acquire(ctxContador, 1); err != nil {
			fmt.Println(err)
		}
		wg.Add(1)
		go func() {
			queue.AddQueue(mensaje)
			semContador.Release(1)
			wg.Done()
		}()
		wg.Wait()
		result := "value"

		wg.Add(1)
		go func() {
			mensaje, err = queue.GetQueue()
			if err != nil {
				fmt.Println(err)
			} else {

				operationCounter(mensaje.Action, int64(mensaje.Cant))
				result = strconv.FormatInt(value, 10) + "\n"

			}
			wg.Done()
		}()
		wg.Wait()

		c.Write([]byte(string(result)))
	}
	c.Close()

}

//----Inicializacion de Servidor TCP (Hilos)
func TCP_Server_Hilos() {
	PORT := ":" + "2020"
	response = make(chan string)

	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())
	fmt.Println("inicializando servidor tcp en puerto 2020")
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("iniciando nueva conexion con el servidor TCP por hilos")
		go process(c)
	}
}

//-----Inicializacion de Servidor UDP
func UDP_Server() {
	PORT := ":" + "2002"

	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()
	buffer := make([]byte, 1024)
	rand.Seed(time.Now().Unix())
	fmt.Println("inicializando servidor udp en puerto 2002")
	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		temp := strings.TrimSpace(string(buffer[0:n]))
		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}
		split := strings.Split(temp, ".")
		cadena := split[0] //inc.25=inc.
		aux := split[1]    //inc.25=25
		val, err := strconv.Atoi(aux)
		t := time.Now().Format(time.RFC3339)
		mensaje := Mensaje{val, cadena, t}

		if err := semContador.Acquire(ctxContador, 1); err != nil {
			fmt.Println(err)
		}
		wg.Add(1)
		go func() {

			queue.AddQueue(mensaje)
			semContador.Release(1)
			wg.Done()
		}()
		wg.Wait()
		result := "value"

		wg.Add(1)
		go func() {
			mensaje, err = queue.GetQueue()
			if err != nil {
				fmt.Println(err)
			} else {

				operationCounter(mensaje.Action, int64(mensaje.Cant))
				result = strconv.FormatInt(value, 10) + "\n"

			}
			wg.Done()
		}()
		wg.Wait()

		data := []byte(string(result))

		_, err = connection.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

//-----Inicializacion de Servicio RPC
func RPC_Service() {
	api := new(API)
	err := rpc.Register(api)
	if err != nil {
		log.Fatal("error registering API", err)
	}
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", "127.0.0.1:8080")

	if err != nil {
		log.Fatal("Listener error", err)
	}
	log.Printf("serving rpc on port %d", 8080)
	http.Serve(listener, nil)

	if err != nil {
		log.Fatal("error serving: ", err)
	}
}

func operationCounter(action string, cant int64) {
	if value == 0 && action == "dec" {
		return
	}
	var act int64
	act = 0

	switch action {
	case "inc":
		act = cant
	case "dec":
		act = -cant
	default:
		act = -value
	}

	value = value + act
	return
}

func counter(w http.ResponseWriter, req *http.Request) {

	consulta, noerr := req.URL.Query()["counter"]

	if noerr {

		split := strings.Split(consulta[0], ".")
		cadena := split[0] //inc.25=inc.
		aux := split[1]    //inc.25=25
		val, _ := strconv.Atoi(aux)
		t := time.Now().Format(time.RFC3339)
		mensaje := Mensaje{val, cadena, t}

		if err := semContador.Acquire(ctxContador, 1); err != nil {
			fmt.Println(err)
		}
		wg.Add(1)
		go func() {

			queue.AddQueue(mensaje)
			semContador.Release(1)
			wg.Done()
		}()
		wg.Wait()

		wg.Add(1)
		go func() {
			mensaje, err := queue.GetQueue()
			if err != nil {
				fmt.Println(err)
			} else {

				operationCounter(mensaje.Action, int64(mensaje.Cant))

			}
			wg.Done()
		}()
		wg.Wait()
		io.WriteString(w, "Valor del Contador, "+strconv.Itoa(int(value)))
	} else {
		io.WriteString(w, "consulta mal formulada")
	}
}
func information(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, info())
}
func main() {
	channelColaDeMensajes = make(chan Mensaje) //Canal de comunicacion de la cola
	channelSemaforo = make(chan Mensaje)       //canal de comunicacion del semaforo
	queue = &Queue{
		maxSize: 1000,
		front:   -1,
		rear:    -1,
	}

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println((" Error Cargando .env "))
	}
	//Se lee la variable del sistema, el valor del contador
	val, err := strconv.ParseInt(os.Getenv("COUNTER"), 10, 64)
	value = val
	if err != nil {
		fmt.Println((" Error en la lectura de counter "))
	}
	//-------------------------
	//Iniciamos TCP por hilos
	go TCP_Server_Hilos()

	//Iniciamos UDP
	go UDP_Server()
	//--------------------
	//Iniciamos Servicio RPC
	go RPC_Service()
	//----------------------------------------------------------------------------------
	http.HandleFunc("/cambiar", counter)
	http.HandleFunc("/info", information)
	//Configuraciones Completas
	fmt.Println((" Configuraciones Completas "))
	fmt.Println((" Servidor Local en el puerto 1234 "))
	http.ListenAndServe(":1234", nil)

	//------------------------------
	//Preparando para cambiar el valor del contador antes de que el servidor se desconecte
	env, err := godotenv.Unmarshal(fmt.Sprintf("%s%d", "COUNTER=", value))
	if err != nil {
		fmt.Println((" Error en la preparaciòn "))
	}
	//Guardando nuevo valor para contador
	erro := godotenv.Write(env, "./.env")
	if erro != nil {
		fmt.Println((" Error en la escritura "))
	}
}
