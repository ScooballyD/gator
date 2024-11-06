# gator
Blog Aggregator

-In order to run gator you will need to install go1.23+ and Postgres-

intalling goose via go install will streamline much of the database setup.
run goose to v5.

to install gator, simply run go instal from the root of the program files.

gator will require a config file located at "~/.gatorconfig.json".
the file should be structured like this:
{
  "db_url": "connect_string",
  "current_user_name": "username_goes_here"
}

-replace "connect_string" with your connection string.
-it should loobe structured like this: protocol://username:password@host:port/database?sslmode=disable
-ex: postgres://postgres:pgpassword@localhost:5432/gator?sslmode=disable

==Commands==
gator login #
    -logs in #
gator register #
    -registers # as a user
gator reset
    -deletes all users
gator users
    -lists all users
gator agg #
    -continuosly saves posts at # interval from users feeds
    -interval should be structured like 30s or like 1m
gator addfeed # #
    -adds feed to database, requires input name and url
gator  feeds
    -lists all feeds and the users who added them
gator follow #
    -follows feed with matching url
gator following
    -lists all the feeds the current user is following
gator unfollow #
    -unfollows feed with matching url
browse #
    -lists the most recent # number of saved posts from the users followed feeds
    -input is optional, if non is given, # will default to 2