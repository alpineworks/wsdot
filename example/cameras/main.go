package main

import (
	"fmt"
	"os"

	"alpineworks.io/wsdot"
	"alpineworks.io/wsdot/cameras"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("API_KEY environment variable is required")
	}

	// Create a new WSDOT client
	wsdotClient, err := wsdot.NewWSDOTClient(
		wsdot.WithAPIKey(apiKey),
	)

	if err != nil {
		panic(err)
	}

	// Create a new Cameras client
	camerasClient, err := cameras.NewCamerasClient(wsdotClient)
	if err != nil {
		panic(err)
	}

	// Get the cameras
	cameras, err := camerasClient.GetCameras()
	if err != nil {
		panic(err)
	}

	if len(cameras) > 0 {
		fmt.Println(cameras[0].CameraID)
		fmt.Println(cameras[0].Title)
		fmt.Println(cameras[0].ImageURL)
	}

	// Get a specific camera
	camera, err := camerasClient.GetCamera(cameras[0].CameraID)
	if err != nil {
		panic(err)
	}

	fmt.Println(camera.CameraID)
	fmt.Println(camera.Title)
	fmt.Println(camera.ImageURL)

}
