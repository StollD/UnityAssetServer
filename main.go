package main

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "github.com/spf13/cast"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)


func main() {
    // Load the configuration
    if err := LoadSettings("config/config.json", &Config); err != nil {
        // handle error
        log.Fatal(err)
        return
    }

    // Start the webserver
    e := echo.New()
    e.Debug = true
    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Routes
    e.POST("/:type", build)

    // Start server
    e.Logger.Fatal(e.Start(Config.Host + ":" + cast.ToString(Config.Port)))
}

/*
 Path: /:type
 Request: GET
 Params: type - which asset should be built
 */
func build(c echo.Context) error {
    // Get the type of the asset
    assetType := c.Param("type")

    // Check if the type was valid
    if _, ok := Config.Modes[assetType]; !ok {
        return c.String(http.StatusNotFound, "Asset type not found")
    }
    buildcfg := Config.Modes[assetType]

    // The asset type exists, check if the unity version is supported
    unity := c.FormValue("unity")
    found := false
    for _, value := range Config.UnityVersions {
        if value == unity {
            found = true
        }
    }
    if !found {
        return c.String(http.StatusNotFound, "Invalid Unity Version!")
    }
    directory := strings.Replace(buildcfg.Directory, "{UNITY}", unity, -1)

    // Get the files from the request
    for name, path := range buildcfg.Files {
        data, err := c.FormFile(name)
        if err != nil {
            return err
        }
        handle, err := data.Open()
        if err != nil {
            return err
        }

        // Save the file to its destination
        dest, err := os.Create(filepath.Join(directory, path))
        if err != nil {
            return err
        }

        // Copy the file
        _, err = io.Copy(dest, handle)
        if err != nil  {
            return err
        }
    }

    // Everything is fine, we can run the build
    cmd := exec.Command("/bin/sh", "-c", "u3d --verbose --trace -- -executeMethod " + buildcfg.Command + " -quit -batchmode -nographics")
    cmd.Dir = directory
    err := cmd.Run()
    if err != nil {
        return err
    }

    // Return the generated file
    return c.File(filepath.Join(directory, buildcfg.Output))
}



func LoadSettings(filename string, out interface{}) error {
    file, _ := os.Open(filename)
    defer file.Close()
    decoder := json.NewDecoder(file)
    return decoder.Decode(out)
}
