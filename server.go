package main

import (
	"strconv"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"strings"
)

var OMDBAPIKEY = LoadAPIKey("apikey.txt")
var videosFolder = LoadAPIKey("path.txt")
var vidsrcUrl  = "https://vsembed.ru/embed/"

type PageData struct {
	Videos []string
	IframeURL string
}

type OmdbResponse struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	ImdbID   string `json:"imdbID"`
	Type     string `json:"Type"`
	Response string `json:"Response"`
	Error    string `json:"Error"`
}

func LoadAPIKey(filename string) (string) {
	data, err := os.ReadFile(filename)
	if err != nil {
//		return "", err
	}

	a := strings.TrimSpace(string(data))
	return a
}

// searchOMDb now takes an extra "mediaType" param.
// • mediaType == "movie", "series", or "episode" → adds &type=...
// • anything else / empty → no filter (OMDb chooses best match).
func searchOMDb(title, apiKey, mediaType string) (*OmdbResponse, error) {
	baseURL := "http://www.omdbapi.com/"
	params := url.Values{}
	params.Set("t", title)
	params.Set("apikey", apiKey)

	switch mediaType {
	case "movie", "series", "episode":
		params.Set("type", mediaType)
	}

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("http error: %v", err)
		//return
	}
	defer resp.Body.Close()

	var result OmdbResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}

	/*if result.Response != "True" {
	    return nil, fmt.Errorf("OMDb error: %s", result.Error)
	}*/

	return &result, nil
}

func resetMovie() {
	//fmt.Println("nigger")
        if err := exec.Command("bash", "resetMovie.sh").Run(); err != nil {
                log.Printf("failed to kill run hitPlay.sh: %v", err)
        }
}

func hitPlay() {
	if err := exec.Command("bash", "hitPlay.sh").Run(); err != nil {
		log.Printf("failed to kill run hitPlay.sh: %v", err)
	}
}

func hideMouse() {
	if err := exec.Command("xdotool", "mousemove", "940", "20").Run(); err != nil {
		log.Printf("failed to kill run command: %v", err)
	}
}

var name bool

func fullScreen() {
	if name {
		// Run fullscreen2.sh and set name to false
		if err := exec.Command("bash", "fullscreen2.sh").Run(); err != nil {
			log.Printf("Error running fullscreen2.sh: %v", err)
		}
		name = false
	} else {
		// Run fullscreen.sh and set name to true
		if err := exec.Command("bash", "fullscreen.sh").Run(); err != nil {
			log.Printf("Error running fullscreen.sh: %v", err)
		}
		name = true
	}
}

// kill the fox
func executeFirefox() {
	if err := exec.Command("pkill", "firefox").Run(); err != nil {
		log.Printf("failed to kill firefox: %v", err)
	}
}

//type PageData struct {
//        IframeURL string
//}
//the url for the iframe in the viewer for roryflix
//the viewer that is displayed on the tv
var globalViewerUrl string

func viewer(w http.ResponseWriter, r *http.Request) {
        data := PageData{
                IframeURL: globalViewerUrl, // your movie लिंक
        }

        tmpl := template.Must(template.New("page").Parse(pageHTML))
        tmpl.Execute(w, data)
}

var pageHTML = `
<!DOCTYPE html>
<html>
<head>
     <style>
html, body {
            margin: 0;
            padding: 0;
            height: 100%;
            background: black;
            overflow: hidden; /* removes scrollbars */
        }
        iframe {
            width: 100vw;   /* full viewport width */
            height: 99vh;  /* full viewport height */
            border: none;
        }
    </style>
    </head>
<body>
<iframe src="{{.IframeURL}}" width="100%" height="100%"></iframe>
</body>
</html>
`


