# Habitifocus

Sync OmniFocus tasks to Habitica using AppleScript and the Habitica JSON API.

## Installing

    go get github.com/BrianHicks/habitifocus

Then run:

    habitifocus --userid=youruserid --apikey=yourapikey

Where `youruserid` and `yourapikey` are your values from [Habitica's API page](https://habitica.com/#/options/settings/api).

If you don't want to provide these on the command line every time, create `~/.habitifocus.yaml` with the following contents:

    ---
    userid: youruserid
    apikey: yourapikey

## License

Habitifocus is licensed under a BSD 3-Clause license, located at [LICENSE](LICENSE)

That said, please don't use this if you don't know what's going on under the
covers. It could eat all your OmniFocus or Habitica tasks if there's a bug.
