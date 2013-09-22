# sv-frontend

This is service for [Small Victories](https://smallvictori.es) that
serves HTML pages created dynamically from a user's Drobpox folder.

It also contains the marketing website, Dropbox OAuth flow and help pages.

## Design

- Serve raw HTML out of Redis to incoming requests quickly.
- Allow a user to authenticate and go through the Dropbox OAuth flow,
catching and displaying issues to them along the way
- Show a "waiting" page when the initial sync is underway

## Hacking and Deploying

This service is written entirely in [Go](), so you'll need to download
and install that.

You'll also need [Foreman]() to start the application with the proper environment
variables.

To configure the service, you can use the configuration example and create
a `.env` file containing your credentials.

To run, you then use foreman:

    $ go build
    $ foreman run ./sv-frontend
    ...

To run the tests:

    $ go test
    ...

## Contributions

Small Victories being open source is mostly educational, as there is
unknown intent for further development.

If you're interested in maintaining or contributing to the project and
website, please contact us and we'll chat about it. Thanks!

[computers@smallvictori.es](mailto:computers@smallvictori.es)

## License

See [license](LICENSE.md) file.
