<p align="center">
    <a href="https://smallvictori.es"><img src="https://f.cloud.github.com/assets/846194/1235472/24f70d94-29a7-11e3-835a-84f55972b657.png" /></a>
</p>

This is service for [Small Victories](https://smallvictori.es) that
serves HTML pages created dynamically from a users Dropbox folder. The pages are retrieved from Redis after being stored there by
the [worker](https://github.com/pearkes/sv-fetcher) service.

It also contains the marketing website, Dropbox OAuth flow and help pages!

## Design

- Serve raw HTML out of Redis to incoming requests quickly.
- Allow a user to authenticate and go through the Dropbox OAuth flow,
catching and displaying issues to them along the way
- Show a "waiting" page when the initial sync is underway

## Hacking and Deploying

This service is written entirely in [Go](http://golang.org/), so you'll need to download
and install that.

You'll also need [Foreman](http://ddollar.github.io/foreman/) to start the application with the proper environment
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

To deploy:

    $ heroku create -b https://github.com/kr/heroku-buildpack-go.git
    ...
    $ git push heroku master
    ...

You'll need:

- Heroku account
- Heroku Postgres database
- Redis add-on of some kind
- Librato Metrics account

Use `heroku config:add` to add environment variables based on your
various tokens from the above to create the production environment,
as it is described in `.env.example`.

### Shared Credentials

Keep in mind these credentials should be shared with the [sv-fetcher](https://github.com/pearkes/sv-fetcher),
so you'll need to add the same environment variables there.

## Contributions

Small Victories being open source is mostly educational, as there is
unknown intent for further development.

If you're interested in maintaining or contributing to the project and
website, please contact us and we'll chat about it. Thanks!

[computers@smallvictori.es](mailto:computers@smallvictori.es)

## License

See [license](LICENSE.md) file.
