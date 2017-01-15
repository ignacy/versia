# Versia

Do you use [PaperTrail](https://github.com/airblade/paper_trail) to
track changes to your models?  Do you often get bug reports like this:
"I'm seeing this value here but it shouldn't be possible"?  Do you
start solving those by thinking "Hmm.. that is in fact strange.. maybe
someone set this value through Admin panel - and avoided some of the
regular validations"

If yes Versia can help you by making it easy to inspect changes made
to a model that PaperTrail keep.

# Configuration

You need to set those configuration settings:

        VERSIA_ADMIN_PASSWORD - basic auth password
        VERSIA_MODEL_NAME - model name you want to have this for
        VERSIA_PG_STRING - PostgreSQL connection string
        PORT - http port versia is supposed to start on


Versia is ready to be used on Heroku `heroku create` `git push heroku
master` and after that set configuration variables.


# Additional info

Connection strings for PostgreSQL are described [here](https://godoc.org/github.com/lib/pq)