//the 3 execute funcs here are used to launch shows, we used to kill firefox and open a new window with our intended url, instead we shall change the src of an iframe and maybe have  abutton to open and close our window on the server from the remote.
// use firefox to launch movie
func executeMovie(id string) {
	name = false
	//executeFirefox()
	time.Sleep(2 * time.Second)
	var url string
	if serverMode == "vidsrc" {
        	url = fmt.Sprintf("https://vsembed.ru/embed/%s", id)
	} else {
                url = fmt.Sprintf("https://moviesapi.to/movie/%s", id)
        }
	
	time.Sleep(2 * time.Second)
	
	globalViewerUrl = url
}

func executeImdbShow(id, season, episode, mediaType string) {
        name = false
        
	var url string
        if mediaType == "movie" {
        
		if serverMode == "vidsrc" {
                	url = fmt.Sprintf("%s%s", vidsrcUrl, id)
        	} else {
                	url = fmt.Sprintf("https://moviesapi.to/movie/%s", id)
        	}

	} else if mediaType == "series" {
                
		if serverMode == "vidsrc" {
                	url = fmt.Sprintf("%s%s/%s-%s", vidsrcUrl, id, season, episode)
        	} else {
        	        url = fmt.Sprintf("https://moviesapi.to/tv/%s-%s-%s", id, season, episode)
        	}

	}
	
	globalViewerUrl = url

}


// use firefox to launch tv show
func executeShow(id, season, episode string) {
	name = false
	//kill firefox
	//executeFirefox()
	//time.Sleep(2 * time.Second)
	var url string
	if serverMode == "vidsrc" {
		url = fmt.Sprintf("https://vsembed.ru/embed/%s/%s-%s", id, season, episode)
	} else {
		url = fmt.Sprintf("https://moviesapi.to/tv/%s-%s-%s", id, season, episode)
	}

	globalViewerUrl = url
}

// I hate that they're two different functions but fuck it we ball
// show search, sel should set to "movie" / "series" / "episode" as needed
func TvSearch(name string, sel string, season string, episode string) {
	apiKey := OMDBAPIKEY

	title := name
	mediaType := sel

	result, err := searchOMDb(title, apiKey, mediaType)
	if err != nil {
		log.Fatal(err)
	}

	//launch show in firefox
	executeShow(result.ImdbID, season, episode)
}

// Movie search func, sel should set to "movie" as needed
func movieSearch(name, sel string) {
	apiKey := OMDBAPIKEY

	title := name
	mediaType := sel

	result, err := searchOMDb(title, apiKey, mediaType)
	if err != nil {
		log.Fatal(err)
	}

	//launch show in firefox
	executeMovie(result.ImdbID)
}

func returnIp() net.IP {
	// Connect to an external address (Google DNS); no data is sent
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Extract the local address (what your system would use)
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	//fmt.Println("Primary IP address:", localAddr.IP)
	return localAddr.IP
}
//x and y are mouse cords
var x = 950
var y = 590
var isOpen = false

