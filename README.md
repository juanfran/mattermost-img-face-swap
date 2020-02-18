# Mattermost FaceSwap Plugin

Change the faces of an image for the ones you want.

## Configuration

Clone the repository.

```
git clone git@github.com:juanfran/mattermost-img-face-swap.git
```

Create your configuration file. In this file you will configure each face that will be available to put over another one.

```
cd mattermost-img-face-swap
cp assets/faces.example.json assets/faces.json 
```

Open the configuration file `assets/faces.json` 

```json
{
    "faces": [
        {
            "name": "name",
            "images": [
                {
                    "path": "faces/face1.png",
                    "width": 80,
                    "paddingLeft": 30,
                    "paddingTop": 0
                }
            ]
        },
                {
            "name": "name2",
            "images": [
                {
                    "path": "faces/face2.png",
                },
                {
                    "path": "faces/face2.png",
                    "paddingLeft": 5
                }
            ]
        }
    ],
    "width": 80,
    "paddingLeft": 0,
    "paddingTop": 0
}
```

`width`, `paddingLeft`, `paddingTop` are parameters that you can configure globally or on each face that need a custom position or size. This modifiers are relative to the face width.

To add new face you have to add it to the `faces` array with a name that will be used later. In the `images` param you can add multiples images for each face and each one with their custom position and width if neccesary. The faces should be cropped and have a transparent background.

## Installation

Just run `make` in the project folder.

```
make
```

When the build is complete a file will be generated in the `dist` folder. You can upload this file in the Mattermost system console to install the plugin.

## Usage

Upload an image to a channel, then run the command `/faceswap` and a bot will answer with the image with the faces randomly replaces with the faces in your `assets/faces.json` file. If you want to choose which faces to use just run `/faceswap name1,name2`
