package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/skriptble/froxy"
	"github.com/skriptble/froxy/cmd/froxy/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/skriptble/froxy/cmd/froxy/Godeps/_workspace/src/github.com/spf13/viper"
)

var proxy froxy.ProxyBuilder

var rootCmd = &cobra.Command{
	Use:   "froxy",
	Short: "Froxy implements a reverse file proxy.",
	Long: `Froxy implements a reverse file proxy.
Use the server subcommand to run the actual server.`,
}

func init() {
	var serverCmd = &cobra.Command{
		Use:   "server [local directory]",
		Short: "Froxy server is a reverse file proxy",
		Long: `Froxy server is a reverse file proxy supporting a single http
interface for multiple types of file backends, including local
files and remote files.`,
		Run: server,
	}

	rootCmd.PersistentFlags().StringP("local-dir", "l", ".", "directory for the local file proxy [FROXY_LOCAL_DIR]")
	rootCmd.PersistentFlags().StringP("remote-url", "r", "", "url of the remote file proxy [FROXY_REMOTE_URL")
	rootCmd.PersistentFlags().IntP("port", "p", 8080, "port for the server [PORT]")
	rootCmd.AddCommand(serverCmd)

	// Setup viper flags and environment variables

	// local directory
	viper.BindPFlag("local-dir", rootCmd.PersistentFlags().Lookup("local-dir"))
	viper.BindEnv("local-dir", "FROXY_LOCAL_DIR")

	// remote url
	viper.BindPFlag("remote-url", rootCmd.PersistentFlags().Lookup("remote-url"))
	viper.BindEnv("remote-url", "FROXY_REMOTE_URL")

	// port
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindEnv("port", "PORT")
}

func main() {
	rootCmd.Execute()
}

func server(cmd *cobra.Command, args []string) {
	local := froxy.Dir(viper.GetString("local-dir"))
	if len(args) > 0 {
		local = froxy.Dir(args[0])
	}

	proxy = froxy.NewProxy()
	proxy.AddFileSource(local, "local")
	log.Printf("Local Proxy Directory: %v", local)

	if viper.GetString("remote-url") != "" {
		href, err := url.Parse(viper.GetString("remote-url"))
		if err != nil {
			log.Panicf("%v is not a valid url!", viper.GetString("remote-url"))
		}
		rmt := froxy.NewRemote(*href)
		proxy.AddFileSource(rmt, "remote")
		log.Printf("Remote URL Proxy: %v", href.String())
	}

	port := 8080
	if viper.IsSet("port") {
		port = viper.GetInt("port")
	}
	addr := ":" + strconv.Itoa(port)
	http.HandleFunc("/", handleRequest)
	log.Printf("Listening on port %v", addr)
	http.ListenAndServe(addr, nil)
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	paths := strings.SplitN(req.URL.Path, "/", 3)
	if len(paths) < 3 {
		// We don't have enough pieces of the url. Return a 404 Not Found.
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, froxy.NotFound)
		return
	}
	source := paths[1]
	name := paths[2]
	file, err := proxy.RetrieveFile(name, source)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, err.Error())
		log.Println(err)
		return
	}
	io.Copy(w, file)
}