var serverMode string
var remoteControlIp = ""
var timeStamp int 
var pixelDistance int
// Tv Remote function
func remote(w http.ResponseWriter, r *http.Request) {
	addr := r.RemoteAddr
        ip := strings.Split(addr, ":")[0]
	now := time.Now()
	
	speedStr := r.FormValue("speed")
	
	if speedStr == "" && pixelDistance <= 10 {
		pixelDistance = 10
	}

	if speedStr != "" {
	    if n, err := strconv.Atoi(speedStr); err == nil {
	        pixelDistance = n
	    }
	}
	
	fmt.Println(pixelDistance)
	
	if remoteControlIp == "" {
                        remoteControlIp = ip
                        timeStamp = now.Hour()
                        fmt.Println("Ip address and time saved")
                        fmt.Println("user:", remoteControlIp, " controls the remote until ", timeStamp + 1, ":00")
                }
	if r.Method == http.MethodPost && ip == remoteControlIp {
		action := r.URL.Query().Get("action")
		

		switch action {
		
		case "pause":
			hitPlay()
			//fmt.Println("hit play condition hit")
			w.WriteHeader(http.StatusOK)
			return

		case "hideMouse":
			hideMouse()
			//fmt.Println("hideMouse triggered")
			w.WriteHeader(http.StatusOK)
			return

		case "fullScreen":
			fullScreen()
			//fmt.Println("fullScreen triggered")
			w.WriteHeader(http.StatusOK)
			return
		
		case "resetMovie":
                        resetMovie()
                        //fmt.Println("fullScreen triggered")
                        w.WriteHeader(http.StatusOK)
                        return
		case "imdbIdSearch":
			imdbIdSearch := r.FormValue("imdbIdSearch")
			imdbseason := r.FormValue("season")
			imdbepisode := r.FormValue("episode")
			imdbmediaType := r.FormValue("mediaType")
			//fmt.Println(imdbIdSearch)
			executeImdbShow(imdbIdSearch, imdbseason, imdbepisode, imdbmediaType)
			w.WriteHeader(http.StatusOK)
			return
		case "openViewer":
			if isOpen != true {
				isOpen = true
				if err := exec.Command("firefox", "http://localhost:8080/viewer").Run(); err != nil {
              				log.Printf("failed to launch firefox: %v", err)
					isOpen = false
        			}
			} else {
				executeFirefox()
				isOpen = false
			}
			w.WriteHeader(http.StatusOK)
			return
		case "killFirefox":
                        executeFirefox()
			w.WriteHeader(http.StatusOK)
                        return
		case "refresh":
			if err := exec.Command("bash", "refresh.sh").Run(); err != nil {
                		log.Printf("failed to kill run command: %v", err)
		        }
			w.WriteHeader(http.StatusOK)
                        return
		case "lowerVolume":
                        if err := exec.Command("amixer", "-D", "pulse", "sset", "Master", "5%-").Run(); err != nil {
                                log.Printf("failed to kill run command: %v", err)
                        }
                        w.WriteHeader(http.StatusOK)
                        return
		case "raiseVolume":
                        if err := exec.Command("amixer", "-D", "pulse", "sset", "Master", "5%+").Run(); err != nil {
                                log.Printf("failed to kill run command: %v", err)
                        }
			w.WriteHeader(http.StatusOK)
                        return
		
		case "click":
			if err := exec.Command("xdotool", "click", "1").Run(); err != nil {
                                log.Printf("failed to kill run command: %v", err)
                        }
			w.WriteHeader(http.StatusOK)
                        return

		case "up":
			y = y - pixelDistance
			stringX := fmt.Sprintf("%d", x)
			stringY := fmt.Sprintf("%d", y)
			//fmt.Println(stringX)
			if err := exec.Command("xdotool", "mousemove", stringX, stringY).Run(); err != nil {
                                log.Printf("failed to kill run command: %v", err)
                        }
                        w.WriteHeader(http.StatusOK)
                        return

		case "down":
			y = y + pixelDistance
                        stringX := fmt.Sprintf("%d", x)
                        stringY := fmt.Sprintf("%d", y)
                        //fmt.Println(stringX)
                        if err := exec.Command("xdotool", "mousemove", stringX, stringY).Run(); err != nil {
                                log.Printf("failed to kill run command: %v", err)
                        }
                        w.WriteHeader(http.StatusOK)
                        return

		case "left":
			x = x - pixelDistance
                        stringX := fmt.Sprintf("%d", x)
			stringY := fmt.Sprintf("%d", y)
                        
			if err := exec.Command("xdotool", "mousemove", stringX, stringY).Run(); err != nil {
                                log.Printf("failed to kill run command: %v", err)
                        }
                        w.WriteHeader(http.StatusOK)
                        return

		case "right":
			x = x + pixelDistance
                        stringX := fmt.Sprintf("%d", x)
                        stringY := fmt.Sprintf("%d", y)

                        if err := exec.Command("xdotool", "mousemove", stringX, stringY).Run(); err != nil {
                                log.Printf("failed to kill run command: %v", err)
                        }
                        w.WriteHeader(http.StatusOK)
                        return
		case "upleft":
			y = y - 10
			x = x - 10
                        stringX := fmt.Sprintf("%d", x)
                        stringY := fmt.Sprintf("%d", y)
                        //fmt.Println(stringX)
                        if err := exec.Command("xdotool", "mousemove", stringX, stringY).Run(); err != nil {
                                log.Printf("failed to kill run command: %v", err)
                        }
                        w.WriteHeader(http.StatusOK)
                        return

		case "upright":
			y = y - 10
                        x = x + 10
                        stringX := fmt.Sprintf("%d", x)
                        stringY := fmt.Sprintf("%d", y)
                        //fmt.Println(stringX)
                        if err := exec.Command("xdotool", "mousemove", stringX, stringY).Run(); err != nil {
                                log.Printf("failed to kill run command: %v", err)
                        }
                        w.WriteHeader(http.StatusOK)
                        return

		case "downleft":
			y = y + 10
                        x = x - 10
                        stringX := fmt.Sprintf("%d", x)
                        stringY := fmt.Sprintf("%d", y)
                        //fmt.Println(stringX)
                        if err := exec.Command("xdotool", "mousemove", stringX, stringY).Run(); err != nil {
                                log.Printf("failed to kill run command: %v", err)
                        }
                        w.WriteHeader(http.StatusOK)
                        return

		case "downright":
			y = y + 10
                        x = x + 10
                        stringX := fmt.Sprintf("%d", x)
                        stringY := fmt.Sprintf("%d", y)
                        //fmt.Println(stringX)
                        if err := exec.Command("xdotool", "mousemove", stringX, stringY).Run(); err != nil {
                                log.Printf("failed to kill run command: %v", err)
                        }
                        w.WriteHeader(http.StatusOK)
                        return

		}
		
		// If no recognized action, handle the form POST
		showName := r.FormValue("title")
		mediaType := r.FormValue("mediaType")
		season := r.FormValue("season")
		episode := r.FormValue("episode")
                serverType := r.FormValue("serverType")
	
		//serverMode is a global var, which should be a string containing which mode to use.
		serverMode = serverType
		//fmt.Println(serverType)
		//security mechanism 
		//takes time stamp and ip address and stores them as global vars
		//if the current hour is greater than or less than the time stamped hour
		//the remoteControlIp is reset to nothing, which should restart the security mechanism.
		//timeStamp is now an int which holds the time in 24 hours. 12 Am being 24 10 being 22
		//etc
		if timeStamp < now.Hour() ||timeStamp > now.Hour() {
			remoteControlIp = ""
		}
		if remoteControlIp == "" {
			remoteControlIp = ip
			timeStamp = now.Hour()
			fmt.Println("Ip address and time saved")
			fmt.Println("user:", remoteControlIp, " controls the remote until ", timeStamp + 1, ":00")
		}
		fmt.Println("Data:", mediaType, showName, "Season:", season, "Episode:", episode)

		if mediaType == "movie" && ip == remoteControlIp {
			/*if ip != remoteControlIp {
				return
			}*/
			//fmt.Println(timeStamp)
			movieSearch(showName, mediaType)
		} else if mediaType == "series" && ip == remoteControlIp {
			/*if ip != remoteControlIp {
                                return
                        }*/
			TvSearch(showName, mediaType, season, episode)
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	// Non-POST: serve HTML
	http.ServeFile(w, r, "roryflixWebPages/remote3.html")
}

// Search function for main pages
func handleIMDbSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	showName := r.FormValue("show")
	showType := r.FormValue("type") // <-- Get "movie", "series", or "episode"

	if showName == "" {
		http.Error(w, "Missing show name", http.StatusBadRequest)
		return
	}

	imdbID, err := fetchImdbID(showName, showType, "d510153b") // pass type now
	if err != nil {
		log.Println("Error fetching IMDb ID:", err)
		http.Error(w, "IMDb search failed", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, vidsrcUrl+imdbID, http.StatusSeeOther)
}

func fetchImdbID(showName string, showType string, apiKey string) (string, error) {
	endpoint := "http://www.omdbapi.com/?apikey=" + apiKey +
		"&t=" + url.QueryEscape(showName) +
		"&type=" + url.QueryEscape(showType)

	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result OmdbResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Response != "True" {
		return "", fmt.Errorf("OMDb error: %s", result.Error)
	}

	return result.ImdbID, nil
}

func Figlet(text string, font string) (string) {
	// Run `figlet` with the specified text and font
	cmd := exec.Command("figlet", "-f", font, text)
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return string(output)
}

func main() {

	http.HandleFunc("/imdbsearch", handleIMDbSearch)
	http.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir("./videos"))))
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/Tv", serveTv)
	//http.HandleFunc("/Movies", serveMovies)
	http.HandleFunc("/Playlist", servePlaylist)
	http.HandleFunc("/Videos", serveVideos)
	http.HandleFunc("/remote", remote)
	http.HandleFunc("/viewer", viewer)
	fmt.Println("Service running on ",returnIp(), ":8080")
	http.ListenAndServe(":8080", nil)

}

