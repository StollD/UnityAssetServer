package main

/*
 The configuration for the server
 */
type ConfigData struct {
     Host string `json:"host"`
     Port int32  `json:"port"`
     UnityVersions []string `json:"versions"`
     Modes map[string]BuildConfig `json:"modes"`
}

type BuildConfig struct {
     Directory string `json:"dir"`
     Files map[string]string `json:"files"`
     Output string `json:"output"`
     Command string `json:"cmd"`
}

/*
 The global instance of the configuration
 */
var Config ConfigData
