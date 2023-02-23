package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/k0kubun/pp"
	"github.com/nyushi/dockerhub-go"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use: "dockerhub",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unreleased")
	},
}

var ImageCmd = &cobra.Command{
	Use: "image",
}

var ImageSummaryCmd = &cobra.Command{
	Use:  "summary [image name]",
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		c := getClient(ctx)
		img := parseImage(args[0])
		summary, err := c.GetImageSummary(ctx, img.namespace, img.repository)
		if err != nil {
			log.Fatalf("failed to get summary: %s", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"", "number of images"})
		table.AppendBulk([][]string{
			{"active", fmt.Sprint(summary.Statistics.Active)},
			{"inactive", fmt.Sprint(summary.Statistics.Inactive)},
		})
		table.SetFooter([]string{"total", fmt.Sprint(summary.Statistics.Total)})
		table.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
		fmt.Printf("%s\n  active from %s\n", img.String(), summary.ActiveFrom)
		table.Render()
	},
}

var ImageDetailCmd = &cobra.Command{
	Use:  "detail [image name]",
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		c := getClient(ctx)
		img := parseImage(args[0])

		resp, err := c.GetImageDetails(ctx, img.namespace, img.repository)
		if err != nil {
			log.Fatalf("failed to get image detail: %s", err)
			return
		}
		pp.Println(resp)
	},
}

var ImageTagCmd = &cobra.Command{
	Use:  "tag [image name] [digest]",
	Args: cobra.MatchAll(cobra.ExactArgs(2)),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		c := getClient(ctx)
		img := parseImage(args[0])
		pp.Println(c.GetImageTags(ctx, img.namespace, img.repository, args[1]))
	},
}

var TagCmd = &cobra.Command{
	Use: "tag",
}

var TagListCmd = &cobra.Command{
	Use:  "list [image name]",
	Args: cobra.MatchAll(cobra.ExactArgs(1)),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		c := getClient(ctx)
		img := parseImage(args[0])

		tags := []*dockerhub.TagResponse{}
		resp, err := c.ListRepository(ctx, img.namespace, img.repository)
		for {
			if err != nil {
				log.Fatalf("failed to get image detail: %s", err)
				return
			}
			if resp == nil {
				break
			}
			for _, tag := range resp.Results {
				tags = append(tags, tag)
			}
			resp, err = resp.GetNext(ctx)
		}

		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Name < tags[j].Name
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"name", "last_updated"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t")
		table.SetNoWhiteSpace(true)
		for _, v := range tags {
			lastUpdated := ""
			if v.LastUpdated != nil {
				lastUpdated = *v.LastUpdated
			}
			table.Append([]string{v.Name, lastUpdated})
		}
		table.Render()
	},
}

func init() {
	viper.SetEnvPrefix("dockerhub")
	viper.AutomaticEnv()
	rootCmd.PersistentFlags().StringP("user", "", "", "docker hub user")
	rootCmd.PersistentFlags().StringP("pass", "", "", "docker hub pass")
	rootCmd.PersistentFlags().StringP("endpoint", "", "https://hub.docker.com/", "docker hub baseurl")
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("pass", rootCmd.PersistentFlags().Lookup("pass"))
	viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))

	ImageCmd.AddCommand(ImageSummaryCmd)
	ImageCmd.AddCommand(ImageDetailCmd)
	ImageCmd.AddCommand(ImageTagCmd)

	TagCmd.AddCommand(TagListCmd)

	rootCmd.AddCommand(ImageCmd)
	rootCmd.AddCommand(TagCmd)
	rootCmd.AddCommand(versionCmd)
}

type image struct {
	namespace  string
	repository string
	tag        string
}

func (i *image) String() string {
	s := fmt.Sprintf("%s/%s", i.namespace, i.repository)
	if i.tag != "" {
		s = fmt.Sprintf("%s:%s", s, i.tag)
	}
	return s
}

func parseImage(s string) *image {
	i := image{}
	a := strings.SplitN(s, "/", 2)
	if len(a) == 1 {
		i.namespace = "library"
		i.repository = a[0]
	} else {
		i.namespace = a[0]
		i.repository = a[1]
	}

	b := strings.SplitN(i.repository, ":", 2)
	if len(b) == 1 {
		return &i
	}

	i.repository = b[0]
	i.tag = b[1]
	return &i
}

func getClient(ctx context.Context) *dockerhub.Client {
	c := dockerhub.NewClient(viper.GetString("endpoint"))
	if err := c.UsersLogin(ctx, viper.GetString("user"), viper.GetString("pass")); err != nil {
		log.Fatalf("error at login: %s", err)
	}
	return c
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