func serveIndex(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		
		http.ServeFile(w, r, "roryflixWebPages/index.html")
	}
}

func serveTv(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		http.ServeFile(w, r, "roryflixWebPages/roryflixTvImproved.html")
	}

}
/*
func serveMovies(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		http.ServeFile(w, r, "roryflixWebPages/roryflixMoviesImproved.html")
	}

}*/

func servePlaylist(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		http.ServeFile(w, r, "roryflixWebPages/roryflixPlaylistSaver.html")
	}

}

//this is the page for serving downloaded content

func serveVideos(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(videosFolder)
	if err != nil {
		http.Error(w, "Unable to read videos directory", 500)
		return
	}

	var videoFiles []string
	for _, file := range files {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if ext == ".mp4" || ext == ".webm" || ext == ".ogg" {
				videoFiles = append(videoFiles, file.Name())
			}
		}
	}

	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>RoryFlix Video List</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
/* Basic RoryFlix styles for navbar */
body {
  background-color: #141414;
  color: white;
  font-family: Arial, sans-serif;
  margin: 0;
  padding: 0;
}
header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 50px;
  background: rgba(0, 0, 0, 0.8);
  position: fixed;
  width: 100%;
  z-index: 1000;
}
header .logo {
  font-size: 24px;
  font-weight: bold;
  color: #E50914;
}
header nav {
  display: flex;
  gap: 20px;
}
header nav a {
  color: #fff;
  text-decoration: none;
  font-size: 16px;
  transition: color 0.3s;
}
header nav a:hover {
  color: #E50914;
}
.container {
  padding-top: 100px;
  max-width: 800px;
  margin: auto;
}
ul {
  list-style: none;
  padding: 0;
}
li {
  margin: 15px 0;
}
a.video-link {
  color: #1db954;
  text-decoration: none;
  font-size: 20px;
}
a.video-link:hover {
  text-decoration: underline;
}
</style>
</head>
<body>

<header>
  <div class="logo">RoryFlix</div>
  <nav>

    <a href="/">Home</a>
    <a href="/Tv">TV Shows & Movies</a>
    
    <a href="/Playlist">Playlist</a>
    <a href="/Videos">Videos</a>
    </nav>
</header>

<div class="container">
  <h1>Available Videos</h1>
  <ul>
    {{range .Videos}}
    <li><a class="video-link" href="/videos/{{.}}" target="_blank">{{.}}</a></li>
    {{end}}
  </ul>
</div>

</body>
</html>`

	t, err := template.New("page").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", 500)
		return
	}

	data := PageData{Videos: videoFiles}
	t.Execute(w, data)
}
