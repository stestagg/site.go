package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/stestagg/site.go/site"
	"github.com/stestagg/site.go/log"
	)


type options struct {
	Verbose bool
	Debug bool
	SiteRoot string
}

var Options options



var RootCmd = &cobra.Command{
	Use:   "sitego",
	Short: "Sitego is a static site generator",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if Options.Debug {
			log.Verbosity = 10
		}else if Options.Verbose{
			log.Verbosity = 2
		} else {
			log.Verbosity = 1
		}

		if Options.SiteRoot != "" {
			abspath, err := filepath.Abs(Options.SiteRoot)
			if err != nil {
				log.Panic("Could not find site root %s: %s", Options.SiteRoot, err)
			}
			site.Site.SetRoot(abspath)
			log.Info("Setting site root to %s", Options.SiteRoot)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: HEAD")
	},
}

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Find and list files that are prart of the site",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Finding Files")
		for node := range site.Site.DiscoverFiles() {
			log.Out(node.Path)
		}
	},
}

var buildCmd = &cobra.Command{
	Use: "build",
	Short: "Build the site",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Finding Files")
		for node := range site.Site.DiscoverFiles() {
			log.Out(node.Path)
		}
	},
}


func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(discoverCmd)
	RootCmd.AddCommand(buildCmd)
	RootCmd.PersistentFlags().BoolVarP(&Options.Verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().BoolVarP(&Options.Debug, "debug", "d", false, "debugging output")
	RootCmd.PersistentFlags().StringVarP(&Options.SiteRoot, "site", "s", "", "Site root")
}