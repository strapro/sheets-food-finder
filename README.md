<!-- Add some text documenting the build step -->
# Sheets food finder

This is a command line tool written in Go which can be used to fetch some info from a Google sheets document

## Download

You can download the latest version from 

https://github.com/strapro/sheets-food-finder/releases

## Usage 

```console
./sheets-food-finder 'Your name'
```

## Building the project

To build the project, use the following command:

```console
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags "-X main.SpreadsheetId=<Your Spreadsheet ID> \
            -X sheetsFoodFinder/pkg/authhelper.ClientID=<Your Client ID> \
            -X sheetsFoodFinder/pkg/authhelper.ClientSecret=<Your Client Secret> \
            -X sheetsFoodFinder/pkg/authhelper.AuthURL=<Your Auth URL> \
            -X sheetsFoodFinder/pkg/authhelper.TokenURL=<Your Token URL>" \
  -o ./build sheets-food-finder
```

The `authhelper` build time constants can be acquired from Google Cloud Console APIS and services section

Once the code is build you need to add execute permission to the file

```console
chmod +x ./sheets-food-finder
```

## Running locally

To run locally you can create a `.env` file by copying the `.env.sample` and populating the values. Then you can execute

```console
run go sheetsfoodfinder.go 
```

## Caveats 

When executing the program for the first time, a browser will open requesting you to login to your google account in order to get readonly access to your google sheets

The program will then create a `token.json` file which will contain the `access_token` and the `refresh_token` for the selected account.

> :warning: **This info will be stored unencrypted in your filesystem** 

