# fillin [![Travis Build Status](https://travis-ci.org/itchyny/fillin.svg?branch=master)](https://travis-ci.org/itchyny/fillin)
fill-in your command line

## Motivation
We rely on shell history very much.
We search from shell history and execute the command dozens times in a day.

However, shell history sometimes contains some authorization token that we don't care while searching the commands.
Some incremental fuzzy searchers have troubles when there are many random tokens in shell history.
Yeah, I know that I should not type some authorization token directly into command line, but it's much easier than creating some shell script snippets.

Another hint to implement `fillin` is that programmers execute same commands switching hosts.
We do not just login with `ssh {{host}}`, we also connect to database with `psql -h {{hostname}} {{dbname}} -U {{username}}` and to Redis host with `redis-cli -h {{hostname}} -p {{port}}`.
We switch the host argument from the localhost (you may omit this), staging hosts and production hosts.

The main idea is that splitting the command history and template variable history.
With `fillin` command line tool, you can

- make your command template and it will make incremental shell history searching easy.
- fill in template variables interactively and its history will be stored locally.

## Usage
The interface of `fillin` command is very simple.
Prepend `fillin` and create template variables with `{{...}}`.
So the hello world for `fillin` command goes like as follows.
```sh
 $ fillin echo {{message}}
message: Hello, world!        # you type here
Hello, world!                 # fillin executes: echo Hello, world!
```
The value of `message` variable is stored locally.
You can use the lastly used value with the upwards key (this may be replaced with more rich interface in the future but I'm not sure).

The `{{message}}` is called as template parts of the command.
As the identifier, you can use alphabets, numbers, underscore and hypens.
Thus `{{sample-id}}`, `{{SAMPLE_ID}}`, `{{X01}}` and `{{FOO_example-identifier0123}}` are all valid template parts.

Another important feature of `fillin` is scope grouping of the variables.
Let's look into more practical example.
When you connect to PostgreSQL host, you can use:
```sh
 $ fillin psql -h {{psql:hostname}} {{psql:dbname}} -U {{psql:username}}
[psql] hostname: example.com
[psql] dbname: example-db
[psql] username: example-user
```
What's the benefit of `psql:` prefix?
You'll notice the answer when you execute the command again:
```sh
 $ fillin psql -h {{psql:hostname}} {{psql:dbname}} -U {{psql:username}}
[psql] hostname, dbname, username: example.com, example-db, example-user   # you can select the most recently used entry with the upwards key
```
The identifiers with the same scope name (`psql` scope here) can be selected as pairs.
You can input each values after clearing the input.
```sh
 $ fillin psql -h {{psql:hostname}} {{psql:dbname}} -U {{psql:username}}
[psql] hostname, dbname, username:             # just type enter to input values for each identifiers
[psql] hostname: example.org
[psql] dbname: example-org-db
[psql] username: example-org-user
```

This scope grouping behaviour is very useful with some authorization keys.
```sh
 $ fillin curl {{example-api:hostname}}/api/example -H 'Api-Key: {{example-api:api-key}}'
[example-api] hostname, api-key: example.com, apikeyabcde012345
```
The `host` and `api-key` are stored as tuples so you can easily switch local, staging and production environment authorization.

In order to have the benefit of this grouping behaviour, it's strongly recommended to prepend the scope name.
The `psql:` prefix on connecting to PostgreSQL database host, `redis:` prefix for Redis host are useful best practice in my opinion.

## Disclaimer
This command line utility is in its early developing stage.
The user interface may be changed without any announcement.

## Installation
### Download binary from GitHub Releases
[Releases・itchyny/fillin - GitHub](https://github.com/itchyny/fillin/releases)

### Build from source
```bash
 $ go get -u github.com/itchyny/fillin
```

## Bug Tracker
Report bug at [Issues・itchyny/fillin - GitHub](https://github.com/itchyny/fillin/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
