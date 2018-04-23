package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/rs/xid"
	"github.com/urfave/cli"
)

const appVersion = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Author = "Repometric"
	app.Copyright = "all rights reserved 2018"
	app.Email = "hi@repomteric.com"
	app.Description = ""
	app.Name = "lhconf"
	app.Usage = "Linterhub Manager Core Component"
	app.Version = appVersion
	app.BashComplete = func(c *cli.Context) {
		cli.ShowCommandHelp(c, "lhconf")
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "project,p",
			Usage: "Path to file.",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"c"},
			Usage:   "Add engine or rule of ingore to config.",
			Subcommands: []cli.Command{
				{
					Name:    "engine",
					Aliases: []string{"e"},
					Usage:   "Add engine",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name,n",
							Usage: "Name of engine",
						},
						cli.StringFlag{
							Name:  "config,c",
							Usage: "Path to file config of engine",
						},
						cli.StringFlag{
							Name:   "install,i",
							Value:  "local",
							EnvVar: "local,global,container",
							Usage:  "Settings to run engine",
						},
						cli.BoolFlag{
							Name:  "activate,a",
							Usage: "Activate or deactivate engine for use",
						},
					},
					Action: func(c *cli.Context) error {
						var (
							name     = c.String("name")
							config   = c.String("config")
							instal   = c.String("install")
							activate = c.Bool("activate")
							path     = c.GlobalString("project")
						)

						pathValidate(path)

						if len(name) > 0 {
							con := getConfig(path)

							find := false
							for i := range con.Engines {
								if con.Engines[i].Name == name {
									find = true
									break
								}
							}

							if !find {
								con.Engines = append(con.Engines, Engine{
									Name:    name,
									Locally: instal,
									Active:  activate,
									Config:  config,
								})

								saveConfig(path, con)

								fmt.Println("Engine", name, "added.", "Total count", len(con.Engines))
							}

						} else {
							log.Fatal("Name обязательное поле")
							os.Exit(1)
						}
						return nil
					},
					BashComplete: func(c *cli.Context) {
						cli.ShowCommandHelp(c, "engine")
					},
				},
				{
					Name:    "ignore",
					Aliases: []string{"i"},
					Usage:   "Add rule of ingore",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "mask,m",
							Usage: "Mask for ingore",
						},
						cli.IntFlag{
							Name:  "line,l",
							Usage: "Line in file for ignore",
						},
						cli.StringFlag{
							Name:  "ruleId,r",
							Usage: "Rule of engine for ignore",
						},
					},
					Action: func(c *cli.Context) error {
						var (
							mask   = c.String("mask")
							line   = c.Int("line")
							ruleid = c.String("ruleId")
							path   = c.GlobalString("project")
						)

						pathValidate(path)

						if len(mask) > 0 || line > 0 || len(ruleid) > 0 {
							con := getConfig(path)

							id := xid.New()
							con.Ingores = append(con.Ingores, Ingore{
								Mask:     mask,
								Line:     line,
								Ruleid:   ruleid,
								Ingoreid: id.String(),
							})

							saveConfig(path, con)

							fmt.Println("Rule of ingore added.", "Total count", len(con.Ingores))
						}
						return nil
					},
					BashComplete: func(c *cli.Context) {
						cli.ShowCommandHelp(c, "ignore")
					},
				},
			},
			BashComplete: func(c *cli.Context) {
				cli.ShowCommandHelp(c, "add")
			},
		},
		{
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "Remove engine or ignore rule",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "engine,e",
					Usage: "Remove engine by name",
				},
				cli.StringSliceFlag{
					Name:  "ingoreId,i",
					Usage: "Remove rule of ignore by `ignoreId`",
				},
			},
			Action: func(c *cli.Context) error {
				var (
					engine = c.String("engine")
					ignore = c.StringSlice("ingoreId")
					path   = c.GlobalString("project")
				)

				pathValidate(path)

				con := getConfig(path)

				if len(engine) > 0 {

					for i := range con.Engines {
						if con.Engines[i].Name == engine {
							con.Engines = append(con.Engines[:i], con.Engines[i+1:]...)
							break
						}
					}

					saveConfig(path, con)

					fmt.Println("Engine", engine, "removed.", "Total count", len(con.Engines))
					return nil
				}

				if len(ignore) > 0 {

					for i := range con.Ingores {
						for id := range ignore {
							if con.Ingores[i].Ingoreid == ignore[id] {
								con.Ingores = append(con.Ingores[:i], con.Ingores[i+1:]...)
							}
						}
					}

					saveConfig(path, con)

					fmt.Println("Rule of ignore", ignore, "removed.", "Total count", len(con.Ingores))
					return nil
				}
				return nil
			},
			BashComplete: func(c *cli.Context) {
				cli.ShowCommandHelp(c, "remove")
			},
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "Get engines, rules of ingore or all configuration",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "list,l",
					Usage:  "Get all config",
					Value:  "all",
					EnvVar: "all,engines,rules",
				},
				cli.StringSliceFlag{
					Name:   "keys,k",
					Usage:  "Keys for output",
					Value:  &cli.StringSlice{"all"},
					EnvVar: "all,id,name,mask,line,ruleId,activate,config,install",
				},
				cli.StringFlag{
					Name:  "engine,e",
					Usage: "Get engine by Name",
				},
				cli.StringSliceFlag{
					Name:  "ingoreId,i",
					Usage: "Get ingore rule by ingoreId",
				},
			},
			Action: func(c *cli.Context) error {
				var (
					list   = c.String("list")
					engine = c.String("engine")
					ignore = c.StringSlice("ingoreId")
					path   = c.GlobalString("project")
				)

				pathValidate(path)

				con := getConfig(path)

				if len(engine) > 0 {

					for i := range con.Engines {
						if con.Engines[i].Name == engine {
							engineJSON, _ := json.Marshal(con.Engines[i])
							fmt.Println(string(engineJSON))
							break
						}
					}
					return nil
				}

				if len(ignore) > 0 {

					for i := range con.Engines {
						for id := range ignore {
							if con.Ingores[i].Ingoreid == ignore[id] {
								ingoreJSON, _ := json.Marshal(con.Ingores[i])
								fmt.Println(string(ingoreJSON))
							}
						}
					}
					return nil
				}

				if len(list) > 0 {

					switch list {
					case "all":
						allJSON, _ := json.Marshal(con)
						fmt.Println(string(allJSON))
					case "engines":
						enginesJSON, _ := json.Marshal(con.Engines)
						fmt.Println(string(enginesJSON))
					case "ingores":
						ingoresJSON, _ := json.Marshal(con.Ingores)
						fmt.Println(string(ingoresJSON))
					}
					return nil
				}
				return nil
			},
			BashComplete: func(c *cli.Context) {
				cli.ShowCommandHelp(c, "get")
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getConfig(path string) Config {
	raw, err := ioutil.ReadFile(path)

	check(err)

	var c Config
	json.Unmarshal(raw, &c)
	return c
}

func saveConfig(path string, con Config) bool {
	conJSON, _ := json.Marshal(con)
	err := os.Remove(path)
	check(err)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	check(err)
	f.Write(conJSON)
	defer f.Close()
	return true
}

func pathValidate(path string) {
	if len(path) == 0 {
		log.Fatal("File path doesn't find")
		os.Exit(1)
	}
}
