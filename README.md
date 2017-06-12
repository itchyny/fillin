# fillin [![Travis Build Status](https://travis-ci.org/itchyny/fillin.svg?branch=master)](https://travis-ci.org/itchyny/fillin)
### fill-in your command and execute
#### ― _separate action and environment of your command!_ ―

## Motivation
We rely on shell history in our terminal operation.
We search from our shell history and execute commands dozens times in a day.

However, shell history sometimes contains authorization tokens that we don't care while searching the commands.
Some incremental fuzzy searchers have troubles when there are many random tokens in the shell history.
Yeah, I know that I should not type a authorization token directly in the command line, but it's much easier than creating some shell script snippets.

Another hint to implement `fillin` is that programmers execute same commands switching servers.
We do not just login with `ssh {{hostname}}`, we also connect to the database with `psql -h {{psql:hostname}} -U {{psql:username}} -d {{psql:dbname}}` and to Redis server with `redis-cli -h {{redis:hostname}} -p {{redis:port}}`.
We switch the host argument from the localhost (you may omit this), staging and production servers.

Some command line tools allow us to login cloud services and retrieve data from our terminal.
Most of such cli tools accept an option to switching between the account.
For example, AWS command line tool has `--profile` option.
Other typical names of options are `--region`, `--conf` and `--account`.
When we specify these options directly, there are quadratic number of commands; the number of accounts times the number of actions.
The `fillin` allows us to save the command something like `aws --profile {{aws:profile}} ec2 describe-instances` so we'll not be bothered by the quadratic combinations of commands while searching through the shell history.

The core concept of `fillin` lies in separating the action (do what) in the command and the environment (to where).
With this `fillin` command line tool, you can

- make your commands reusable and it will make incremental shell history searching easy.
- fill in the template variables interactively and their history will be stored locally.
- invoke the same action switching multiple environment (local, staging and production servers, configuration path, cloud service accounts or whatever)

## Installation
### Homebrew
```sh
 $ brew install itchyny/fillin/fillin
```

### Download binary from GitHub Releases
[Releases・itchyny/fillin - GitHub](https://github.com/itchyny/fillin/releases)

### Build from source
```sh
 $ go get -u github.com/itchyny/fillin
```

## Usage
The interface of the `fillin` command is very simple.
Prepend `fillin` to the command and create template variables with `{{...}}`.
So the hello world for the `fillin` command is as follows.
```sh
 $ fillin echo {{message}}
message: Hello, world!        # you type here
Hello, world!                 # fillin executes: echo 'Hello, world!'
```
The value of `message` variable is stored locally.
You can use the recently used value with the upwards key.

The `{{message}}` is a template part of the command.
As the identifier, you can use alphabets, numbers, underscore and hyphen.
Thus `{{sample-id}}`, `{{SAMPLE_ID}}`, `{{X01}}` and `{{FOO_example-identifier0123}}` are all valid template parts.

One of the important features of `fillin` is variable scope grouping.
Let's look into more practical example.
When you connect to PostgreSQL server, you can use:
```sh
 $ fillin psql -h {{psql:hostname}} -U {{psql:username}} -d {{psql:dbname}}
[psql] hostname: localhost
[psql] username: example-user
[psql] dbname: example-db
```
What's the benefit of `psql:` prefix?
You'll notice the answer when you execute the command again:
```sh
 $ fillin psql -h {{psql:hostname}} -U {{psql:username}} -d {{psql:dbname}}
[psql] hostname, username, dbname: localhost, example-user, example-db   # you can select the most recently used entry with the upwards key
```
The identifiers with the same scope name (`psql` scope here) can be selected as pairs.
You can input individual values to create a new pair after skipping the multi input prompt.
```sh
 $ fillin psql -h {{psql:hostname}} -U {{psql:username}} -d {{psql:dbname}}
[psql] hostname, username, dbname:             # just type enter to input values for each identifiers
[psql] hostname: example.org
[psql] username: example-org-user
[psql] dbname: example-org-db
```

The scope grouping behaviour is useful with some authorization keys.
```sh
 $ fillin curl {{example-api:base-url}}/api/1/example/info -H 'Authorization: Bearer {{example-api:access-token}}'
[example-api] base-url, access-token: example.com, accesstokenabcde012345
```
The `base-url` and `access-token` are stored as tuples so you can easily switch between local, staging and production environment authorization.
Without the grouping behaviour, variable history searching will lead you to an unmatched pair of `base-url` and `access-token`.
Since the curl endpoint are stored in the shell history and authorization keys are stored in `fillin` history, we'll not be bothered by the quadratic number of the command history.

In order to have the benefit of this grouping behaviour, it's strongly recommended to prepend the scope name.
The `psql:` prefix on connecting to PostgreSQL database server, `redis:` prefix for Redis server are useful best practice in my opinion.

## Pipe, redirection and subshell
The terminal interface of `fillin` has problem with some shell functionalities.
For example, the following command gets stuck the terminal interface.
```sh
 $ fillin echo {{message}} | jq .
{}^M^M^C
```
This is because the interface of `fillin` is rely on the standard output.
Instead of connecting the standard output of `fillin` to another command, pass the pipe character as an argument.
```sh
 $ fillin echo {{message}} \| jq .
 $ # or
 $ fillin echo {{message}} '|' jq .
```
Another solution is to quote the whole command.
```sh
 $ fillin 'echo {{message}} | jq .'
```
This is much simple and general method because you can use on subshells as well.
```sh
 $ fillin 'echo $(cat {{foo}} {{bar}})'
 $ fillin '(echo {{foo}}; cat {{bar}}) | grep func'
```

The same problem occurs with redirection so escape `>` or quote the whole command.
```sh
 $ fillin echo {{message}} \> /tmp/message
 $ fillin echo {{message}} '>' /tmp/message
 $ fillin 'echo {{message}} > /tmp/message'
```

## Disclaimer
This command line tool is in its early developing stage.
The user interface may be changed without any announcement.

This tool is not an encryption tool.
The command saves the inputted values in a JSON file with no encryption.

## Bug Tracker
Report bug at [Issues・itchyny/fillin - GitHub](https://github.com/itchyny/fillin/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
